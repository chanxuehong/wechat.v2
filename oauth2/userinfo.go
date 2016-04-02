// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 获取用户信息.
//  lang 可以为空值.
func (clt *Client) GetUserInfo(userinfo interface{}, lang string) (err error) {
	if clt.Config == nil {
		err = errors.New("nil Config")
		return
	}

	tk, err := clt.getToken()
	if err != nil {
		return
	}

	// 过期自动刷新 Token
	if tk.AccessTokenExpired() {
		if tk, err = clt.tokenRefresh(tk); err != nil {
			return
		}
	}

	httpResp, err := clt.httpClient().Get(clt.Config.UserInfoURL(tk.AccessToken, tk.OpenId, lang))
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	httpRespBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return
	}

	var errResult Error
	if err = json.Unmarshal(httpRespBytes, &errResult); err != nil {
		return
	}
	if errResult.ErrCode != ErrCodeOK {
		return &errResult
	}
	return json.Unmarshal(httpRespBytes, userinfo)
}
