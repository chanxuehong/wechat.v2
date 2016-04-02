// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mp

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
	Token string // 请求消息所属公众号的 Token

	HttpRequest *http.Request // 可以为 nil, 因为某些 http 框架没有提供此参数
	QueryValues url.Values    // 回调请求 URL 中的查询参数集合

	Signature string // 回调请求 URL 中的签名: signature
	Timestamp int64  // 回调请求 URL 中的时间戳: timestamp
	Nonce     string // 回调请求 URL 中的随机数: nonce

	EncryptType string        // 请求 URL 中的加密方式: encrypt_type
	RawMsgXML   []byte        // 消息的 XML 文本, 对于加密模式是解密后的消息
	MixedMsg    *MixedMessage // RawMsgXML 解析后的消息

	// 下面的字段是 AES 模式才有的
	MsgSignature string   // 请求 URL 中的消息体签名: msg_signature
	AESKey       [32]byte // 当前消息 AES 加密的 key
	Random       []byte   // 当前消息加密时所用的 random, 16 bytes
	AppId        string   // 当前消息加密时所用的 AppId
}

// 微信服务器推送过来的消息(事件)通用的消息头
type MessageHeader struct {
	ToUserName   string `xml:"ToUserName"   json:"ToUserName"`
	FromUserName string `xml:"FromUserName" json:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"   json:"CreateTime"`
	MsgType      string `xml:"MsgType"      json:"MsgType"`
}

// 微信服务器推送过来的消息(事件)的合集.
type MixedMessage struct {
	XMLName struct{} `xml:"xml" json:"-"`
	MessageHeader

	MsgId int64 `xml:"MsgId" json:"MsgId"`
	MsgID int64 `xml:"MsgID" json:"MsgID"` // 群发消息和模板消息的消息ID, 而不是此事件消息的ID

	Content      string  `xml:"Content"      json:"Content"`
	MediaId      string  `xml:"MediaId"      json:"MediaId"`
	PicURL       string  `xml:"PicUrl"       json:"PicUrl"`
	Format       string  `xml:"Format"       json:"Format"`
	Recognition  string  `xml:"Recognition"  json:"Recognition"`
	ThumbMediaId string  `xml:"ThumbMediaId" json:"ThumbMediaId"`
	LocationX    float64 `xml:"Location_X"   json:"Location_X"`
	LocationY    float64 `xml:"Location_Y"   json:"Location_Y"`
	Scale        int     `xml:"Scale"        json:"Scale"`
	Label        string  `xml:"Label"        json:"Label"`
	Title        string  `xml:"Title"        json:"Title"`
	Description  string  `xml:"Description"  json:"Description"`
	URL          string  `xml:"Url"          json:"Url"`

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

	Ticket      string  `xml:"Ticket"      json:"Ticket"`
	Latitude    float64 `xml:"Latitude"    json:"Latitude"`
	Longitude   float64 `xml:"Longitude"   json:"Longitude"`
	Precision   float64 `xml:"Precision"   json:"Precision"`
	Status      string  `xml:"Status"      json:"Status"`
	TotalCount  int     `xml:"TotalCount"  json:"TotalCount"`
	FilterCount int     `xml:"FilterCount" json:"FilterCount"`
	SentCount   int     `xml:"SentCount"   json:"SentCount"`
	ErrorCount  int     `xml:"ErrorCount"  json:"ErrorCount"`

	// merchant
	OrderId     string `xml:"OrderId"     json:"OrderId"`
	OrderStatus int    `xml:"OrderStatus" json:"OrderStatus"`
	ProductId   string `xml:"ProductId"   json:"ProductId"`
	SKUInfo     string `xml:"SkuInfo"     json:"SkuInfo"`

	// card
	CardId          string `xml:"CardId"          json:"CardId"`
	IsGiveByFriend  int    `xml:"IsGiveByFriend"  json:"IsGiveByFriend"`
	FriendUserName  string `xml:"FriendUserName"  json:"FriendUserName"`
	UserCardCode    string `xml:"UserCardCode"    json:"UserCardCode"`
	OldUserCardCode string `xml:"OldUserCardCode" json:"OldUserCardCode"`
	ConsumeSource   string `xml:"ConsumeSource"   json:"ConsumeSource"`
	OuterId         int64  `xml:"OuterId"         json:"OuterId"`

	// poi
	UniqId string `xml:"UniqId" json:"UniqId"`
	PoiId  int64  `xml:"PoiId"  json:"PoiId"`
	Result string `xml:"Result" json:"Result"`
	Msg    string `xml:"Msg"    json:"Msg"`

	// dkf
	KfAccount     string `xml:"KfAccount"     json:"KfAccount"`
	FromKfAccount string `xml:"FromKfAccount" json:"FromKfAccount"`
	ToKfAccount   string `xml:"ToKfAccount"   json:"ToKfAccount"`

	// shakearound
	ChosenBeacon  ChosenBeacon   `xml:"ChosenBeacon"                         json:"ChosenBeacon"`
	AroundBeacons []AroundBeacon `xml:"AroundBeacons>AroundBeacon,omitempty" json:"AroundBeacons,omitempty"`

	// bizwifi
	ConnectTime int64  `xml:"ConnectTime" json:"ConnectTime"`
	ExpireTime  int64  `xml:"ExpireTime"  json:"ExpireTime"`
	VendorId    string `xml:"VendorId"    json:"VendorId"`
	PlaceId     int64  `xml:"PlaceId"     json:"PlaceId"`
	DeviceNo    string `xml:"DeviceNo"    json:"DeviceNo"`
}

// 和 github.com/chanxuehong/wechat/mp/shakearound.ChosenBeacon 一样, 同步修改
type ChosenBeacon struct {
	UUID     string  `xml:"Uuid"     json:"Uuid"`
	Major    int     `xml:"Major"    json:"Major"`
	Minor    int     `xml:"Minor"    json:"Minor"`
	Distance float64 `xml:"Distance" json:"Distance"`
}

// 和 github.com/chanxuehong/wechat/mp/shakearound.AroundBeacon 一样, 同步修改
type AroundBeacon struct {
	UUID     string  `xml:"Uuid"     json:"Uuid"`
	Major    int     `xml:"Major"    json:"Major"`
	Minor    int     `xml:"Minor"    json:"Minor"`
	Distance float64 `xml:"Distance" json:"Distance"`
}
