// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

// 被动响应消息.
package response

import (
	"errors"
	"fmt"
)

const (
	MsgTypeText  = "text"  // 文本消息
	MsgTypeImage = "image" // 图片消息
	MsgTypeVoice = "voice" // 语音消息
	MsgTypeVideo = "video" // 视频消息
	MsgTypeNews  = "news"  // 图文消息
)

type MessageHeader struct {
	ToUserName   string `xml:"ToUserName"   json:"ToUserName"`
	FromUserName string `xml:"FromUserName" json:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"   json:"CreateTime"`
	MsgType      string `xml:"MsgType"      json:"MsgType"`
}

type Text struct {
	XMLName struct{} `xml:"xml" json:"-"`
	MessageHeader

	Content string `xml:"Content" json:"Content"` // 文本消息内容
}

func NewText(to, from string, timestamp int64, content string) (text *Text) {
	return &Text{
		MessageHeader: MessageHeader{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   timestamp,
			MsgType:      MsgTypeText,
		},
		Content: content,
	}
}

type Image struct {
	XMLName struct{} `xml:"xml" json:"-"`
	MessageHeader

	Image struct {
		MediaId string `xml:"MediaId" json:"MediaId"` // 图片文件id, 可以调用上传媒体文件接口获取
	} `xml:"Image" json:"Image"`
}

func NewImage(to, from string, timestamp int64, mediaId string) (image *Image) {
	image = &Image{
		MessageHeader: MessageHeader{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   timestamp,
			MsgType:      MsgTypeImage,
		},
	}
	image.Image.MediaId = mediaId
	return
}

type Voice struct {
	XMLName struct{} `xml:"xml" json:"-"`
	MessageHeader

	Voice struct {
		MediaId string `xml:"MediaId" json:"MediaId"` // 语音文件id, 可以调用上传媒体文件接口获取
	} `xml:"Voice" json:"Voice"`
}

func NewVoice(to, from string, timestamp int64, mediaId string) (voice *Voice) {
	voice = &Voice{
		MessageHeader: MessageHeader{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   timestamp,
			MsgType:      MsgTypeVoice,
		},
	}
	voice.Voice.MediaId = mediaId
	return
}

type Video struct {
	XMLName struct{} `xml:"xml" json:"-"`
	MessageHeader

	Video struct {
		MediaId     string `xml:"MediaId"               json:"MediaId"`               // 视频文件id, 可以调用上传媒体文件接口获取
		Title       string `xml:"Title,omitempty"       json:"Title,omitempty"`       // 视频消息的标题
		Description string `xml:"Description,omitempty" json:"Description,omitempty"` // 视频消息的描述
	} `xml:"Video" json:"Video"`
}

func NewVideo(to, from string, timestamp int64, mediaId, title, description string) (video *Video) {
	video = &Video{
		MessageHeader: MessageHeader{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   timestamp,
			MsgType:      MsgTypeVideo,
		},
	}
	video.Video.MediaId = mediaId
	video.Video.Title = title
	video.Video.Description = description
	return
}

type Article struct {
	Title       string `xml:"Title,omitempty"       json:"Title,omitempty"`       // 图文消息标题
	Description string `xml:"Description,omitempty" json:"Description,omitempty"` // 图文消息描述
	PicURL      string `xml:"PicUrl,omitempty"      json:"PicUrl,omitempty"`      // 图片链接, 支持JPG, PNG格式, 较好的效果为大图360*200, 小图200*200
	URL         string `xml:"Url,omitempty"         json:"Url,omitempty"`         // 点击图文消息跳转链接
}

const (
	NewsArticleCountLimit = 10 // 被动回复图文消息的文章数据最大数
)

type News struct {
	XMLName struct{} `xml:"xml" json:"-"`
	MessageHeader

	ArticleCount int       `xml:"ArticleCount"            json:"ArticleCount"` // 图文条数, 默认第一条为大图. 图文数不能超过10, 否则将会无响应
	Articles     []Article `xml:"Articles>item,omitempty" json:"Articles,omitempty"`
}

func NewNews(to, from string, timestamp int64, articles []Article) (news *News) {
	news = &News{
		MessageHeader: MessageHeader{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   timestamp,
			MsgType:      MsgTypeNews,
		},
	}
	news.Articles = articles
	news.ArticleCount = len(articles)
	return
}

// 检查 News 是否有效, 有效返回 nil, 否则返回错误信息
func (news *News) CheckValid() (err error) {
	n := len(news.Articles)
	if n != news.ArticleCount {
		err = fmt.Errorf("图文消息的 ArticleCount == %d, 实际文章个数为 %d", news.ArticleCount, n)
		return
	}
	if n <= 0 {
		err = errors.New("图文消息里没有文章")
		return
	}
	if n > NewsArticleCountLimit {
		err = fmt.Errorf("图文消息的文章个数不能超过 %d, 现在为 %d", NewsArticleCountLimit, n)
		return
	}
	return
}
