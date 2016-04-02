// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package oauth2

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chanxuehong/wechat/mp"
)

type Token struct {
	AccessToken  string // 网页授权接口调用凭证, 注意: 此access_token与基础支持的access_token不同
	ExpiresAt    int64  // 过期时间, unixtime, 分布式系统要求时间同步, 建议使用 NTP
	RefreshToken string // 用户刷新access_token

	OpenId  string   // 用户唯一标识, 请注意, 在未关注公众号时, 用户访问公众号的网页, 也会产生一个用户和公众号唯一的OpenID
	UnionId string   // UnionID机制
	Scopes  []string // 用户授权的作用域
}

// 判断 Token.AccessToken 是否过期, 过期返回 true, 否则返回 false
func (token *Token) AccessTokenExpired() bool {
	return time.Now().Unix() >= token.ExpiresAt
}

// 通过code换取网页授权access_token.
//  返回的 token == clt.Token
func (clt *Client) Exchange(code string) (token *Token, err error) {
	if clt.Config == nil {
		err = errors.New("nil Config")
		return
	}

	var tk *Token
	if clt.TokenStorage != nil {
		if tk, _ = clt.TokenStorage.Get(); tk == nil {
			tk = clt.Token
		} else {
			clt.Token = tk // update local
		}
	} else {
		tk = clt.Token
	}
	if tk == nil {
		tk = new(Token)
	}

	if err = clt.updateToken(tk, clt.Config.ExchangeTokenURL(code)); err != nil {
		return
	}

	if err = clt.putToken(tk); err != nil {
		return
	}
	token = tk
	return
}

// 刷新access_token(如果需要).
//  返回的 token == clt.Token
func (clt *Client) TokenRefresh() (token *Token, err error) {
	if clt.Config == nil {
		err = errors.New("nil Config")
		return
	}

	tk, err := clt.getToken()
	if err != nil {
		return
	}

	return clt.tokenRefresh(tk)
}

func (clt *Client) tokenRefresh(tk *Token) (token *Token, err error) {
	if err = clt.updateToken(tk, clt.Config.RefreshTokenURL(tk.RefreshToken)); err != nil {
		return
	}

	if err = clt.putToken(tk); err != nil {
		return
	}
	token = tk
	return
}

// 从服务器获取新的 token 更新 tk, 通过 code 换取 token 或者刷新 token
func (clt *Client) updateToken(tk *Token, url string) (err error) {
	var result struct {
		mp.Error

		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenId       string `json:"openid"`
		UnionId      string `json:"unionid"`
		Scope        string `json:"scope"`
	}
	if err = clt.getJSON(url, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		return &result.Error
	}

	// 由于网络的延时, 分布式服务器之间的时间可能不是绝对同步, access_token 过期时间留了一个缓冲区;
	switch {
	case result.ExpiresIn > 31556952: // 60*60*24*365.2425
		err = errors.New("expires_in too large: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
	case result.ExpiresIn > 60*60:
		result.ExpiresIn -= 60 * 20
	case result.ExpiresIn > 60*30:
		result.ExpiresIn -= 60 * 10
	case result.ExpiresIn > 60*15:
		result.ExpiresIn -= 60 * 5
	case result.ExpiresIn > 60*5:
		result.ExpiresIn -= 60
	case result.ExpiresIn > 60:
		result.ExpiresIn -= 20
	default:
		err = errors.New("expires_in too small: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
	}

	tk.AccessToken = result.AccessToken
	if result.RefreshToken != "" {
		tk.RefreshToken = result.RefreshToken
	}
	tk.ExpiresAt = time.Now().Unix() + result.ExpiresIn

	tk.OpenId = result.OpenId
	tk.UnionId = result.UnionId

	strArr := strings.Split(result.Scope, ",")
	tk.Scopes = make([]string, 0, len(strArr))
	for _, str := range strArr {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		tk.Scopes = append(tk.Scopes, str)
	}
	return
}

// 检验授权凭证(access_token)是否有效.
func (clt *Client) CheckAccessTokenValid() (valid bool, err error) {
	if clt.Config == nil {
		err = errors.New("nil Config")
		return
	}

	tk, err := clt.getToken()
	if err != nil {
		return
	}

	var result mp.Error

	_url := "https://api.weixin.qq.com/sns/auth?access_token=" + url.QueryEscape(tk.AccessToken) +
		"&openid=" + url.QueryEscape(tk.OpenId)
	if err = clt.getJSON(_url, &result); err != nil {
		return
	}

	switch result.ErrCode {
	case mp.ErrCodeOK:
		valid = true
		return
	case 40001:
		//valid = false
		return
	default:
		err = &result
		return
	}
}
