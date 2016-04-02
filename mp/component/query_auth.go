// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package component

import (
	"github.com/chanxuehong/wechat/mp"
)

// 權限集
type Function struct {
	ScopeCategory struct {
		Id int64 `json:"id"`
	} `json:"funcscope_category"`
}

type AuthorizationInfo struct {
	AuthorizerAccessTokenInfo
	AuthorizerAppId string     `json:"authorizer_appid"`
	FunctionInfo    []Function `json:"func_info"`
}

// 使用授权码换取公众号的授权信息.
func (clt *Client) QueryAuth(authCode string) (info *AuthorizationInfo, err error) {
	request := struct {
		ComponentAppId string `json:"component_appid"`
		AuthCode       string `json:"authorization_code"`
	}{
		ComponentAppId: clt.AppId,
		AuthCode:       authCode,
	}

	var result struct {
		mp.Error
		AuthorizationInfo `json:"authorization_info"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.AuthorizationInfo
	return
}
