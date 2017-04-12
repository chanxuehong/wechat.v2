// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mp

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"

	"github.com/chanxuehong/wechat/internal/util"
)

// 转义微信中的错误编码，防止json decode的时候出错
func EscapeCtrl(ctrl []byte) (esc []byte) {
	u := []byte(`\u0000`)
	for i, ch := range ctrl {
		if ch <= 31 {
			if esc == nil {
				esc = append(make([]byte, 0, len(ctrl)+len(u)), ctrl[:i]...)
			}
			esc = append(esc, u...)
			hex.Encode(esc[len(esc)-2:], ctrl[i:i+1])
			continue
		}
		if esc != nil {
			esc = append(esc, ch)
		}
	}
	if esc == nil {
		return ctrl
	}
	return esc
}

// 明文模式下回复消息给微信服务器.
//  要求 msg 是有效的消息数据结构(经过 encoding/xml marshal 后符合微信消息格式);
//  如果有必要可以修改 Request 里面的某些值, 比如 Timestamp, Nonce, Random.
func WriteRawResponse(w http.ResponseWriter, r *Request, msg interface{}) (err error) {
	if w == nil {
		return errors.New("nil http.ResponseWriter")
	}
	if msg == nil {
		return errors.New("nil message")
	}
	return xml.NewEncoder(w).Encode(msg)
}

// 安全模式下回复消息的 http body
type ResponseHttpBody struct {
	XMLName struct{} `xml:"xml" json:"-"`

	EncryptedMsg string `xml:"Encrypt"      json:"Encrypt"`
	MsgSignature string `xml:"MsgSignature" json:"MsgSignature"`
	Timestamp    int64  `xml:"TimeStamp"    json:"TimeStamp"`
	Nonce        string `xml:"Nonce"        json:"Nonce"`
}

// 安全模式下回复消息给微信服务器.
//  要求 msg 是有效的消息数据结构(经过 encoding/xml marshal 后符合微信消息格式);
//  如果有必要可以修改 Request 里面的某些值, 比如 Timestamp, Nonce, Random.
func WriteAESResponse(w http.ResponseWriter, r *Request, msg interface{}) (err error) {
	if w == nil {
		return errors.New("nil http.ResponseWriter")
	}
	if r == nil {
		return errors.New("nil Request")
	}
	if msg == nil {
		return errors.New("nil message")
	}

	rawMsgXML, err := xml.Marshal(msg)
	if err != nil {
		return
	}

	encryptedMsg := util.AESEncryptMsg(r.Random, rawMsgXML, r.AppId, r.AESKey)
	base64EncryptedMsg := base64.StdEncoding.EncodeToString(encryptedMsg)

	responseHttpBody := ResponseHttpBody{
		EncryptedMsg: base64EncryptedMsg,
		Timestamp:    r.Timestamp,
		Nonce:        r.Nonce,
	}

	TimestampStr := strconv.FormatInt(responseHttpBody.Timestamp, 10)
	responseHttpBody.MsgSignature = util.MsgSign(r.Token, TimestampStr, responseHttpBody.Nonce, responseHttpBody.EncryptedMsg)

	return xml.NewEncoder(w).Encode(&responseHttpBody)
}
