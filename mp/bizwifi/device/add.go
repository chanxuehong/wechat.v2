// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package device

import (
	"github.com/chanxuehong/wechat/mp"
)

type AddParameters struct {
	ShopId   int64  `json:"shop_id"`  // 必须, 门店ID
	SSID     string `json:"ssid"`     // 必须, 无线网络设备的ssid。非认证公众号添加的ssid必需是“WX”开头(“WX”为大写字母)，认证公众号和第三方平台无此限制；所有ssid均不能包含中文字符
	Password string `json:"password"` // 必须, 无线网络设备的密码，大于8个字符，不能包含中文字符
	BSSID    string `json:"bssid"`    // 必须, 无线网络设备无线mac地址，格式冒号分隔，字符长度17个，并且字母小写，例如：00:1f:7a:ad:5c:a8
}

// 添加设备
func Add(clt *mp.Client, para *AddParameters) (err error) {
	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/bizwifi/device/add?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
