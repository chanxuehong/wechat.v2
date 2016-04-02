// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package component

import (
	"github.com/chanxuehong/wechat/mp"
)

// 设置授权方的选项信息.
func (clt *Client) SetAuthorizerOption(authorizerAppId, optionName, optionValue string) (err error) {
	request := struct {
		ComponentAppId  string `json:"component_appid"`
		AuthorizerAppId string `json:"authorizer_appid"`
		OptionName      string `json:"option_name"`
		OptionValue     string `json:"option_value"`
	}{
		ComponentAppId:  clt.AppId,
		AuthorizerAppId: authorizerAppId,
		OptionName:      optionName,
		OptionValue:     optionValue,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/component/api_set_authorizer_option?component_access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
