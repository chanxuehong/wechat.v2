// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package oauth2

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/chanxuehong/wechat/corp"
)

// 构造获取code的URL.
//  corpId:      企业的CorpID
//  redirectURL: 授权后重定向的回调链接地址, 员工点击后, 页面将跳转至
//               redirect_uri/?code=CODE&state=STATE, 企业可根据code参数获得员工的userid.
//  scope:       应用授权作用域, 此时固定为: snsapi_base
//  state:       重定向后会带上state参数, 企业可以填写a-zA-Z0-9的参数值, 长度不可超过128个字节
func AuthCodeURL(corpId, redirectURL, scope, state string) string {
	return "https://open.weixin.qq.com/connect/oauth2/authorize" +
		"?appid=" + url.QueryEscape(corpId) +
		"&redirect_uri=" + url.QueryEscape(redirectURL) +
		"&response_type=code&scope=" + url.QueryEscape(scope) +
		"&state=" + url.QueryEscape(state) +
		"#wechat_redirect"
}

type Client corp.Client

func NewClient(srv corp.AccessTokenServer, clt *http.Client) *Client {
	return (*Client)(corp.NewClient(srv, clt))
}

type UserInfo struct {
	UserId   string `json:"UserId"`   // 员工UserID
	DeviceId string `json:"DeviceId"` // 手机设备号(由微信在安装时随机生成)
}

// 根据code获取成员信息.
//  agentId: 跳转链接时所在的企业应用ID
//  code:    通过员工授权获取到的code, 每次员工授权带上的code将不一样,
//           code只能使用一次, 5分钟未被使用自动过期
func (clt *Client) UserInfo(agentId int64, code string) (info *UserInfo, err error) {
	var result struct {
		corp.Error
		UserInfo
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?agentid=" +
		strconv.FormatInt(agentId, 10) + "&code=" + url.QueryEscape(code) +
		"&access_token="
	if err = ((*corp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.UserInfo
	return
}
