// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package request

import (
	"github.com/chanxuehong/wechat/mp"
)

const (
	// 微信服务器推送过来的消息类型
	MsgTypeText       = "text"       // 文本消息
	MsgTypeImage      = "image"      // 图片消息
	MsgTypeVoice      = "voice"      // 语音消息
	MsgTypeVideo      = "video"      // 视频消息
	MsgTypeShortVideo = "shortvideo" // 小视频消息
	MsgTypeLocation   = "location"   // 地理位置消息
	MsgTypeLink       = "link"       // 链接消息
)

// 文本消息
type Text struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	MsgId   int64  `xml:"MsgId"   json:"MsgId"`   // 消息id, 64位整型
	Content string `xml:"Content" json:"Content"` // 文本消息内容
}

func GetText(msg *mp.MixedMessage) *Text {
	return &Text{
		MessageHeader: msg.MessageHeader,
		MsgId:         msg.MsgId,
		Content:       msg.Content,
	}
}

// 图片消息
type Image struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	MsgId   int64  `xml:"MsgId"   json:"MsgId"`   // 消息id, 64位整型
	MediaId string `xml:"MediaId" json:"MediaId"` // 图片消息媒体id, 可以调用多媒体文件下载接口拉取数据.
	PicURL  string `xml:"PicUrl"  json:"PicUrl"`  // 图片链接
}

func GetImage(msg *mp.MixedMessage) *Image {
	return &Image{
		MessageHeader: msg.MessageHeader,
		MsgId:         msg.MsgId,
		MediaId:       msg.MediaId,
		PicURL:        msg.PicURL,
	}
}

// 语音消息
type Voice struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	MsgId   int64  `xml:"MsgId"   json:"MsgId"`   // 消息id, 64位整型
	MediaId string `xml:"MediaId" json:"MediaId"` // 语音消息媒体id, 可以调用多媒体文件下载接口拉取该媒体
	Format  string `xml:"Format"  json:"Format"`  // 语音格式, 如amr, speex等

	// 语音识别结果, UTF8编码,
	// NOTE: 需要开通语音识别功能, 否则该字段为空, 即使开通了语音识别该字段还是有可能为空
	Recognition string `xml:"Recognition,omitempty" json:"Recognition,omitempty"`
}

func GetVoice(msg *mp.MixedMessage) *Voice {
	return &Voice{
		MessageHeader: msg.MessageHeader,
		MsgId:         msg.MsgId,
		MediaId:       msg.MediaId,
		Format:        msg.Format,
		Recognition:   msg.Recognition,
	}
}

// 视频消息
type Video struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	MsgId        int64  `xml:"MsgId"        json:"MsgId"`        // 消息id, 64位整型
	MediaId      string `xml:"MediaId"      json:"MediaId"`      // 视频消息媒体id, 可以调用多媒体文件下载接口拉取数据.
	ThumbMediaId string `xml:"ThumbMediaId" json:"ThumbMediaId"` // 视频消息缩略图的媒体id, 可以调用多媒体文件下载接口拉取数据.
}

func GetVideo(msg *mp.MixedMessage) *Video {
	return &Video{
		MessageHeader: msg.MessageHeader,
		MsgId:         msg.MsgId,
		MediaId:       msg.MediaId,
		ThumbMediaId:  msg.ThumbMediaId,
	}
}

// 小视频消息
type ShortVideo struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	MsgId        int64  `xml:"MsgId"        json:"MsgId"`        // 消息id, 64位整型
	MediaId      string `xml:"MediaId"      json:"MediaId"`      // 视频消息媒体id, 可以调用多媒体文件下载接口拉取数据.
	ThumbMediaId string `xml:"ThumbMediaId" json:"ThumbMediaId"` // 视频消息缩略图的媒体id, 可以调用多媒体文件下载接口拉取数据.
}

func GetShortVideo(msg *mp.MixedMessage) *ShortVideo {
	return &ShortVideo{
		MessageHeader: msg.MessageHeader,
		MsgId:         msg.MsgId,
		MediaId:       msg.MediaId,
		ThumbMediaId:  msg.ThumbMediaId,
	}
}

// 地理位置消息
type Location struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	MsgId     int64   `xml:"MsgId"      json:"MsgId"`      // 消息id, 64位整型
	LocationX float64 `xml:"Location_X" json:"Location_X"` // 地理位置纬度
	LocationY float64 `xml:"Location_Y" json:"Location_Y"` // 地理位置经度
	Scale     int     `xml:"Scale"      json:"Scale"`      // 地图缩放大小
	Label     string  `xml:"Label"      json:"Label"`      // 地理位置信息
}

func GetLocation(msg *mp.MixedMessage) *Location {
	return &Location{
		MessageHeader: msg.MessageHeader,
		MsgId:         msg.MsgId,
		LocationX:     msg.LocationX,
		LocationY:     msg.LocationY,
		Scale:         msg.Scale,
		Label:         msg.Label,
	}
}

// 链接消息
type Link struct {
	XMLName struct{} `xml:"xml" json:"-"`
	mp.MessageHeader

	MsgId       int64  `xml:"MsgId"       json:"MsgId"`       // 消息id, 64位整型
	Title       string `xml:"Title"       json:"Title"`       // 消息标题
	Description string `xml:"Description" json:"Description"` // 消息描述
	URL         string `xml:"Url"         json:"Url"`         // 消息链接
}

func GetLink(msg *mp.MixedMessage) *Link {
	return &Link{
		MessageHeader: msg.MessageHeader,
		MsgId:         msg.MsgId,
		Title:         msg.Title,
		Description:   msg.Description,
		URL:           msg.URL,
	}
}
