// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

// +build wechatdebug

package component

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/chanxuehong/wechat/mp"
)

type Client struct {
	AccessTokenServer
	AppId      string
	HttpClient *http.Client
}

// 创建一个新的 Client.
//  如果 clt == nil 则默认用 http.DefaultClient
func NewClient(appId string, srv AccessTokenServer, clt *http.Client) *Client {
	if srv == nil {
		panic("nil AccessTokenServer")
	}
	if clt == nil {
		clt = http.DefaultClient
	}

	return &Client{
		AccessTokenServer: srv,
		AppId:             appId,
		HttpClient:        clt,
	}
}

// 用 encoding/json 把 request marshal 为 JSON, 放入 http 请求的 body 中,
// POST 到微信服务器, 然后将微信服务器返回的 JSON 用 encoding/json 解析到 response.
//
//  NOTE:
//  1. 一般不用调用这个方法, 请直接调用高层次的封装方法;
//  2. 最终的 URL == incompleteURL + component_access_token;
//  3. response 格式有要求, 要么是 *mp.Error, 要么是下面结构体的指针(注意 Error 必须是第一个 Field):
//      struct {
//          mp.Error
//          ...
//      }
func (clt *Client) PostJSON(incompleteURL string, request interface{}, response interface{}) (err error) {
	buf := textBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer textBufferPool.Put(buf)

	if err = json.NewEncoder(buf).Encode(request); err != nil {
		return
	}
	requestBytes := buf.Bytes()

	token, err := clt.Token()
	if err != nil {
		return
	}

	hasRetried := false
RETRY:
	finalURL := incompleteURL + url.QueryEscape(token)

	mp.LogInfoln("[WECHAT_DEBUG] request url:", finalURL)
	mp.LogInfoln("[WECHAT_DEBUG] request json:", string(requestBytes))

	httpResp, err := clt.HttpClient.Post(finalURL, "application/json; charset=utf-8", bytes.NewReader(requestBytes))
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return
	}
	mp.LogInfoln("[WECHAT_DEBUG] response json:", string(respBody))

	if err = json.Unmarshal(respBody, response); err != nil {
		return
	}

	var ErrorStructValue reflect.Value // mp.Error

	// 下面的代码对 response 有特定要求, 见此函数 NOTE
	responseStructValue := reflect.ValueOf(response).Elem()
	if v := responseStructValue.Field(0); v.Kind() == reflect.Struct {
		ErrorStructValue = v
	} else {
		ErrorStructValue = responseStructValue
	}

	switch ErrCode := ErrorStructValue.Field(0).Int(); ErrCode {
	case mp.ErrCodeOK:
		return
	case mp.ErrCodeInvalidCredential, mp.ErrCodeAccessTokenExpired:
		ErrMsg := ErrorStructValue.Field(1).String()
		mp.LogInfoln("[WECHAT_RETRY] err_code:", ErrCode, ", err_msg:", ErrMsg)
		mp.LogInfoln("[WECHAT_RETRY] current token:", token)

		if !hasRetried {
			hasRetried = true

			if token, err = clt.TokenRefresh(); err != nil {
				return
			}
			mp.LogInfoln("[WECHAT_RETRY] new token:", token)

			responseStructValue.Set(reflect.New(responseStructValue.Type()).Elem())
			goto RETRY
		}
		mp.LogInfoln("[WECHAT_RETRY] fallthrough, current token:", token)
		fallthrough
	default:
		return
	}
}

// GET 微信资源, 然后将微信服务器返回的 JSON 用 encoding/json 解析到 response.
//
//  NOTE:
//  1. 一般不用调用这个方法, 请直接调用高层次的封装方法;
//  2. 最终的 URL == incompleteURL + component_access_token;
//  3. response 格式有要求, 要么是 *mp.Error, 要么是下面结构体的指针(注意 Error 必须是第一个 Field):
//      struct {
//          mp.Error
//          ...
//      }
func (clt *Client) GetJSON(incompleteURL string, response interface{}) (err error) {
	token, err := clt.Token()
	if err != nil {
		return
	}

	hasRetried := false
RETRY:
	finalURL := incompleteURL + url.QueryEscape(token)

	httpResp, err := clt.HttpClient.Get(finalURL)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return
	}
	mp.LogInfoln("[WECHAT_DEBUG] request url:", finalURL)
	mp.LogInfoln("[WECHAT_DEBUG] response json:", string(respBody))

	if err = json.Unmarshal(respBody, response); err != nil {
		return
	}

	var ErrorStructValue reflect.Value // mp.Error

	// 下面的代码对 response 有特定要求, 见此函数 NOTE
	responseStructValue := reflect.ValueOf(response).Elem()
	if v := responseStructValue.Field(0); v.Kind() == reflect.Struct {
		ErrorStructValue = v
	} else {
		ErrorStructValue = responseStructValue
	}

	switch ErrCode := ErrorStructValue.Field(0).Int(); ErrCode {
	case mp.ErrCodeOK:
		return
	case mp.ErrCodeInvalidCredential, mp.ErrCodeAccessTokenExpired:
		ErrMsg := ErrorStructValue.Field(1).String()
		mp.LogInfoln("[WECHAT_RETRY] err_code:", ErrCode, ", err_msg:", ErrMsg)
		mp.LogInfoln("[WECHAT_RETRY] current token:", token)

		if !hasRetried {
			hasRetried = true

			if token, err = clt.TokenRefresh(); err != nil {
				return
			}
			mp.LogInfoln("[WECHAT_RETRY] new token:", token)

			responseStructValue.Set(reflect.New(responseStructValue.Type()).Elem())
			goto RETRY
		}
		mp.LogInfoln("[WECHAT_RETRY] fallthrough, current token:", token)
		fallthrough
	default:
		return
	}
}
