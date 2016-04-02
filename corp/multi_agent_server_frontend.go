// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package corp

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

// 多个 AgentServer 的前端, http.Handler 的实现.
//
//  MultiAgentServerFrontend 可以处理多个企业号应用的消息(事件), 但是要求在回调 URL 上加上一个
//  查询参数(参数名与 urlAgentServerQueryName 一致), 通过这个参数的值来索引对应的 AgentServer.
//
//  例如回调 URL 为(urlAgentServerQueryName == "agent_server"):
//    http://www.xxx.com/weixin?agent_server=1234567890
//  那么就可以在后端调用
//   MultiAgentServerFrontend.SetAgentServer("1234567890", AgentServer)
//  来增加一个 AgentServer 来处理 agent_server=1234567890 的消息(事件).
//
//  MultiAgentServerFrontend 并发安全, 可以在运行中动态增加和删除 AgentServer.
type MultiAgentServerFrontend struct {
	urlAgentServerQueryName string

	errHandler  ErrorHandler
	interceptor Interceptor

	rwmutex        sync.RWMutex
	agentServerMap map[string]AgentServer
}

// NewMultiAgentServerFrontend 创建一个新的 MultiAgentServerFrontend.
//  urlAgentServerQueryName: 回调 URL 上参数名, 这个参数的值就是索引 AgentServer 的 key
//  errHandler:              错误处理 handler, 可以为 nil
//  interceptor:             拦截器, 可以为 nil
func NewMultiAgentServerFrontend(urlAgentServerQueryName string, errHandler ErrorHandler, interceptor Interceptor) *MultiAgentServerFrontend {
	if urlAgentServerQueryName == "" {
		urlAgentServerQueryName = "agent_server"
	}
	if errHandler == nil {
		errHandler = DefaultErrorHandler
	}

	return &MultiAgentServerFrontend{
		urlAgentServerQueryName: urlAgentServerQueryName,
		errHandler:              errHandler,
		interceptor:             interceptor,
		agentServerMap:          make(map[string]AgentServer),
	}
}

func (frontend *MultiAgentServerFrontend) SetAgentServer(serverKey string, server AgentServer) (err error) {
	if serverKey == "" {
		return errors.New("empty serverKey")
	}
	if server == nil {
		return errors.New("nil AgentServer")
	}

	frontend.rwmutex.Lock()
	frontend.agentServerMap[serverKey] = server
	frontend.rwmutex.Unlock()
	return
}

func (frontend *MultiAgentServerFrontend) DeleteAgentServer(serverKey string) {
	frontend.rwmutex.Lock()
	delete(frontend.agentServerMap, serverKey)
	frontend.rwmutex.Unlock()
}

func (frontend *MultiAgentServerFrontend) DeleteAllAgentServer() {
	frontend.rwmutex.Lock()
	frontend.agentServerMap = make(map[string]AgentServer)
	frontend.rwmutex.Unlock()
}

func (frontend *MultiAgentServerFrontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		frontend.errHandler.ServeError(w, r, err)
		return
	}

	if interceptor := frontend.interceptor; interceptor != nil && !interceptor.Intercept(w, r, queryValues) {
		return
	}

	serverKey := queryValues.Get(frontend.urlAgentServerQueryName)
	if serverKey == "" {
		err := fmt.Errorf("the url query value with name %s is empty", frontend.urlAgentServerQueryName)
		frontend.errHandler.ServeError(w, r, err)
		return
	}

	frontend.rwmutex.RLock()
	agentServer := frontend.agentServerMap[serverKey]
	frontend.rwmutex.RUnlock()

	if agentServer == nil {
		err := fmt.Errorf("Not found AgentServer for %s == %s", frontend.urlAgentServerQueryName, serverKey)
		frontend.errHandler.ServeError(w, r, err)
		return
	}

	ServeHTTP(w, r, queryValues, agentServer, frontend.errHandler)
}
