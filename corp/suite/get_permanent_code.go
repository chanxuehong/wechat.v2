// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package suite

import (
	"github.com/chanxuehong/wechat/corp"
)

type PermanentCodeInfo struct {
	CorpAccessTokenInfo
	PermanentCode string       `json:"permanent_code"`
	AuthCorpInfo  AuthCorpInfo `json:"auth_corp_info"`
	AuthInfo      AuthInfo     `json:"auth_info"`
	AuthUserInfo  AuthUserInfo `json:"auth_user_info"`
}

type AuthUserInfo struct {
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

// 获取企业号的永久授权码
//  authCode: 临时授权码会在授权成功时附加在redirect_uri中跳转回应用提供商网站.
func (clt *Client) GetPermanentCode(authCode string) (info *PermanentCodeInfo, err error) {
	request := struct {
		SuiteId  string `json:"suite_id"`
		AuthCode string `json:"auth_code"`
	}{
		SuiteId:  clt.SuiteId,
		AuthCode: authCode,
	}

	var result struct {
		corp.Error
		PermanentCodeInfo
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/service/get_permanent_code?suite_access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.PermanentCodeInfo
	return
}
