// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mch

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

// 多个 Server 的前端, http.Handler 的实现.
//
//  MultiServerFrontend 可以处理多个APP的消息(事件), 但是要求在回调 URL 上加上一个
//  查询参数(参数名与 urlServerQueryName 一致), 通过这个参数的值来索引对应的 Server.
//
//  例如回调 URL 为(urlServerQueryName == "mch_server"):
//    http://www.xxx.com/weixin?mch_server=1234567890
//  那么就可以在后端调用
//   MultiServerFrontend.SetServer("1234567890", Server)
//  来增加一个 Server 来处理 mch_server=1234567890 的消息(事件).
//
//  MultiServerFrontend 并发安全, 可以在运行中动态增加和删除 Server.
type MultiServerFrontend struct {
	urlServerQueryName string

	errHandler  ErrorHandler
	interceptor Interceptor

	rwmutex   sync.RWMutex
	serverMap map[string]Server
}

// NewMultiServerFrontend 创建一个新的 MultiServerFrontend.
//  urlServerQueryName: 回调 URL 上参数名, 这个参数的值就是索引 Server 的 key
//  errHandler:         错误处理 handler, 可以为 nil
//  interceptor:        拦截器, 可以为 nil
func NewMultiServerFrontend(urlServerQueryName string, errHandler ErrorHandler, interceptor Interceptor) *MultiServerFrontend {
	if urlServerQueryName == "" {
		urlServerQueryName = "mch_server"
	}
	if errHandler == nil {
		errHandler = DefaultErrorHandler
	}

	return &MultiServerFrontend{
		urlServerQueryName: urlServerQueryName,
		errHandler:         errHandler,
		interceptor:        interceptor,
		serverMap:          make(map[string]Server),
	}
}

func (frontend *MultiServerFrontend) SetServer(serverKey string, server Server) (err error) {
	if serverKey == "" {
		return errors.New("empty serverKey")
	}
	if server == nil {
		return errors.New("nil Server")
	}

	frontend.rwmutex.Lock()
	frontend.serverMap[serverKey] = server
	frontend.rwmutex.Unlock()
	return
}

func (frontend *MultiServerFrontend) DeleteServer(serverKey string) {
	frontend.rwmutex.Lock()
	delete(frontend.serverMap, serverKey)
	frontend.rwmutex.Unlock()
}

func (frontend *MultiServerFrontend) DeleteAllServer() {
	frontend.rwmutex.Lock()
	frontend.serverMap = make(map[string]Server)
	frontend.rwmutex.Unlock()
}

func (frontend *MultiServerFrontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		frontend.errHandler.ServeError(w, r, err)
		return
	}

	if interceptor := frontend.interceptor; interceptor != nil && !interceptor.Intercept(w, r, queryValues) {
		return
	}

	serverKey := queryValues.Get(frontend.urlServerQueryName)
	if serverKey == "" {
		err := fmt.Errorf("the url query value with name %s is empty", frontend.urlServerQueryName)
		frontend.errHandler.ServeError(w, r, err)
		return
	}

	frontend.rwmutex.RLock()
	server := frontend.serverMap[serverKey]
	frontend.rwmutex.RUnlock()

	if server == nil {
		err := fmt.Errorf("Not found Server for %s == %s", frontend.urlServerQueryName, serverKey)
		frontend.errHandler.ServeError(w, r, err)
		return
	}

	ServeHTTP(w, r, queryValues, server, frontend.errHandler)
}
