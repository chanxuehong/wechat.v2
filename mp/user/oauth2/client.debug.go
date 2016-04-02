// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

// +build wechatdebug

package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chanxuehong/wechat/mp"
)

type TokenStorage interface {
	Get() (*Token, error)
	Put(*Token) error
}

type Client struct {
	Config Config

	// TokenStorage, Token 两个字段正常情况下只用指定一个, 如果两个同时被指定了, 优先使用 TokenStorage;
	// Client 会自动将最新的 Token 更新到 Client.Token 字段, 不管 Token 字段一开始是否被指定!!!
	TokenStorage TokenStorage
	Token        *Token

	HttpClient *http.Client // 如果 HttpClient == nil 则默认用 http.DefaultClient
}

func (clt *Client) getToken() (tk *Token, err error) {
	if clt.TokenStorage != nil {
		if tk, err = clt.TokenStorage.Get(); err != nil {
			return
		}
		if tk == nil {
			err = errors.New("Incorrect TokenStorage.Get()")
			return
		}
		clt.Token = tk // update local
		return
	}

	tk = clt.Token
	if tk == nil {
		err = errors.New("nil TokenStorage and nil Token")
		return
	}
	return
}

func (clt *Client) putToken(tk *Token) (err error) {
	if clt.TokenStorage != nil {
		if err = clt.TokenStorage.Put(tk); err != nil {
			return
		}
	}
	clt.Token = tk
	return
}

func (clt *Client) httpClient() *http.Client {
	if clt.HttpClient != nil {
		return clt.HttpClient
	}
	return http.DefaultClient
}

func (clt *Client) getJSON(url string, response interface{}) (err error) {
	mp.LogInfoln("[WECHAT_DEBUG] request url:", url)

	httpResp, err := clt.httpClient().Get(url)
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

	return json.Unmarshal(respBody, response)
}
