// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package component

const (
	// 微信服务器推送过来的消息类型
	MsgTypeVerifyTicket = "component_verify_ticket" // 推送 component_verify_ticket 协议
	MsgTypeUnauthorized = "unauthorized"            // 取消授权的通知
)

type VerifyTicketMessage struct {
	XMLName struct{} `xml:"xml" json:"-"`

	AppId      string `xml:"AppId"      json:"AppId"`
	CreateTime int64  `xml:"CreateTime" json:"CreateTime"`
	InfoType   string `xml:"InfoType"   json:"InfoType"`

	VerifyTicket string `xml:"ComponentVerifyTicket" json:"ComponentVerifyTicket"`
}

func GetVerifyTicketMessage(msg *MixedMessage) *VerifyTicketMessage {
	return &VerifyTicketMessage{
		AppId:        msg.AppId,
		CreateTime:   msg.CreateTime,
		InfoType:     msg.InfoType,
		VerifyTicket: msg.VerifyTicket,
	}
}

type UnauthorizedMessage struct {
	XMLName struct{} `xml:"xml" json:"-"`

	AppId      string `xml:"AppId"      json:"AppId"`
	CreateTime int64  `xml:"CreateTime" json:"CreateTime"`
	InfoType   string `xml:"InfoType"   json:"InfoType"`

	AuthorizerAppId string `xml:"AuthorizerAppid" json:"AuthorizerAppid"`
}

func GetUnauthorizedMessage(msg *MixedMessage) *UnauthorizedMessage {
	return &UnauthorizedMessage{
		AppId:           msg.AppId,
		CreateTime:      msg.CreateTime,
		InfoType:        msg.InfoType,
		AuthorizerAppId: msg.AuthorizerAppId,
	}
}
