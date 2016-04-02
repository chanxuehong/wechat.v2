// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mp

import (
	"net/http"
	"net/url"
)

// ServerFrontend 实现了 http.Handler, 处理一个公众号的消息(事件)请求.
type ServerFrontend struct {
	server      Server
	errHandler  ErrorHandler
	interceptor Interceptor
}

// NOTE: errHandler, interceptor 均可以为 nil
func NewServerFrontend(server Server, errHandler ErrorHandler, interceptor Interceptor) *ServerFrontend {
	if server == nil {
		panic("nil Server")
	}
	if errHandler == nil {
		errHandler = DefaultErrorHandler
	}

	return &ServerFrontend{
		server:      server,
		errHandler:  errHandler,
		interceptor: interceptor,
	}
}

func (frontend *ServerFrontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		frontend.errHandler.ServeError(w, r, err)
		return
	}

	if interceptor := frontend.interceptor; interceptor != nil && !interceptor.Intercept(w, r, queryValues) {
		return
	}

	ServeHTTP(w, r, queryValues, frontend.server, frontend.errHandler)
}
