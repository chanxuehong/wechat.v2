// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

// 群发给所有用户的消息数据结构.
package mass2all

const (
	MsgTypeText   = "text"
	MsgTypeImage  = "image"
	MsgTypeVoice  = "voice"
	MsgTypeVideo  = "mpvideo"
	MsgTypeNews   = "mpnews"
	MsgTypeWxCard = "wxcard"
)

type MessageHeader struct {
	Filter struct {
		IsToAll bool `json:"is_to_all"`
	} `json:"filter"`
	MsgType string `json:"msgtype"`
}

type Text struct {
	MessageHeader
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

func NewText(content string) *Text {
	var msg Text
	msg.MsgType = MsgTypeText
	msg.Filter.IsToAll = true
	msg.Text.Content = content
	return &msg
}

type Image struct {
	MessageHeader
	Image struct {
		MediaId string `json:"media_id"`
	} `json:"image"`
}

func NewImage(mediaId string) *Image {
	var msg Image
	msg.MsgType = MsgTypeImage
	msg.Filter.IsToAll = true
	msg.Image.MediaId = mediaId
	return &msg
}

type Voice struct {
	MessageHeader
	Voice struct {
		MediaId string `json:"media_id"`
	} `json:"voice"`
}

func NewVoice(mediaId string) *Voice {
	var msg Voice
	msg.MsgType = MsgTypeVoice
	msg.Filter.IsToAll = true
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
func NewVideo(mediaId string) *Video {
	var msg Video
	msg.MsgType = MsgTypeVideo
	msg.Filter.IsToAll = true
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
func NewNews(mediaId string) *News {
	var msg News
	msg.MsgType = MsgTypeNews
	msg.Filter.IsToAll = true
	msg.News.MediaId = mediaId
	return &msg
}

// 卡券消息
type WxCard struct {
	MessageHeader
	WxCard struct {
		CardId string `json:"card_id"`
	} `json:"wxcard"`
}

// 新建卡券, 特别注意: 目前该接口仅支持填入非自定义code的卡券和预存模式的自定义code卡券.
func NewWxCard(cardId string) *WxCard {
	var msg WxCard
	msg.MsgType = MsgTypeWxCard
	msg.Filter.IsToAll = true
	msg.WxCard.CardId = cardId
	return &msg
}
