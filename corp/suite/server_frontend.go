// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package suite

import (
	"net/http"
	"net/url"

	"github.com/chanxuehong/wechat/corp"
)

type ServerFrontend struct {
	server      Server
	errHandler  corp.ErrorHandler
	interceptor corp.Interceptor
}

// handler, interceptor 均可以为 nil
func NewServerFrontend(server Server, handler corp.ErrorHandler, interceptor corp.Interceptor) *ServerFrontend {
	if server == nil {
		panic("nil Server")
	}
	if handler == nil {
		handler = corp.DefaultErrorHandler
	}

	return &ServerFrontend{
		server:      server,
		errHandler:  handler,
		interceptor: interceptor,
	}
}

// 实现 http.Handler.
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
