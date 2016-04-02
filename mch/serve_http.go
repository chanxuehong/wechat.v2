// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mch

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/chanxuehong/util/security"

	"github.com/chanxuehong/util"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request, queryValues url.Values, srv Server, errHandler ErrorHandler) {
	switch r.Method {
	case "POST":
		RawMsgXML, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errHandler.ServeError(w, r, err)
			return
		}

		msg, err := util.DecodeXMLToMap(bytes.NewReader(RawMsgXML))
		if err != nil {
			errHandler.ServeError(w, r, err)
			return
		}

		ReturnCode, ok := msg["return_code"]
		if ReturnCode == ReturnCodeSuccess || !ok {
			haveAppId := msg["appid"]
			wantAppId := srv.AppId()
			if wantAppId != "" && !security.SecureCompareString(haveAppId, wantAppId) {
				err = fmt.Errorf("the message's appid mismatch, have: %s, want: %s", haveAppId, wantAppId)
				errHandler.ServeError(w, r, err)
				return
			}

			haveMchId := msg["mch_id"]
			wantMchId := srv.MchId()
			if wantMchId != "" && !security.SecureCompareString(haveMchId, wantMchId) {
				err = fmt.Errorf("the message's mch_id mismatch, have: %s, want: %s", haveMchId, wantMchId)
				errHandler.ServeError(w, r, err)
				return
			}

			// 认证签名
			signature1, ok := msg["sign"]
			if !ok {
				err = errors.New("no sign parameter")
				errHandler.ServeError(w, r, err)
				return
			}
			signature2 := Sign(msg, srv.APIKey(), nil)
			if !security.SecureCompareString(signature1, signature2) {
				err = fmt.Errorf("check signature failed, \r\ninput: %q, \r\nlocal: %q", signature1, signature2)
				errHandler.ServeError(w, r, err)
				return
			}
		}

		req := &Request{
			HttpRequest: r,

			RawMsgXML: RawMsgXML,
			Msg:       msg,
		}
		srv.MessageHandler().ServeMessage(w, req)

	default:
		errHandler.ServeError(w, r, errors.New("Not expect Request.Method: "+r.Method))
	}
}
