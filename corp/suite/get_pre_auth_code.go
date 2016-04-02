// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package suite

import (
	"github.com/chanxuehong/wechat/corp"
)

type PreAuthCode struct {
	Value     string `json:"pre_auth_code"`
	ExpiresIn int64  `json:"expires_in"`
}

// 获取预授权码.
//  appIdList: 应用id, 本参数选填, 表示用户能对本套件内的哪些应用授权, 不填时默认用户有全部授权权限
func (clt *Client) GetPreAuthCode(appIdList []int64) (code *PreAuthCode, err error) {
	request := struct {
		SuiteId   string  `json:"suite_id"`
		AppIdList []int64 `json:"appid,omitempty"`
	}{
		SuiteId:   clt.SuiteId,
		AppIdList: appIdList,
	}

	var result struct {
		corp.Error
		PreAuthCode
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/service/get_pre_auth_code?suite_access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	code = &result.PreAuthCode
	return
}
