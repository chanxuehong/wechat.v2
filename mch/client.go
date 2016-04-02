// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

// +build !wechatdebug

package mch

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/chanxuehong/util"
)

type Proxy struct {
	appId      string
	mchId      string
	apiKey     string
	httpClient *http.Client
}

func (pxy *Proxy) AppId() string {
	return pxy.appId
}
func (pxy *Proxy) MchId() string {
	return pxy.mchId
}

// 创建一个新的 Proxy.
//  如果 httpClient == nil 则默认用 http.DefaultClient.
func NewProxy(appId, mchId, apiKey string, httpClient *http.Client) *Proxy {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Proxy{
		appId:      appId,
		mchId:      mchId,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// 微信支付通用请求方法.
//  注意: err == nil 表示协议状态都为 SUCCESS(return_code == SUCCESS).
func (pxy *Proxy) PostXML(url string, req map[string]string) (resp map[string]string, err error) {
	bodyBuf := textBufferPool.Get().(*bytes.Buffer)
	bodyBuf.Reset()
	defer textBufferPool.Put(bodyBuf)

	if err = util.EncodeXMLFromMap(bodyBuf, req, "xml"); err != nil {
		return
	}

	httpResp, err := pxy.httpClient.Post(url, "text/xml; charset=utf-8", bodyBuf)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	if resp, err = util.DecodeXMLToMap(httpResp.Body); err != nil {
		return
	}

	// 判断协议状态
	ReturnCode, ok := resp["return_code"]
	if !ok {
		err = errors.New("no return_code parameter")
		return
	}
	if ReturnCode != ReturnCodeSuccess {
		err = &Error{
			ReturnCode: ReturnCode,
			ReturnMsg:  resp["return_msg"],
		}
		return
	}

	// 安全考虑, 做下验证
	appId, ok := resp["appid"]
	if ok && appId != pxy.appId {
		err = fmt.Errorf("appid mismatch, have: %q, want: %q", appId, pxy.appId)
		return
	}
	mchId, ok := resp["mch_id"]
	if ok && mchId != pxy.mchId {
		err = fmt.Errorf("mch_id mismatch, have: %q, want: %q", mchId, pxy.mchId)
		return
	}

	// 认证签名
	signature1, ok := resp["sign"]
	if !ok {
		err = errors.New("no sign parameter")
		return
	}
	signature2 := Sign(resp, pxy.apiKey, nil)
	if signature1 != signature2 {
		err = fmt.Errorf("check signature failed, \r\ninput: %q, \r\nlocal: %q", signature1, signature2)
		return
	}
	return
}
