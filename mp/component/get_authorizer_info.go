// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package component

import (
	"github.com/chanxuehong/wechat/mp"
)

type AuthorizerInfo struct {
	NickName        string `json:"nick_name"`
	HeadImage       string `json:"head_img"`
	ServiceTypeInfo struct {
		Id int64 `json:"id"`
	} `json:"service_type_info"`
	VerifyTypeInfo struct {
		Id int64 `json:"id"`
	} `json:"verify_type_info"`
	UserName string `json:"user_name"`
	Alias    string `json:"alias"`
}

type AuthorizerInfoEx struct {
	AuthorizerInfo    AuthorizerInfo    `json:"authorizer_info"`
	QrCodeURL         string            `json:"qrcode_url"`
	AuthorizationInfo AuthorizationInfo `json:"authorization_info"`
}

// 获取授权方的账户信息.
func (clt *Client) GetAuthorizerInfo(authorizerAppId string) (info *AuthorizerInfoEx, err error) {
	request := struct {
		ComponentAppId  string `json:"component_appid"`
		AuthorizerAppId string `json:"authorizer_appid"`
	}{
		ComponentAppId:  clt.AppId,
		AuthorizerAppId: authorizerAppId,
	}

	var result struct {
		mp.Error
		AuthorizerInfoEx
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?component_access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.AuthorizerInfoEx
	return
}
