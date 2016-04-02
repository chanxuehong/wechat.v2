// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mass

import (
	"github.com/chanxuehong/wechat/mp"
)

const (
	EventTypeMassSendJobFinish = "MASSSENDJOBFINISH"
)

// 高级群发消息, 事件推送群发结果
type MassSendJobFinishEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	Event string `xml:"Event" json:"Event"` // 事件信息, 此处为 MASSSENDJOBFINISH

	MsgId int64 `xml:"MsgId" json:"MsgId"` // 群发的消息ID, 64位整型

	// 群发的结构, 为 "send success" 或 "send fail" 或 "err(num)".
	// 但 send success 时, 也有可能因用户拒收公众号的消息, 系统错误等原因造成少量用户接收失败.
	// err(num) 是审核失败的具体原因, 可能的情况如下:
	// err(10001), //涉嫌广告
	// err(20001), //涉嫌政治
	// err(20004), //涉嫌社会
	// err(20002), //涉嫌色情
	// err(20006), //涉嫌违法犯罪
	// err(20008), //涉嫌欺诈
	// err(20013), //涉嫌版权
	// err(22000), //涉嫌互推(互相宣传)
	// err(21000), //涉嫌其他
	Status string `xml:"Status" json:"Status"`

	TotalCount int `xml:"TotalCount" json:"TotalCount"` // group_id 下粉丝数, 或者 openid_list 中的粉丝数

	// 过滤(过滤是指特定地区, 性别的过滤, 用户设置拒收的过滤; 用户接收已超4条的过滤)后,
	// 准备发送的粉丝数, 原则上, FilterCount = SentCount + ErrorCount
	FilterCount int `xml:"FilterCount" json:"FilterCount"`
	SentCount   int `xml:"SentCount"   json:"SentCount"`  // 发送成功的粉丝数
	ErrorCount  int `xml:"ErrorCount"  json:"ErrorCount"` // 发送失败的粉丝数
}

func GetMassSendJobFinishEvent(msg *mp.MixedMessage) *MassSendJobFinishEvent {
	return &MassSendJobFinishEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		MsgId:         msg.MsgID, // NOTE
		Status:        msg.Status,
		TotalCount:    msg.TotalCount,
		FilterCount:   msg.FilterCount,
		SentCount:     msg.SentCount,
		ErrorCount:    msg.ErrorCount,
	}
}
