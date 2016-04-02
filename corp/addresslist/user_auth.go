// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package addresslist

import (
	"net/url"

	"github.com/chanxuehong/wechat/corp"
)

// 二次验证
//  企业在开启二次验证时, 必须填写企业二次验证页面的url, 此url的域名必须设置为企业小助手的可信域名.
//  当员工绑定通讯录中的帐号后, 会收到一条图文消息, 引导员工到企业的验证页面验证身份.
//  在跳转到企业的验证页面时, 会带上如下参数: code=CODE&state=STATE, 企业可以调用oauth2接口,
//  根据code和agentid获取员工的userid.
//
//  企业在员工验证成功后, 调用如下接口即可让员工关注成功.
func (clt *Client) UserAuthSuccess(userId string) (err error) {
	var result corp.Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/user/authsucc?userid=" +
		url.QueryEscape(userId) + "&access_token="
	if err = ((*corp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result
		return
	}
	return
}
