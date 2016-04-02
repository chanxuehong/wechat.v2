// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package device

import (
	"github.com/chanxuehong/wechat/mp"
)

// 删除设备
func Delete(clt *mp.Client, bssid string) (err error) {
	request := struct {
		BSSID string `json:"bssid"`
	}{
		BSSID: bssid,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/bizwifi/device/delete?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
