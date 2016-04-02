// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package send

import (
	"errors"
	"fmt"
)

const (
	MsgTypeText   = "text"
	MsgTypeImage  = "image"
	MsgTypeVoice  = "voice"
	MsgTypeVideo  = "video"
	MsgTypeFile   = "file"
	MsgTypeNews   = "news"
	MsgTypeMPNews = "mpnews"
)

type MessageHeader struct {
	ToUser  string `json:"touser,omitempty"`  // 非必须; 员工ID列表(消息接收者, 多个接收者用‘|’分隔, 最多支持1000个). 特殊情况: 指定为@all, 则向关注该企业应用的全部成员发送
	ToParty string `json:"toparty,omitempty"` // 非必须; 部门ID列表, 多个接收者用‘|’分隔, 最多支持100个. 当touser为@all时忽略本参数
	ToTag   string `json:"totag,omitempty"`   // 非必须; 标签ID列表, 多个接收者用‘|’分隔. 当touser为@all时忽略本参数

	MsgType string `json:"msgtype"`        // 必须; 消息类型
	AgentId int64  `json:"agentid"`        // 必须; 企业应用的id, 整型
	Safe    *int   `json:"safe,omitempty"` // 非必须; 表示是否是保密消息, 0表示否, 1表示是, 默认0
}

type Text struct {
	MessageHeader

	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

type Image struct {
	MessageHeader

	Image struct {
		MediaId string `json:"media_id"` // 图片媒体文件id, 可以调用上传媒体文件接口获取
	} `json:"image"`
}

type Voice struct {
	MessageHeader

	Voice struct {
		MediaId string `json:"media_id"` // 语音文件id, 可以调用上传媒体文件接口获取
	} `json:"voice"`
}

type Video struct {
	MessageHeader

	Video struct {
		MediaId     string `json:"media_id"`              // 视频媒体文件id, 可以调用上传媒体文件接口获取
		Title       string `json:"title,omitempty"`       // 视频消息的标题
		Description string `json:"description,omitempty"` // 视频消息的描述
	} `json:"video"`
}

type File struct {
	MessageHeader

	File struct {
		MediaId string `json:"media_id"` // 媒体文件id, 可以调用上传媒体文件接口获取
	} `json:"file"`
}

type NewsArticle struct {
	Title       string `json:"title,omitempty"`       // 图文消息标题
	Description string `json:"description,omitempty"` // 图文消息描述
	URL         string `json:"url,omitempty"`         // 点击后跳转的链接.
	PicURL      string `json:"picurl,omitempty"`      // 图文消息的图片链接, 支持JPG, PNG格式, 较好的效果为大图640*320, 小图80*80. 如不填, 在客户端不显示图片
}

const NewsArticleCountLimit = 10

// News 消息, 注意沒有 Safe 字段.
type News struct {
	MessageHeader

	News struct {
		Articles []NewsArticle `json:"articles,omitempty"` // 图文消息, 一个图文消息支持1到10条图文
	} `json:"news"`
}

// 检查 News 是否有效, 有效返回 nil, 否则返回错误信息
func (this *News) CheckValid() (err error) {
	n := len(this.News.Articles)
	if n <= 0 {
		err = errors.New("没有有效的图文消息")
		return
	}
	if n > NewsArticleCountLimit {
		err = fmt.Errorf("图文消息的文章个数不能超过 %d, 现在为 %d", NewsArticleCountLimit, n)
		return
	}
	return
}

type MPNewsArticle struct {
	ThumbMediaId     string `json:"thumb_media_id"`               // 图文消息缩略图的media_id, 可以在上传多媒体文件接口中获得. 此处thumb_media_id即上传接口返回的media_id
	Title            string `json:"title"`                        // 图文消息的标题
	Author           string `json:"author,omitempty"`             // 图文消息的作者
	ContentSourceURL string `json:"content_source_url,omitempty"` // 图文消息点击"阅读原文"之后的页面链接
	Content          string `json:"content"`                      // 图文消息的内容, 支持html标签
	Digest           string `json:"digest,omitempty"`             // 图文消息的描述
	ShowCoverPic     int    `json:"show_cover_pic"`               // 是否显示封面, 1为显示, 0为不显示
}

func (article *MPNewsArticle) SetShowCoverPic(b bool) {
	if b {
		article.ShowCoverPic = 1
	} else {
		article.ShowCoverPic = 0
	}
}

// MPNews 消息与 News 消息类似, 不同的是图文消息内容存储在微信后台, 并且支持保密选项.
type MPNews struct {
	MessageHeader

	MPNews struct {
		Articles []MPNewsArticle `json:"articles,omitempty"` // 多条图文消息信息, 默认第一个item为大图, 注意, 如果图文数超过10, 则将会无响应
	} `json:"mpnews"`
}

// 检查 MPNews 是否有效, 有效返回 nil, 否则返回错误信息
func (this *MPNews) CheckValid() (err error) {
	n := len(this.MPNews.Articles)
	if n <= 0 {
		err = errors.New("没有有效的图文消息")
		return
	}
	if n > NewsArticleCountLimit {
		err = fmt.Errorf("图文消息的文章个数不能超过 %d, 现在为 %d", NewsArticleCountLimit, n)
		return
	}
	return
}
