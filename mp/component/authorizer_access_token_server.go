// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package component

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/chanxuehong/wechat/mp"
)

var _ mp.AccessTokenServer = (*AuthorizerAccessTokenServer)(nil)

// authorizer_access_token 中控服务器, mp.AccessTokenServer 的简单实现.
//  NOTE:
//  1. 用于单进程环境.
//  2. 因为 AuthorizerAccessTokenServer 同时也是一个简单的中控服务器, 而不是仅仅实现 mp.AccessTokenServer 接口,
//     所以整个系统只能存在一个 AuthorizerAccessTokenServer 实例!
type AuthorizerAccessTokenServer struct {
	client          *Client
	authorizerAppId string

	resetTickerChan chan time.Duration // 用于重置 tokenDaemon 里的 ticker

	tokenGet struct {
		sync.Mutex
		LastTokenInfo AuthorizerAccessTokenInfo // 最后一次成功从微信服务器获取的 authorizer_access_token 信息
		LastTimestamp int64                     // 最后一次成功从微信服务器获取 authorizer_access_token 的时间戳
	}

	tokenCache struct {
		sync.RWMutex
		Token        string
		RefreshToken string // 最新的 authorizer_refresh_token
	}
}

// 创建一个新的 AuthorizerAccessTokenServer.
func NewAuthorizerAccessTokenServer(clt *Client, authorizerAppId, authorizerRefreshToken string) (srv *AuthorizerAccessTokenServer) {
	if clt == nil {
		panic("nil Client")
	}

	srv = &AuthorizerAccessTokenServer{
		client:          clt,
		authorizerAppId: authorizerAppId,
		resetTickerChan: make(chan time.Duration),
	}
	srv.tokenCache.RefreshToken = authorizerRefreshToken

	go srv.tokenDaemon(time.Hour * 24) // 启动 tokenDaemon
	return
}

func (srv *AuthorizerAccessTokenServer) TagCE90001AFE9C11E48611A4DB30FED8E1() {}

// 獲取 authorizer_access_token
func (srv *AuthorizerAccessTokenServer) Token() (token string, err error) {
	srv.tokenCache.RLock()
	token = srv.tokenCache.Token
	srv.tokenCache.RUnlock()

	if token != "" {
		return
	}
	return srv.TokenRefresh()
}

// 刷新 authorizer_access_token
func (srv *AuthorizerAccessTokenServer) TokenRefresh() (token string, err error) {
	tokenInfo, cached, err := srv.getToken()
	if err != nil {
		return
	}
	if !cached {
		srv.resetTickerChan <- time.Duration(tokenInfo.ExpiresIn) * time.Second
	}
	token = tokenInfo.Token
	return
}

func (srv *AuthorizerAccessTokenServer) tokenDaemon(tickDuration time.Duration) {
NEW_TICK_DURATION:
	ticker := time.NewTicker(tickDuration)

	for {
		select {
		case tickDuration = <-srv.resetTickerChan:
			ticker.Stop()
			goto NEW_TICK_DURATION

		case <-ticker.C:
			tokenInfo, cached, err := srv.getToken()
			if err != nil {
				break
			}
			if !cached {
				newTickDuration := time.Duration(tokenInfo.ExpiresIn) * time.Second
				if tickDuration != newTickDuration {
					tickDuration = newTickDuration
					ticker.Stop()
					goto NEW_TICK_DURATION
				}
			}
		}
	}
}

type AuthorizerAccessTokenInfo struct {
	Token        string `json:"authorizer_access_token"`
	ExpiresIn    int64  `json:"expires_in"` // 有效时间, seconds
	RefreshToken string `json:"authorizer_refresh_token"`
}

// 从微信服务器获取 authorizer_access_token.
//  同一时刻只能一个 goroutine 进入, 防止没必要的重复获取.
func (srv *AuthorizerAccessTokenServer) getToken() (token AuthorizerAccessTokenInfo, cached bool, err error) {
	srv.tokenGet.Lock()
	defer srv.tokenGet.Unlock()

	timeNowUnix := time.Now().Unix()

	// 在收敛周期内直接返回最近一次获取的 authorizer_access_token, 这里的收敛时间设定为4秒.
	if n := srv.tokenGet.LastTimestamp; n <= timeNowUnix && timeNowUnix < n+4 {
		// 因为只有成功获取后才会更新 srv.tokenGet.LastTimestamp, 所以这些都是有效数据
		token = AuthorizerAccessTokenInfo{
			Token:        srv.tokenGet.LastTokenInfo.Token,
			ExpiresIn:    srv.tokenGet.LastTokenInfo.ExpiresIn - timeNowUnix + n,
			RefreshToken: srv.tokenGet.LastTokenInfo.RefreshToken,
		}
		cached = true
		return
	}

	srv.tokenCache.RLock()
	authorizerRefreshToken := srv.tokenCache.RefreshToken
	srv.tokenCache.RUnlock()

	request := struct {
		ComponentAppId         string `json:"component_appid"`
		AuthorizerAppId        string `json:"authorizer_appid"`
		AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
	}{
		ComponentAppId:         srv.client.AppId,
		AuthorizerAppId:        srv.authorizerAppId,
		AuthorizerRefreshToken: authorizerRefreshToken,
	}

	var result struct {
		mp.Error
		AuthorizerAccessTokenInfo
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token="
	if err = srv.client.PostJSON(incompleteURL, &request, &result); err != nil {
		srv.tokenCache.Lock()
		srv.tokenCache.Token = ""
		srv.tokenCache.Unlock()
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		srv.tokenCache.Lock()
		srv.tokenCache.Token = ""
		srv.tokenCache.Unlock()

		err = &result.Error
		return
	}

	// 由于网络的延时, authorizer_access_token 过期时间留了一个缓冲区
	switch {
	case result.ExpiresIn > 31556952: // 60*60*24*365.2425
		srv.tokenCache.Lock()
		srv.tokenCache.Token = ""
		srv.tokenCache.Unlock()

		err = errors.New("expires_in too large: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
	case result.ExpiresIn > 60*60:
		result.ExpiresIn -= 60 * 10
	case result.ExpiresIn > 60*30:
		result.ExpiresIn -= 60 * 5
	case result.ExpiresIn > 60*5:
		result.ExpiresIn -= 60
	case result.ExpiresIn > 60:
		result.ExpiresIn -= 10
	default:
		srv.tokenCache.Lock()
		srv.tokenCache.Token = ""
		srv.tokenCache.Unlock()

		err = errors.New("expires_in too small: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
	}

	// 更新 tokenGet 信息
	srv.tokenGet.LastTokenInfo = result.AuthorizerAccessTokenInfo
	srv.tokenGet.LastTimestamp = timeNowUnix

	// 更新缓存
	srv.tokenCache.Lock()
	srv.tokenCache.Token = result.AuthorizerAccessTokenInfo.Token
	if authorizerRefreshToken := result.AuthorizerAccessTokenInfo.RefreshToken; authorizerRefreshToken != "" {
		srv.tokenCache.RefreshToken = authorizerRefreshToken
	}
	srv.tokenCache.Unlock()

	token = result.AuthorizerAccessTokenInfo
	return
}
