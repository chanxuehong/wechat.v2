// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package oauth2

import (
	"net/url"
	"strings"
)

type Config interface {
	AuthCodeURL(state string, redirectURIExt url.Values) string // 请求用户授权的地址, 获取code; redirectURIExt 用于扩展回调地址的参数
	ExchangeTokenURL(code string) string                        // 通过code换取access_token的地址
	RefreshTokenURL(refreshToken string) string                 // 刷新access_token的地址
	UserInfoURL(accessToken, openId, lang string) string        // 获取用户信息的地址
}

var _ Config = (*OAuth2Config)(nil)

type OAuth2Config struct {
	AppId     string
	AppSecret string

	// 用户授权后跳转的目的地址
	// 用户授权后跳转到 RedirectURI?code=CODE&state=STATE
	// 用户禁止授权跳转到 RedirectURI?state=STATE
	RedirectURI string

	// 应用授权作用域, snsapi_base, snsapi_userinfo
	Scopes []string
}

func NewOAuth2Config(AppId, AppSecret, RedirectURI string, Scope ...string) *OAuth2Config {
	return &OAuth2Config{
		AppId:       AppId,
		AppSecret:   AppSecret,
		RedirectURI: RedirectURI,
		Scopes:      Scope,
	}
}

func (cfg *OAuth2Config) AuthCodeURL(state string, redirectURIExt url.Values) string {
	return AuthCodeURL(cfg.AppId, cfg.RedirectURI, strings.Join(cfg.Scopes, ","), state, redirectURIExt)
}

func AuthCodeURL(appId, redirectURI, scope, state string, redirectURIExt url.Values) string {
	if redirectURIExt != nil {
		if strings.Contains(redirectURI, "?") {
			redirectURI += "&" + redirectURIExt.Encode()
		} else {
			redirectURI += "?" + redirectURIExt.Encode()
		}
	}

	return "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + url.QueryEscape(appId) +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&response_type=code&scope=" + url.QueryEscape(scope) +
		"&state=" + url.QueryEscape(state) +
		"#wechat_redirect"
}

func (cfg *OAuth2Config) ExchangeTokenURL(code string) string {
	return "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + url.QueryEscape(cfg.AppId) +
		"&secret=" + url.QueryEscape(cfg.AppSecret) +
		"&grant_type=authorization_code&code=" + url.QueryEscape(code)
}

func (cfg *OAuth2Config) RefreshTokenURL(refreshToken string) string {
	return "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=" + url.QueryEscape(cfg.AppId) +
		"&grant_type=refresh_token&refresh_token=" + url.QueryEscape(refreshToken)
}

func (cfg *OAuth2Config) UserInfoURL(accessToken, openId, lang string) string {
	return "https://api.weixin.qq.com/sns/userinfo?access_token=" + url.QueryEscape(accessToken) +
		"&openid=" + url.QueryEscape(openId) +
		"&lang=" + url.QueryEscape(lang)
}
