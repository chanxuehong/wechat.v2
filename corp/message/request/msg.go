// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package request

import (
	"github.com/chanxuehong/wechat/corp"
)

const (
	// 微信服务器推送过来的消息类型
	MsgTypeText     = "text"     // 文本消息
	MsgTypeImage    = "image"    // 图片消息
	MsgTypeVoice    = "voice"    // 语音消息
	MsgTypeVideo    = "video"    // 视频消息
	MsgTypeLocation = "location" // 地理位置消息
)

type Text struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	MsgId   int64  `xml:"MsgId"   json:"MsgId"`   // 消息id, 64位整型
	Content string `xml:"Content" json:"Content"` // 文本消息内容
}

func GetText(msg *corp.MixedMessage) *Text {
	return &Text{
		MessageHeader: msg.MessageHeader,
		MsgId:         msg.MsgId,
		Content:       msg.Content,
	}
}

type Image struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	MsgId   int64  `xml:"MsgId"   json:"MsgId"`   // 消息id, 64位整型
	MediaId string `xml:"MediaId" json:"MediaId"` // 图片媒体文件id, 可以调用获取媒体文件接口拉取数据
	PicURL  string `xml:"PicUrl"  json:"PicUrl"`  // 图片链接
}

func GetImage(msg *corp.MixedMessage) *Image {
	return &Image{
		MessageHeader: msg.MessageHeader,
		MsgId:         msg.MsgId,
		MediaId:       msg.MediaId,
		PicURL:        msg.PicURL,
	}
}

type Voice struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	MsgId   int64  `xml:"MsgId"   json:"MsgId"`   // 消息id, 64位整型
	MediaId string `xml:"MediaId" json:"MediaId"` // 语音媒体文件id, 可以调用获取媒体文件接口拉取数据
	Format  string `xml:"Format"  json:"Format"`  // 语音格式, 如amr, speex等
}

func GetVoice(msg *corp.MixedMessage) *Voice {
	return &Voice{
		MessageHeader: msg.MessageHeader,
		MsgId:         msg.MsgId,
		MediaId:       msg.MediaId,
		Format:        msg.Format,
	}
}

type Video struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	MsgId        int64  `xml:"MsgId"        json:"MsgId"`        // 消息id, 64位整型
	MediaId      string `xml:"MediaId"      json:"MediaId"`      // 视频媒体文件id, 可以调用获取媒体文件接口拉取数据
	ThumbMediaId string `xml:"ThumbMediaId" json:"ThumbMediaId"` // 视频消息缩略图的媒体id, 可以调用获取媒体文件接口拉取数据
}

func GetVideo(msg *corp.MixedMessage) *Video {
	return &Video{
		MessageHeader: msg.MessageHeader,
		MsgId:         msg.MsgId,
		MediaId:       msg.MediaId,
		ThumbMediaId:  msg.ThumbMediaId,
	}
}

type Location struct {
	XMLName struct{} `xml:"xml" json:"-"`
	corp.MessageHeader

	MsgId     int64   `xml:"MsgId"      json:"MsgId"`      // 消息id, 64位整型
	LocationX float64 `xml:"Location_X" json:"Location_X"` // 地理位置纬度
	LocationY float64 `xml:"Location_Y" json:"Location_Y"` // 地理位置经度
	Scale     int     `xml:"Scale"      json:"Scale"`      // 地图缩放大小
	Label     string  `xml:"Label"      json:"Label"`      // 地理位置信息
}

func GetLocation(msg *corp.MixedMessage) *Location {
	return &Location{
		MessageHeader: msg.MessageHeader,
		MsgId:         msg.MsgId,
		LocationX:     msg.LocationX,
		LocationY:     msg.LocationY,
		Scale:         msg.Scale,
		Label:         msg.Label,
	}
}
