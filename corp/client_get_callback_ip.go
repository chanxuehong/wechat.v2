// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package corp

// 获取微信服务器的ip段
func (clt *Client) GetCallbackIP() (ipList []string, err error) {
	var result struct {
		Error
		IPList []string `json:"ip_list"`
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/getcallbackip?access_token="
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
