// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package addresslist

import (
	"github.com/chanxuehong/wechat/corp"
)

// 邀请成员关注
//  UserId:     必须, 用户的userid
//  InviteTips: 可选, 推送到微信上的提示语(只有认证号可以使用).
//              当使用微信推送时, 该字段默认为"请关注XXX企业号", 邮件邀请时, 该字段无效.
//  Type:       1:微信邀请 2.邮件邀请
func (clt *Client) InviteSend(UserId, InviteTips string) (Type int, err error) {
	var request = struct {
		UserId     string `json:"userid"`
		InviteTips string `json:"invite_tips,omitempty"`
	}{
		UserId:     UserId,
		InviteTips: InviteTips,
	}

	var result struct {
		corp.Error
		Type int `json:"type"`
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/invite/send?access_token="
	if err = ((*corp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	Type = result.Type
	return
}
