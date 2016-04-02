// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

// +build !wechatdebug

package mp

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/chanxuehong/util/security"

	"github.com/chanxuehong/wechat/internal/util"
)

// 安全模式, 微信服务器推送过来的 http body
type RequestHttpBody struct {
	XMLName struct{} `xml:"xml" json:"-"`

	ToUserName   string `xml:"ToUserName" json:"ToUserName"`
	EncryptedMsg string `xml:"Encrypt"    json:"Encrypt"`
}

// ServeHTTP 处理 http 消息请求
//  NOTE: 调用者保证所有参数有效
func ServeHTTP(w http.ResponseWriter, r *http.Request, queryValues url.Values, srv Server, errHandler ErrorHandler) {
	switch r.Method {
	case "POST": // 消息处理
		switch encryptType := queryValues.Get("encrypt_type"); encryptType {
		case "aes": // 安全模式, 兼容模式
			signature := queryValues.Get("signature") // 只读取, 不做校验

			msgSignature1 := queryValues.Get("msg_signature")
			if msgSignature1 == "" {
				errHandler.ServeError(w, r, errors.New("msg_signature is empty"))
				return
			}

			timestampStr := queryValues.Get("timestamp")
			if timestampStr == "" {
				errHandler.ServeError(w, r, errors.New("timestamp is empty"))
				return
			}

			timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
			if err != nil {
				err = errors.New("can not parse timestamp to int64: " + timestampStr)
				errHandler.ServeError(w, r, err)
				return
			}

			nonce := queryValues.Get("nonce")
			if nonce == "" {
				errHandler.ServeError(w, r, errors.New("nonce is empty"))
				return
			}

			var requestHttpBody RequestHttpBody
			if err := xml.NewDecoder(r.Body).Decode(&requestHttpBody); err != nil {
				errHandler.ServeError(w, r, err)
				return
			}

			// 安全考虑验证下 ToUserName
			haveToUserName := requestHttpBody.ToUserName
			wantToUserName := srv.OriId()
			if wantToUserName != "" && !security.SecureCompareString(haveToUserName, wantToUserName) {
				err := fmt.Errorf("the RequestHttpBody's ToUserName mismatch, have: %s, want: %s", haveToUserName, wantToUserName)
				errHandler.ServeError(w, r, err)
				return
			}

			token := srv.Token()

			// 验证签名
			msgSignature2 := util.MsgSign(token, timestampStr, nonce, requestHttpBody.EncryptedMsg)
			if !security.SecureCompareString(msgSignature1, msgSignature2) {
				err := fmt.Errorf("check msg_signature failed, input: %s, local: %s", msgSignature1, msgSignature2)
				errHandler.ServeError(w, r, err)
				return
			}

			// 解密
			encryptedMsgBytes, err := base64.StdEncoding.DecodeString(requestHttpBody.EncryptedMsg)
			if err != nil {
				errHandler.ServeError(w, r, err)
				return
			}

			aesKey := srv.CurrentAESKey()
			random, rawMsgXML, haveAppIdBytes, err := util.AESDecryptMsg(encryptedMsgBytes, aesKey)
			if err != nil {
				// 尝试用上一次的 AESKey 来解密
				lastAESKey, isLastAESKeyValid := srv.LastAESKey()
				if !isLastAESKeyValid {
					errHandler.ServeError(w, r, err)
					return
				}

				aesKey = lastAESKey // NOTE

				random, rawMsgXML, haveAppIdBytes, err = util.AESDecryptMsg(encryptedMsgBytes, aesKey)
				if err != nil {
					errHandler.ServeError(w, r, err)
					return
				}
			}
			haveAppId := string(haveAppIdBytes)
			wantAppId := srv.AppId()
			if wantAppId != "" && !security.SecureCompareString(haveAppId, wantAppId) {
				err := fmt.Errorf("the message's appid mismatch, have: %s, want: %s", haveAppId, wantAppId)
				errHandler.ServeError(w, r, err)
				return
			}

			// 解密成功, 解析 MixedMessage
			var mixedMsg MixedMessage
			if err := xml.Unmarshal(rawMsgXML, &mixedMsg); err != nil {
				errHandler.ServeError(w, r, err)
				return
			}

			// 安全考虑再次验证 ToUserName
			if haveToUserName != mixedMsg.ToUserName {
				err := fmt.Errorf("the RequestHttpBody's ToUserName(==%s) mismatch the MixedMessage's ToUserName(==%s)", haveToUserName, mixedMsg.ToUserName)
				errHandler.ServeError(w, r, err)
				return
			}

			// 成功, 交给 MessageHandler
			req := &Request{
				Token: token,

				HttpRequest: r,
				QueryValues: queryValues,

				Signature: signature,
				Timestamp: timestamp,
				Nonce:     nonce,

				EncryptType: encryptType,
				RawMsgXML:   rawMsgXML,
				MixedMsg:    &mixedMsg,

				MsgSignature: msgSignature1,
				AESKey:       aesKey,
				Random:       random,
				AppId:        haveAppId,
			}
			srv.MessageHandler().ServeMessage(w, req)

		case "", "raw": // 明文模式
			signature1 := queryValues.Get("signature")
			if signature1 == "" {
				errHandler.ServeError(w, r, errors.New("signature is empty"))
				return
			}

			timestampStr := queryValues.Get("timestamp")
			if timestampStr == "" {
				errHandler.ServeError(w, r, errors.New("timestamp is empty"))
				return
			}

			timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
			if err != nil {
				err = errors.New("can not parse timestamp to int64: " + timestampStr)
				errHandler.ServeError(w, r, err)
				return
			}

			nonce := queryValues.Get("nonce")
			if nonce == "" {
				errHandler.ServeError(w, r, errors.New("nonce is empty"))
				return
			}

			token := srv.Token()

			signature2 := util.Sign(token, timestampStr, nonce)
			if !security.SecureCompareString(signature1, signature2) {
				err := fmt.Errorf("check signature failed, input: %s, local: %s", signature1, signature2)
				errHandler.ServeError(w, r, err)
				return
			}

			// 验证签名成功, 解析 MixedMessage
			rawMsgXML, err := ioutil.ReadAll(r.Body)
			if err != nil {
				errHandler.ServeError(w, r, err)
				return
			}

			var mixedMsg MixedMessage
			if err := xml.Unmarshal(rawMsgXML, &mixedMsg); err != nil {
				errHandler.ServeError(w, r, err)
				return
			}

			// 安全考虑验证 ToUserName
			haveToUserName := mixedMsg.ToUserName
			wantToUserName := srv.OriId()
			if wantToUserName != "" && !security.SecureCompareString(haveToUserName, wantToUserName) {
				err := fmt.Errorf("the Message's ToUserName mismatch, have: %s, want: %s", haveToUserName, wantToUserName)
				errHandler.ServeError(w, r, err)
				return
			}

			// 成功, 交给 MessageHandler
			req := &Request{
				Token: token,

				HttpRequest: r,
				QueryValues: queryValues,

				Signature: signature1,
				Timestamp: timestamp,
				Nonce:     nonce,

				EncryptType: encryptType,
				RawMsgXML:   rawMsgXML,
				MixedMsg:    &mixedMsg,
			}
			srv.MessageHandler().ServeMessage(w, req)

		default: // 未知的加密类型
			err := errors.New("unknown encrypt_type: " + encryptType)
			errHandler.ServeError(w, r, err)
			return
		}

	case "GET": // 验证回调 url 有效性
		signature1 := queryValues.Get("signature")
		if signature1 == "" {
			errHandler.ServeError(w, r, errors.New("signature is empty"))
			return
		}

		timestamp := queryValues.Get("timestamp")
		if timestamp == "" {
			errHandler.ServeError(w, r, errors.New("timestamp is empty"))
			return
		}

		nonce := queryValues.Get("nonce")
		if nonce == "" {
			errHandler.ServeError(w, r, errors.New("nonce is empty"))
			return
		}

		echostr := queryValues.Get("echostr")
		if echostr == "" {
			errHandler.ServeError(w, r, errors.New("echostr is empty"))
			return
		}

		signature2 := util.Sign(srv.Token(), timestamp, nonce)
		if !security.SecureCompareString(signature1, signature2) {
			err := fmt.Errorf("check signature failed, input: %s, local: %s", signature1, signature2)
			errHandler.ServeError(w, r, err)
			return
		}

		io.WriteString(w, echostr)
	}
}
