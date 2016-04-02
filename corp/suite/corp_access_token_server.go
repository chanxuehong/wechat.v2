// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package suite

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/chanxuehong/wechat/corp"
)

var _ corp.AccessTokenServer = (*CorpAccessTokenServer)(nil)

// corp.AccessTokenServer 的简单实现.
//  NOTE:
//  1. 用于单进程环境.
//  2. 因为 CorpAccessTokenServer 同时也是一个简单的中控服务器, 而不是仅仅实现 corp.AccessTokenServer 接口,
//     所以整个系统只能存在一个 CorpAccessTokenServer 实例!
type CorpAccessTokenServer struct {
	client        *Client
	authCorpId    string
	permanentCode string

	resetTickerChan chan time.Duration // 用于重置 tokenDaemon 里的 ticker

	tokenGet struct {
		sync.Mutex
		LastTokenInfo CorpAccessTokenInfo // 最后一次成功从微信服务器获取的 access_token 信息
		LastTimestamp int64               // 最后一次成功从微信服务器获取 access_token 的时间戳
	}

	tokenCache struct {
		sync.RWMutex
		Token string
	}
}

// 创建一个新的 CorpAccessTokenServer.
func NewCorpAccessTokenServer(clt *Client, authCorpId, permanentCode string) (srv *CorpAccessTokenServer) {
	if clt == nil {
		panic("nil Client")
	}

	srv = &CorpAccessTokenServer{
		client:          clt,
		authCorpId:      authCorpId,
		permanentCode:   permanentCode,
		resetTickerChan: make(chan time.Duration),
	}

	go srv.tokenDaemon(time.Hour * 24) // 启动 tokenDaemon
	return
}

func (srv *CorpAccessTokenServer) Tag6D89F2E2FE9811E49EAAA4DB30FED8E1() {}

func (srv *CorpAccessTokenServer) Token() (token string, err error) {
	srv.tokenCache.RLock()
	token = srv.tokenCache.Token
	srv.tokenCache.RUnlock()

	if token != "" {
		return
	}
	return srv.TokenRefresh()
}

func (srv *CorpAccessTokenServer) TokenRefresh() (token string, err error) {
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

func (srv *CorpAccessTokenServer) tokenDaemon(tickDuration time.Duration) {
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

type CorpAccessTokenInfo struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"` // 有效时间, seconds
}

// 从微信服务器获取 corp_access_token.
//  同一时刻只能一个 goroutine 进入, 防止没必要的重复获取.
func (srv *CorpAccessTokenServer) getToken() (token CorpAccessTokenInfo, cached bool, err error) {
	srv.tokenGet.Lock()
	defer srv.tokenGet.Unlock()

	timeNowUnix := time.Now().Unix()

	// 在收敛周期内直接返回最近一次获取的 corp_access_token, 这里的收敛时间设定为4秒.
	if n := srv.tokenGet.LastTimestamp; n <= timeNowUnix && timeNowUnix < n+4 {
		// 因为只有成功获取后才会更新 srv.tokenGet.LastTimestamp, 所以这些都是有效数据
		token = CorpAccessTokenInfo{
			Token:     srv.tokenGet.LastTokenInfo.Token,
			ExpiresIn: srv.tokenGet.LastTokenInfo.ExpiresIn - timeNowUnix + n,
		}
		cached = true
		return
	}

	request := struct {
		SuiteId       string `json:"suite_id"`
		AuthCorpId    string `json:"auth_corpid"`
		PermanentCode string `json:"permanent_code"`
	}{
		SuiteId:       srv.client.SuiteId,
		AuthCorpId:    srv.authCorpId,
		PermanentCode: srv.permanentCode,
	}

	var result struct {
		corp.Error
		CorpAccessTokenInfo
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/service/get_corp_token?suite_access_token="
	if err = srv.client.PostJSON(incompleteURL, &request, &result); err != nil {
		srv.tokenCache.Lock()
		srv.tokenCache.Token = ""
		srv.tokenCache.Unlock()
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		srv.tokenCache.Lock()
		srv.tokenCache.Token = ""
		srv.tokenCache.Unlock()

		err = &result.Error
		return
	}

	// 由于网络的延时, corp_access_token 过期时间留了一个缓冲区
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
	srv.tokenGet.LastTokenInfo = result.CorpAccessTokenInfo
	srv.tokenGet.LastTimestamp = timeNowUnix

	// 更新缓存
	srv.tokenCache.Lock()
	srv.tokenCache.Token = result.CorpAccessTokenInfo.Token
	srv.tokenCache.Unlock()

	token = result.CorpAccessTokenInfo
	return
}
