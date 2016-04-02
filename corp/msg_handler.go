// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package corp

import (
	"net/http"
	"net/url"
)

// 微信服务器推送过来的消息(事件)处理接口
type MessageHandler interface {
	ServeMessage(http.ResponseWriter, *Request)
}

type MessageHandlerFunc func(http.ResponseWriter, *Request)

func (fn MessageHandlerFunc) ServeMessage(w http.ResponseWriter, r *Request) {
	fn(w, r)
}

// 消息(事件)请求信息
type Request struct {
	AgentToken string // 请求消息所属企业号应用的 Token

	HttpRequest *http.Request // 可以为 nil, 因为某些 http 框架没有提供此参数
	QueryValues url.Values    // 回调请求 URL 中的查询参数集合

	MsgSignature string // 回调请求 URL 中的消息体签名: msg_signature
	Timestamp    int64  // 回调请求 URL 中的时间戳: timestamp
	Nonce        string // 回调请求 URL 中的随机数: nonce

	RawMsgXML []byte        // 消息的"明文"XML 文本
	MixedMsg  *MixedMessage // RawMsgXML 解析后的消息

	AESKey  [32]byte // 当前消息 AES 加密的 key
	Random  []byte   // 当前消息加密时所用的 random, 16 bytes
	CorpId  string   // 当前消息的企业号ID
	AgentId int64    // 当前消息的应用ID
}

// 微信服务器推送过来的消息(事件)通用的消息头
type MessageHeader struct {
	ToUserName   string `xml:"ToUserName"   json:"ToUserName"`
	FromUserName string `xml:"FromUserName" json:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"   json:"CreateTime"`
	MsgType      string `xml:"MsgType"      json:"MsgType"`
	AgentId      int64  `xml:"AgentID"      json:"AgentID"`
}

// 微信服务器推送过来的消息(事件)的合集.
type MixedMessage struct {
	XMLName struct{} `xml:"xml" json:"-"`
	MessageHeader

	MsgId int64 `xml:"MsgId" json:"MsgId"`

	Content      string  `xml:"Content"      json:"Content"`
	MediaId      string  `xml:"MediaId"      json:"MediaId"`
	PicURL       string  `xml:"PicUrl"       json:"PicUrl"`
	Format       string  `xml:"Format"       json:"Format"`
	ThumbMediaId string  `xml:"ThumbMediaId" json:"ThumbMediaId"`
	LocationX    float64 `xml:"Location_X"   json:"Location_X"`
	LocationY    float64 `xml:"Location_Y"   json:"Location_Y"`
	Scale        int     `xml:"Scale"        json:"Scale"`
	Label        string  `xml:"Label"        json:"Label"`

	Event    string `xml:"Event"    json:"Event"`
	EventKey string `xml:"EventKey" json:"EventKey"`

	ScanCodeInfo struct {
		ScanType   string `xml:"ScanType"   json:"ScanType"`
		ScanResult string `xml:"ScanResult" json:"ScanResult"`
	} `xml:"ScanCodeInfo" json:"ScanCodeInfo"`

	SendPicsInfo struct {
		Count   int `xml:"Count" json:"Count"`
		PicList []struct {
			PicMD5Sum string `xml:"PicMd5Sum" json:"PicMd5Sum"`
		} `xml:"PicList>item,omitempty" json:"PicList,omitempty"`
	} `xml:"SendPicsInfo" json:"SendPicsInfo"`

	SendLocationInfo struct {
		LocationX float64 `xml:"Location_X" json:"Location_X"`
		LocationY float64 `xml:"Location_Y" json:"Location_Y"`
		Scale     int     `xml:"Scale"      json:"Scale"`
		Label     string  `xml:"Label"      json:"Label"`
		PoiName   string  `xml:"Poiname"    json:"Poiname"`
	} `xml:"SendLocationInfo" json:"SendLocationInfo"`

	Latitude  float64 `xml:"Latitude"    json:"Latitude"`
	Longitude float64 `xml:"Longitude"   json:"Longitude"`
	Precision float64 `xml:"Precision"   json:"Precision"`
}
