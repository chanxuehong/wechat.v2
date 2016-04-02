// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package suite

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
	SuiteToken string // 请求消息所属套件的 Token

	HttpRequest *http.Request // 可以为 nil, 因为某些 http 框架没有提供此参数
	QueryValues url.Values    // 回调请求 URL 中的查询参数集合

	MsgSignature string // 回调请求 URL 中的消息体签名: msg_signature
	Timestamp    int64  // 回调请求 URL 中的时间戳: timestamp
	Nonce        string // 回调请求 URL 中的随机数: nonce

	RawMsgXML []byte        // 消息的"明文"XML 文本
	MixedMsg  *MixedMessage // RawMsgXML 解析后的消息

	AESKey  [32]byte // 当前消息 AES 加密的 key
	Random  []byte   // 当前消息加密时所用的 random, 16 bytes
	SuiteId string   // 当前消息的套件ID

}

// 微信服务器推送过来的消息(事件)的合集.
type MixedMessage struct {
	XMLName struct{} `xml:"xml" json:"-"`

	SuiteId   string `xml:"SuiteId"   json:"SuiteId"`
	InfoType  string `xml:"InfoType"  json:"InfoType"`
	Timestamp int64  `xml:"TimeStamp" json:"TimeStamp"`

	SuiteTicket string `xml:"SuiteTicket" json:"SuiteTicket"`
	AuthCorpId  string `xml:"AuthCorpId"  json:"AuthCorpId"`
}
