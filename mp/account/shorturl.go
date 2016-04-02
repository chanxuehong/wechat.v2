// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package account

import (
	"github.com/chanxuehong/wechat/mp"
)

// 将一条长链接转成短链接.
func (clt *Client) ShortURL(longURL string) (shortURL string, err error) {
	var request = struct {
		Action  string `json:"action"`
		LongURL string `json:"long_url"`
	}{
		Action:  "long2short",
		LongURL: longURL,
	}

	var result struct {
		mp.Error
		ShortURL string `json:"short_url"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/shorturl?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	shortURL = result.ShortURL
	return
}
