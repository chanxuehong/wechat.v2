// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mp

// 获取微信服务器IP地址.
//  如果公众号基于安全等考虑, 需要获知微信服务器的IP地址列表, 以便进行相关限制,
//  可以通过该接口获得微信服务器IP地址列表.
func (clt *Client) GetCallbackIP() (ipList []string, err error) {
	var result struct {
		Error
		IPList []string `json:"ip_list"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	ipList = result.IPList
	return
}
