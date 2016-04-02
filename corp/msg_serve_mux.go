// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package corp

import (
	"net/http"
	"sync"

	"github.com/chanxuehong/wechat/internal/util"
)

var _ MessageHandler = (*MessageServeMux)(nil)

// MessageServeMux 实现了一个简单的消息(事件)路由器, 同时也是一个 MessageHandler 的实现.
type MessageServeMux struct {
	rwmutex               sync.RWMutex
	messageHandlerMap     map[string]MessageHandler // map[MsgType]MessageHandler
	eventHandlerMap       map[string]MessageHandler // map[EventType]MessageHandler
	defaultMessageHandler MessageHandler
	defaultEventHandler   MessageHandler
}

func NewMessageServeMux() *MessageServeMux {
	return &MessageServeMux{
		messageHandlerMap: make(map[string]MessageHandler),
		eventHandlerMap:   make(map[string]MessageHandler),
	}
}

// 注册特定类型消息的 MessageHandler.
func (mux *MessageServeMux) MessageHandle(msgType string, handler MessageHandler) {
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
func (mux *MessageServeMux) MessageHandleFunc(msgType string, handler func(http.ResponseWriter, *Request)) {
	mux.MessageHandle(msgType, MessageHandlerFunc(handler))
}

// 注册消息的默认 MessageHandler.
func (mux *MessageServeMux) DefaultMessageHandle(handler MessageHandler) {
	if handler == nil {
		panic("nil MessageHandler")
	}

	mux.rwmutex.Lock()
	mux.defaultMessageHandler = handler
	mux.rwmutex.Unlock()
}

// 注册消息的默认 MessageHandler.
func (mux *MessageServeMux) DefaultMessageHandleFunc(handler func(http.ResponseWriter, *Request)) {
	mux.DefaultMessageHandle(MessageHandlerFunc(handler))
}

// 注册特定类型事件的 MessageHandler.
func (mux *MessageServeMux) EventHandle(eventType string, handler MessageHandler) {
	if eventType == "" {
		panic("empty eventType")
	}
	if handler == nil {
		panic("nil MessageHandler")
	}

	mux.rwmutex.Lock()
	if mux.eventHandlerMap == nil {
		mux.eventHandlerMap = make(map[string]MessageHandler)
	}
	mux.eventHandlerMap[util.ToLower(eventType)] = handler
	mux.rwmutex.Unlock()
}

// 注册特定类型事件的 MessageHandler.
func (mux *MessageServeMux) EventHandleFunc(eventType string, handler func(http.ResponseWriter, *Request)) {
	mux.EventHandle(eventType, MessageHandlerFunc(handler))
}

// 注册事件的默认 MessageHandler.
func (mux *MessageServeMux) DefaultEventHandle(handler MessageHandler) {
	if handler == nil {
		panic("nil MessageHandler")
	}

	mux.rwmutex.Lock()
	mux.defaultEventHandler = handler
	mux.rwmutex.Unlock()
}

// 注册事件的默认 MessageHandler.
func (mux *MessageServeMux) DefaultEventHandleFunc(handler func(http.ResponseWriter, *Request)) {
	mux.DefaultEventHandle(MessageHandlerFunc(handler))
}

// 获取 msgType 对应的 MessageHandler, 如果没有找到返回 nil.
func (mux *MessageServeMux) getMessageHandler(msgType string) (handler MessageHandler) {
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

// 获取 eventType 对应的 MessageHandler, 如果没有找到返回 nil.
func (mux *MessageServeMux) getEventHandler(eventType string) (handler MessageHandler) {
	if eventType == "" {
		return nil
	}

	mux.rwmutex.RLock()
	handler = mux.eventHandlerMap[util.ToLower(eventType)]
	if handler == nil {
		handler = mux.defaultEventHandler
	}
	mux.rwmutex.RUnlock()
	return
}

// MessageServeMux 实现了 MessageHandler 接口.
func (mux *MessageServeMux) ServeMessage(w http.ResponseWriter, r *Request) {
	if msgType := r.MixedMsg.MsgType; msgType == "event" {
		handler := mux.getEventHandler(r.MixedMsg.Event)
		if handler == nil {
			return // 返回空串, 符合微信协议
		}
		handler.ServeMessage(w, r)
	} else {
		handler := mux.getMessageHandler(msgType)
		if handler == nil {
			return // 返回空串, 符合微信协议
		}
		handler.ServeMessage(w, r)
	}
}
