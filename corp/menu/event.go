// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package menu

import (
	"github.com/chanxuehong/wechat/corp"
)

const (
	EventTypeClick = "click" // 点击菜单拉取消息时的事件推送
	EventTypeView  = "view"  // 点击菜单跳转链接时的事件推送

	// 请注意, 下面的事件仅支持微信iPhone5.4.1以上版本, 和Android5.4以上版本的微信用户,
	// 旧版本微信用户点击后将没有回应, 开发者也不能正常接收到事件推送.
	EventTypeScanCodePush    = "scancode_push"      // scancode_push: 扫码推事件的事件推送
	EventTypeScanCodeWaitMsg = "scancode_waitmsg"   // scancode_waitmsg: 扫码推事件且弹出"消息接收中"提示框的事件推送
	EventTypePicSysPhoto     = "pic_sysphoto"       // pic_sysphoto: 弹出系统拍照发图的事件推送
	EventTypePicPhotoOrAlbum = "pic_photo_or_album" // pic_photo_or_album: 弹出拍照或者相册发图的事件推送
	EventTypePicWeixin       = "pic_weixin"         // pic_weixin: 弹出微信相册发图器的事件推送
	EventTypeLocationSelect  = "location_select"    // location_select: 弹出地理位置选择器的事件推送
	EventTypeEnterAgent      = "enter_agent"        // enter_agent: 用户进入应用的事件推送
)

// 点击菜单拉取消息的事件推送
type ClickEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, CLICK
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 与自定义菜单接口中KEY值对应
}

func GetClickEvent(msg *corp.MixedMessage) *ClickEvent {
	return &ClickEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		EventKey:      msg.EventKey,
	}
}

// 点击菜单跳转链接的事件推送
type ViewEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, VIEW
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 设置的跳转URL
}

func GetViewEvent(msg *corp.MixedMessage) *ViewEvent {
	return &ViewEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		EventKey:      msg.EventKey,
	}
}

// scancode_push: 扫码推事件的事件推送
type ScanCodePushEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, scancode_push
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 由开发者在创建菜单时设定

	ScanCodeInfo struct {
		ScanType   string `xml:"ScanType"   json:"ScanType"`   // 扫描类型, 一般是qrcode
		ScanResult string `xml:"ScanResult" json:"ScanResult"` // 扫描结果, 即二维码对应的字符串信息
	} `xml:"ScanCodeInfo" json:"ScanCodeInfo"` // 扫描信息
}

func GetScanCodePushEvent(msg *corp.MixedMessage) *ScanCodePushEvent {
	return &ScanCodePushEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		EventKey:      msg.EventKey,
		ScanCodeInfo:  msg.ScanCodeInfo,
	}
}

// scancode_waitmsg: 扫码推事件且弹出"消息接收中"提示框的事件推送
type ScanCodeWaitMsgEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, scancode_waitmsg
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 由开发者在创建菜单时设定

	ScanCodeInfo struct {
		ScanType   string `xml:"ScanType"   json:"ScanType"`   // 扫描类型, 一般是qrcode
		ScanResult string `xml:"ScanResult" json:"ScanResult"` // 扫描结果, 即二维码对应的字符串信息
	} `xml:"ScanCodeInfo" json:"ScanCodeInfo"` // 扫描信息
}

func GetScanCodeWaitMsgEvent(msg *corp.MixedMessage) *ScanCodeWaitMsgEvent {
	return &ScanCodeWaitMsgEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		EventKey:      msg.EventKey,
		ScanCodeInfo:  msg.ScanCodeInfo,
	}
}

// pic_sysphoto: 弹出系统拍照发图的事件推送
type PicSysPhotoEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, pic_sysphoto
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 由开发者在创建菜单时设定

	SendPicsInfo struct {
		Count   int `xml:"Count" json:"Count"` // 发送的图片数量
		PicList []struct {
			PicMD5Sum string `xml:"PicMd5Sum" json:"PicMd5Sum"` // 图片的MD5值, 开发者若需要, 可用于验证接收到图片
		} `xml:"PicList>item,omitempty" json:"PicList,omitempty"` // 图片列表
	} `xml:"SendPicsInfo" json:"SendPicsInfo"` // 发送的图片信息
}

func GetPicSysPhotoEvent(msg *corp.MixedMessage) *PicSysPhotoEvent {
	return &PicSysPhotoEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		EventKey:      msg.EventKey,
		SendPicsInfo:  msg.SendPicsInfo,
	}
}

// pic_photo_or_album: 弹出拍照或者相册发图的事件推送
type PicPhotoOrAlbumEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, pic_photo_or_album
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 由开发者在创建菜单时设定

	SendPicsInfo struct {
		Count   int `xml:"Count" json:"Count"` // 发送的图片数量
		PicList []struct {
			PicMD5Sum string `xml:"PicMd5Sum" json:"PicMd5Sum"` // 图片的MD5值, 开发者若需要, 可用于验证接收到图片
		} `xml:"PicList>item,omitempty" json:"PicList,omitempty"` // 图片列表
	} `xml:"SendPicsInfo" json:"SendPicsInfo"` // 发送的图片信息
}

func GetPicPhotoOrAlbumEvent(msg *corp.MixedMessage) *PicPhotoOrAlbumEvent {
	return &PicPhotoOrAlbumEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		EventKey:      msg.EventKey,
		SendPicsInfo:  msg.SendPicsInfo,
	}
}

// pic_weixin: 弹出微信相册发图器的事件推送
type PicWeixinEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, pic_weixin
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 由开发者在创建菜单时设定

	SendPicsInfo struct {
		Count   int `xml:"Count" json:"Count"` // 发送的图片数量
		PicList []struct {
			PicMD5Sum string `xml:"PicMd5Sum" json:"PicMd5Sum"` // 图片的MD5值, 开发者若需要, 可用于验证接收到图片
		} `xml:"PicList>item,omitempty" json:"PicList,omitempty"` // 图片列表
	} `xml:"SendPicsInfo" json:"SendPicsInfo"` // 发送的图片信息
}

func GetPicWeixinEvent(msg *corp.MixedMessage) *PicWeixinEvent {
	return &PicWeixinEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		EventKey:      msg.EventKey,
		SendPicsInfo:  msg.SendPicsInfo,
	}
}

// location_select: 弹出地理位置选择器的事件推送
type LocationSelectEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, location_select
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 由开发者在创建菜单时设定

	SendLocationInfo struct {
		LocationX float64 `xml:"Location_X" json:"Location_X"` // 地理位置纬度
		LocationY float64 `xml:"Location_Y" json:"Location_Y"` // 地理位置经度
		Scale     int     `xml:"Scale"      json:"Scale"`      // 精度, 可理解为精度或者比例尺, 越精细的话 scale越高
		Label     string  `xml:"Label"      json:"Label"`      // 地理位置的字符串信息
		PoiName   string  `xml:"Poiname"    json:"Poiname"`    // 朋友圈POI的名字, 可能为空
	} `xml:"SendLocationInfo" json:"SendLocationInfo"` // 发送的位置信息
}

func GetLocationSelectEvent(msg *corp.MixedMessage) *LocationSelectEvent {
	return &LocationSelectEvent{
		MessageHeader:    msg.MessageHeader,
		Event:            msg.Event,
		EventKey:         msg.EventKey,
		SendLocationInfo: msg.SendLocationInfo,
	}
}

// enter_agent: 用户进入应用的事件推送
//  本事件只有在应用的回调模式中打开上报开关时上报
type EnterAgentEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, enter_agent
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 此事件该值为空
}

func GetEnterAgentEvent(msg *corp.MixedMessage) *EnterAgentEvent {
	return &EnterAgentEvent{
		MessageHeader: msg.MessageHeader,
		Event:         msg.Event,
		EventKey:      msg.EventKey,
	}
}
