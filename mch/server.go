// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mch

type Server interface {
	AppId() string  // AppId, 用于约束消息的 appid, 如果为空表示不约束
	MchId() string  // MchId, 用于约束消息的 mch_id, 如果为空表示不约束
	APIKey() string // API密钥

	MessageHandler() MessageHandler // 获取 MessageHandler
}

var _ Server = (*DefaultServer)(nil)

type DefaultServer struct {
	appId  string
	mchId  string
	apiKey string

	messageHandler MessageHandler
}

func NewDefaultServer(appId, mchId, apiKey string, handler MessageHandler) *DefaultServer {
	if handler == nil {
		panic("nil MessageHandler")
	}

	return &DefaultServer{
		appId:          appId,
		mchId:          mchId,
		apiKey:         apiKey,
		messageHandler: handler,
	}
}

func (srv *DefaultServer) AppId() string {
	return srv.appId
}
func (srv *DefaultServer) MchId() string {
	return srv.mchId
}
func (srv *DefaultServer) APIKey() string {
	return srv.apiKey
}
func (srv *DefaultServer) MessageHandler() MessageHandler {
	return srv.messageHandler
}
