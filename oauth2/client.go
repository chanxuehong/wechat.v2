// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	CreateAt     int64  `json:"create_at"`  // 创建时间, unixtime, 分布式系统要求时间同步, 建议使用 NTP
	ExpiresIn    int64  `json:"expires_in"` // 超时时间, seconds
	RefreshToken string `json:"refresh_token"`

	OpenId  string   `json:"openid"`
	UnionId string   `json:"unionid,omitempty"`
	Scopes  []string `json:"scopes,omitempty"` // 用户授权的作用域
}

// 判断 Token.AccessToken 是否过期, 过期返回 true, 否则返回 false
func (token *Token) AccessTokenExpired() bool {
	return time.Now().Unix() >= token.CreateAt+token.ExpiresIn
}

type TokenStorage interface {
	Get() (*Token, error)
	Put(*Token) error
}

type Client struct {
	Config Config

	// TokenStorage, Token 两个字段正常情况下只用指定一个, 如果两个同时被指定了, 优先使用 TokenStorage;
	TokenStorage TokenStorage
	Token        *Token // Client 自动将最新的 Token 更新到此字段, 不管 Token 字段一开始是否被指定!!!

	HttpClient *http.Client // 如果 HttpClient == nil 则默认用 http.DefaultClient
}

func (clt *Client) getToken() (tk *Token, err error) {
	if clt.TokenStorage != nil {
		if tk, err = clt.TokenStorage.Get(); err != nil {
			return
		}
		if tk == nil {
			err = errors.New("Incorrect TokenStorage.Get()")
			return
		}
		clt.Token = tk // update local
		return
	}

	tk = clt.Token
	if tk == nil {
		err = errors.New("nil TokenStorage and nil Token")
		return
	}
	return
}

func (clt *Client) putToken(tk *Token) (err error) {
	if clt.TokenStorage != nil {
		if err = clt.TokenStorage.Put(tk); err != nil {
			return
		}
	}
	clt.Token = tk
	return
}

func (clt *Client) httpClient() *http.Client {
	if clt.HttpClient != nil {
		return clt.HttpClient
	}
	return http.DefaultClient
}

func (clt *Client) getJSON(url string, response interface{}) (err error) {
	httpResp, err := clt.httpClient().Get(url)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	return json.NewDecoder(httpResp.Body).Decode(response)
}
