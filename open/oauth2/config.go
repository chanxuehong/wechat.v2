// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package oauth2

import (
	"net/url"
	"strings"

	"github.com/chanxuehong/wechat/oauth2"
)

var _ oauth2.Config = (*Config)(nil)

type Config struct {
	AppId     string
	AppSecret string

	// 用户授权后跳转的目的地址
	// 用户授权后跳转到 RedirectURI?code=CODE&state=STATE
	// 用户禁止授权跳转到 RedirectURI?state=STATE
	RedirectURI string

	// 应用授权作用域, snsapi_base, snsapi_userinfo
	Scopes []string
}

func NewConfig(AppId, AppSecret, RedirectURI string, Scope ...string) *Config {
	return &Config{
		AppId:       AppId,
		AppSecret:   AppSecret,
		RedirectURI: RedirectURI,
		Scopes:      Scope,
	}
}

func (cfg *Config) AuthCodeURL(state string, redirectURIExt url.Values) string {
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

	return "https://open.weixin.qq.com/connect/qrconnect?appid=" + url.QueryEscape(appId) +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&response_type=code&scope=" + url.QueryEscape(scope) +
		"&state=" + url.QueryEscape(state) +
		"#wechat_redirect"
}

func (cfg *Config) ExchangeTokenURL(code string) string {
	return "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + url.QueryEscape(cfg.AppId) +
		"&secret=" + url.QueryEscape(cfg.AppSecret) +
		"&grant_type=authorization_code&code=" + url.QueryEscape(code)
}

func (cfg *Config) RefreshTokenURL(refreshToken string) string {
	return "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=" + url.QueryEscape(cfg.AppId) +
		"&grant_type=refresh_token&refresh_token=" + url.QueryEscape(refreshToken)
}

func (cfg *Config) UserInfoURL(accessToken, openId, lang string) string {
	if lang == "" {
		return "https://api.weixin.qq.com/sns/userinfo?access_token=" + url.QueryEscape(accessToken) +
			"&openid=" + url.QueryEscape(openId)
	}
	return "https://api.weixin.qq.com/sns/userinfo?access_token=" + url.QueryEscape(accessToken) +
		"&openid=" + url.QueryEscape(openId) +
		"&lang=" + url.QueryEscape(lang)
}
