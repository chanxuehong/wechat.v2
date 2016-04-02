// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package template

import (
	"github.com/chanxuehong/wechat/mp"
)

const (
	EventTypeTemplateSendJobFinish = "TEMPLATESENDJOBFINISH"
)

const (
	TemplateSendStatusSuccess            = "success"               // 送达成功时
	TemplateSendStatusFailedUserBlock    = "failed:user block"     // 送达由于用户拒收(用户设置拒绝接收公众号消息)而失败
	TemplateSendStatusFailedSystemFailed = "failed: system failed" // 送达由于其他原因失败
)

// 在模版消息发送任务完成后, 微信服务器会将是否送达成功作为通知, 发送到开发者中心中填写的服务器配置地址中.
type TemplateSendJobFinishEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	Event  string `xml:"Event"  json:"Event"` // 事件信息, 此处为 TEMPLATESENDJOBFINISH
	MsgId  int64  `xml:"MsgId"  json:"MsgId"` // 模板消息ID
	Status string `xml:"Status" json:"Status"`
}

func GetTemplateSendJobFinishEvent(msg *mp.MixedMessage) *TemplateSendJobFinishEvent {
	return &TemplateSendJobFinishEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		MsgId:         msg.MsgID, // NOTE
		Status:        msg.Status,
	}
}
