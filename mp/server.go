// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mp

import (
	"bytes"
	"errors"
	"sync"
)

type Server interface {
	Token() string                          // 获取公众号的 Token, 用于校验签名.
	CurrentAESKey() [32]byte                // 获取当前的 AES 加密 Key
	LastAESKey() (key [32]byte, valid bool) // 获取上一个 AES 加密 Key

	OriId() string // 公众号的原始ID, 用于约束消息的 ToUserName, 如果为空表示不约束
	AppId() string // AppId, 用于约束消息的 AppId, 如果为空表示不约束

	MessageHandler() MessageHandler // 获取 MessageHandler
}

var _ Server = (*DefaultServer)(nil)

type DefaultServer struct {
	oriId string
	appId string
	token string

	rwmutex           sync.RWMutex
	currentAESKey     [32]byte
	lastAESKey        [32]byte
	isLastAESKeyValid bool

	messageHandler MessageHandler
}

// NewDefaultServer 创建一个新的 DefaultServer.
//  如果是明文模式, 则 appId 可以为 "", aesKey 可以为 nil.
func NewDefaultServer(oriId, token, appId string, aesKey []byte, handler MessageHandler) (srv *DefaultServer) {
	if aesKey != nil && len(aesKey) != 32 {
		panic("the length of aesKey must equal to 32")
	}
	if handler == nil {
		panic("nil MessageHandler")
	}

	srv = &DefaultServer{
		oriId:          oriId,
		appId:          appId,
		token:          token,
		messageHandler: handler,
	}
	copy(srv.currentAESKey[:], aesKey)
	return
}

func (srv *DefaultServer) OriId() string {
	return srv.oriId
}
func (srv *DefaultServer) AppId() string {
	return srv.appId
}
func (srv *DefaultServer) Token() string {
	return srv.token
}
func (srv *DefaultServer) MessageHandler() MessageHandler {
	return srv.messageHandler
}
func (srv *DefaultServer) CurrentAESKey() (key [32]byte) {
	srv.rwmutex.RLock()
	key = srv.currentAESKey
	srv.rwmutex.RUnlock()
	return
}
func (srv *DefaultServer) LastAESKey() (key [32]byte, valid bool) {
	srv.rwmutex.RLock()
	key = srv.lastAESKey
	valid = srv.isLastAESKeyValid
	srv.rwmutex.RUnlock()
	return
}

func (srv *DefaultServer) UpdateAESKey(aesKey []byte) (err error) {
	if len(aesKey) != 32 {
		return errors.New("the length of aesKey must equal to 32")
	}

	srv.rwmutex.Lock()
	defer srv.rwmutex.Unlock()

	if bytes.Equal(aesKey, srv.currentAESKey[:]) {
		return
	}

	srv.isLastAESKeyValid = true
	srv.lastAESKey = srv.currentAESKey
	copy(srv.currentAESKey[:], aesKey)
	return
}
