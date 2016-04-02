// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package oauth2

import (
	"net/url"
)

type Config interface {
	AuthCodeURL(state string, redirectURIExt url.Values) string // 请求用户授权的地址, 获取code; redirectURIExt 用于扩展回调地址的参数
	ExchangeTokenURL(code string) string                        // 通过code换取access_token的地址
	RefreshTokenURL(refreshToken string) string                 // 刷新access_token的地址
	UserInfoURL(accessToken, openId, lang string) string        // 获取用户信息的地址, lang可以为空值
}
