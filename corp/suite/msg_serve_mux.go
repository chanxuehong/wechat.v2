// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package suite

import (
	"io"
	"net/http"
	"sync"

	"github.com/chanxuehong/wechat/internal/util"
)

var _ MessageHandler = (*SuiteMessageServeMux)(nil)

// SuiteMessageServeMux 实现了一个简单的消息(事件)路由器, 同时也是一个 MessageHandler 的实现.
type SuiteMessageServeMux struct {
	rwmutex               sync.RWMutex
	messageHandlerMap     map[string]MessageHandler // map[MsgType]MessageHandler
	defaultMessageHandler MessageHandler
}

func NewSuiteMessageServeMux() *SuiteMessageServeMux {
	return &SuiteMessageServeMux{
		messageHandlerMap: make(map[string]MessageHandler),
	}
}

// 注册特定类型消息的 MessageHandler.
func (mux *SuiteMessageServeMux) MessageHandle(msgType string, handler MessageHandler) {
	if msgType == "" {
		panic("empty msgType")
	}
	if handler == nil {
		panic("nil MessageHandler")
	}

	mux.rwmutex.Lock()
	if mux.messageHandlerMap == nil {
		mux.messageHandlerMap = make(map[string]MessageHandler)
	}
	mux.messageHandlerMap[util.ToLower(msgType)] = handler
	mux.rwmutex.Unlock()
}

// 注册特定类型消息的 MessageHandler.
func (mux *SuiteMessageServeMux) MessageHandleFunc(msgType string, handler func(http.ResponseWriter, *Request)) {
	mux.MessageHandle(msgType, MessageHandlerFunc(handler))
}

// 注册消息的默认 MessageHandler.
func (mux *SuiteMessageServeMux) DefaultMessageHandle(handler MessageHandler) {
	if handler == nil {
		panic("nil MessageHandler")
	}

	mux.rwmutex.Lock()
	mux.defaultMessageHandler = handler
	mux.rwmutex.Unlock()
}

// 注册消息的默认 MessageHandler.
func (mux *SuiteMessageServeMux) DefaultMessageHandleFunc(handler func(http.ResponseWriter, *Request)) {
	mux.DefaultMessageHandle(MessageHandlerFunc(handler))
}

// 获取 msgType 对应的 MessageHandler, 如果没有找到返回 nil.
func (mux *SuiteMessageServeMux) getMessageHandler(msgType string) (handler MessageHandler) {
	if msgType == "" {
		return nil
	}

	mux.rwmutex.RLock()
	handler = mux.messageHandlerMap[util.ToLower(msgType)]
	if handler == nil {
		handler = mux.defaultMessageHandler
	}
	mux.rwmutex.RUnlock()
	return
}

// SuiteMessageServeMux 实现了 MessageHandler 接口.
func (mux *SuiteMessageServeMux) ServeMessage(w http.ResponseWriter, r *Request) {
	handler := mux.getMessageHandler(r.MixedMsg.InfoType)
	if handler == nil {
		io.WriteString(w, "success")
		return
	}
	handler.ServeMessage(w, r)
}
