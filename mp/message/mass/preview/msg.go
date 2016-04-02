// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

// 预览接口的消息数据结构.
package preview

const (
	MsgTypeText   = "text"
	MsgTypeImage  = "image"
	MsgTypeVoice  = "voice"
	MsgTypeVideo  = "mpvideo"
	MsgTypeNews   = "mpnews"
	MsgTypeWxCard = "wxcard"
)

type MessageHeader struct {
	ToWxName string `json:"towxname,omitempty"`
	ToUser   string `json:"touser,omitempty"`
	MsgType  string `json:"msgtype"`
}

type Text struct {
	MessageHeader
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

func NewText(touser, content string) *Text {
	var msg Text
	msg.MsgType = MsgTypeText
	msg.ToUser = touser
	msg.Text.Content = content
	return &msg
}

func NewText2(towxname, content string) *Text {
	var msg Text
	msg.MsgType = MsgTypeText
	msg.ToWxName = towxname
	msg.Text.Content = content
	return &msg
}

type Image struct {
	MessageHeader
	Image struct {
		MediaId string `json:"media_id"`
	} `json:"image"`
}

func NewImage(touser, mediaId string) *Image {
	var msg Image
	msg.MsgType = MsgTypeImage
	msg.ToUser = touser
	msg.Image.MediaId = mediaId
	return &msg
}

func NewImage2(towxname, mediaId string) *Image {
	var msg Image
	msg.MsgType = MsgTypeImage
	msg.ToWxName = towxname
	msg.Image.MediaId = mediaId
	return &msg
}

type Voice struct {
	MessageHeader
	Voice struct {
		MediaId string `json:"media_id"`
	} `json:"voice"`
}

func NewVoice(touser, mediaId string) *Voice {
	var msg Voice
	msg.MsgType = MsgTypeVoice
	msg.ToUser = touser
	msg.Voice.MediaId = mediaId
	return &msg
}

func NewVoice2(towxname, mediaId string) *Voice {
	var msg Voice
	msg.MsgType = MsgTypeVoice
	msg.ToWxName = towxname
	msg.Voice.MediaId = mediaId
	return &msg
}

type Video struct {
	MessageHeader
	Video struct {
		MediaId string `json:"media_id"`
	} `json:"mpvideo"`
}

// 新建视频消息
//  NOTE: 对于临时素材, mediaId 应该通过 media.Client.CreateVideo 得到
func NewVideo(touser, mediaId string) *Video {
	var msg Video
	msg.MsgType = MsgTypeVideo
	msg.ToUser = touser
	msg.Video.MediaId = mediaId
	return &msg
}

// 新建视频消息
//  NOTE: 对于临时素材, mediaId 应该通过 media.Client.CreateVideo 得到
func NewVideo2(towxname, mediaId string) *Video {
	var msg Video
	msg.MsgType = MsgTypeVideo
	msg.ToWxName = towxname
	msg.Video.MediaId = mediaId
	return &msg
}

// 图文消息
type News struct {
	MessageHeader
	News struct {
		MediaId string `json:"media_id"`
	} `json:"mpnews"`
}

// 新建图文消息
//  NOTE: 对于临时素材, mediaId 应该通过 media.Client.CreateNews 得到
func NewNews(touser, mediaId string) *News {
	var msg News
	msg.MsgType = MsgTypeNews
	msg.ToUser = touser
	msg.News.MediaId = mediaId
	return &msg
}

// 新建图文消息
//  NOTE: 对于临时素材, mediaId 应该通过 media.Client.CreateNews 得到
func NewNews2(towxname, mediaId string) *News {
	var msg News
	msg.MsgType = MsgTypeNews
	msg.ToWxName = towxname
	msg.News.MediaId = mediaId
	return &msg
}

// 卡券消息
type WxCard struct {
	MessageHeader
	WxCard struct {
		CardId  string `json:"card_id"`
		CardExt string `json:"card_ext,omitempty"`
	} `json:"wxcard"`
}

// 新建卡券, 特别注意: 目前该接口仅支持填入非自定义code的卡券和预存模式的自定义code卡券.
func NewWxCard(toUser, cardId, cardExt string) *WxCard {
	var msg WxCard
	msg.MsgType = MsgTypeWxCard
	msg.ToUser = toUser
	msg.WxCard.CardId = cardId
	msg.WxCard.CardExt = cardExt
	return &msg
}

// 新建卡券, 特别注意: 目前该接口仅支持填入非自定义code的卡券和预存模式的自定义code卡券.
func NewWxCard2(towxname, cardId, cardExt string) *WxCard {
	var msg WxCard
	msg.MsgType = MsgTypeWxCard
	msg.ToWxName = towxname
	msg.WxCard.CardId = cardId
	msg.WxCard.CardExt = cardExt
	return &msg
}
