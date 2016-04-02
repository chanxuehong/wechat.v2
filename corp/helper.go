// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package corp

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"

	"github.com/chanxuehong/wechat/internal/util"
)

// 回复消息的 http body
type ResponseHttpBody struct {
	XMLName struct{} `xml:"xml" json:"-"`

	EncryptedMsg string `xml:"Encrypt"`
	MsgSignature string `xml:"MsgSignature"`
	Timestamp    int64  `xml:"TimeStamp"`
	Nonce        string `xml:"Nonce"`
}

// 回复消息给微信服务器.
//  要求 msg 是有效的消息数据结构(经过 encoding/xml marshal 后符合消息的格式);
//  如果有必要可以修改 Request 里面的某些值, 比如 Timestamp.
func WriteResponse(w http.ResponseWriter, r *Request, msg interface{}) (err error) {
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

	encryptedMsg := util.AESEncryptMsg(r.Random, rawMsgXML, r.CorpId, r.AESKey)
	base64EncryptedMsg := base64.StdEncoding.EncodeToString(encryptedMsg)

	responseHttpBody := ResponseHttpBody{
		EncryptedMsg: base64EncryptedMsg,
		Timestamp:    r.Timestamp,
		Nonce:        r.Nonce,
	}

	TimestampStr := strconv.FormatInt(responseHttpBody.Timestamp, 10)
	responseHttpBody.MsgSignature = util.MsgSign(r.AgentToken, TimestampStr,
		responseHttpBody.Nonce, responseHttpBody.EncryptedMsg)

	return xml.NewEncoder(w).Encode(&responseHttpBody)
}
