// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package request

import (
	"fmt"
	"strings"

	"github.com/chanxuehong/wechat/mp"
)

const (
	// 微信服务器推送过来的事件类型
	EventTypeSubscribe   = "subscribe"   // 关注事件, 包括点击关注和扫描二维码(公众号二维码和公众号带参数二维码)关注
	EventTypeUnsubscribe = "unsubscribe" // 取消关注
	EventTypeScan        = "SCAN"        // 已经关注的用户扫描带参数二维码事件
	EventTypeLocation    = "LOCATION"    // 上报地理位置事件
)

// 关注
type SubscribeEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	Event string `xml:"Event" json:"Event"` // subscribe

	// 下面两个字段只有在扫描带参数二维码进行关注时才有值, 否则为空值!
	EventKey string `xml:"EventKey,omitempty" json:"EventKey,omitempty"` // 事件KEY值, 格式为: qrscene_二维码的参数值
	Ticket   string `xml:"Ticket,omitempty"   json:"Ticket,omitempty"`   // 二维码的ticket, 可用来换取二维码图片
}

func GetSubscribeEvent(msg *mp.MixedMessage) *SubscribeEvent {
	return &SubscribeEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		EventKey:      msg.EventKey,
		Ticket:        msg.Ticket,
	}
}

// 获取二维码参数
func (event *SubscribeEvent) Scene() (scene string, err error) {
	const prefix = "qrscene_"
	if !strings.HasPrefix(event.EventKey, prefix) {
		err = fmt.Errorf("EventKey 应该以 %s 为前缀: %s", prefix, event.EventKey)
		return
	}
	scene = event.EventKey[len(prefix):]
	return
}

// 取消关注
type UnsubscribeEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	Event    string `xml:"Event"              json:"Event"`              // unsubscribe(取消关注)
	EventKey string `xml:"EventKey,omitempty" json:"EventKey,omitempty"` // 事件KEY值, 为空值
}

func GetUnsubscribeEvent(msg *mp.MixedMessage) *UnsubscribeEvent {
	return &UnsubscribeEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
	}
}

// 用户已关注时, 扫描带参数二维码的事件
type ScanEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	Event    string `xml:"Event"    json:"Event"`    // SCAN
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 二维码的参数值(scene_id, scene_str)
	Ticket   string `xml:"Ticket"   json:"Ticket"`   // 二维码的ticket, 可用来换取二维码图片
}

func GetScanEvent(msg *mp.MixedMessage) *ScanEvent {
	return &ScanEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		EventKey:      msg.EventKey,
		Ticket:        msg.Ticket,
	}
}

// 上报地理位置事件
type LocationEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	Event     string  `xml:"Event"     json:"Event"`     // LOCATION
	Latitude  float64 `xml:"Latitude"  json:"Latitude"`  // 地理位置纬度
	Longitude float64 `xml:"Longitude" json:"Longitude"` // 地理位置经度
	Precision float64 `xml:"Precision" json:"Precision"` // 地理位置精度(实际上应该是整数, 但是微信推送过来是浮点数形式)
}

func GetLocationEvent(msg *mp.MixedMessage) *LocationEvent {
	return &LocationEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		Latitude:      msg.Latitude,
		Longitude:     msg.Longitude,
		Precision:     msg.Precision,
	}
}
