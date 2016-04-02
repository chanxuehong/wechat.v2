// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package user

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/chanxuehong/wechat/mp"
)

const (
	Language_zh_CN = "zh_CN" // 简体中文
	Language_zh_TW = "zh_TW" // 繁体中文
	Language_en    = "en"    // 英文
)

const (
	SexUnknown = 0 // 未知
	SexMale    = 1 // 男性
	SexFemale  = 2 // 女性
)

type UserInfo struct {
	IsSubscriber int    `json:"subscribe"` // 用户是否订阅该公众号标识, 值为0时, 代表此用户没有关注该公众号, 拉取不到其余信息
	OpenId       string `json:"openid"`    // 用户的标识, 对当前公众号唯一
	Nickname     string `json:"nickname"`  // 用户的昵称
	Sex          int    `json:"sex"`       // 用户的性别, 值为1时是男性, 值为2时是女性, 值为0时是未知
	Language     string `json:"language"`  // 用户的语言, zh_CN, zh_TW, en
	City         string `json:"city"`      // 用户所在城市
	Province     string `json:"province"`  // 用户所在省份
	Country      string `json:"country"`   // 用户所在国家

	// 用户头像, 最后一个数值代表正方形头像大小(有0, 46, 64, 96, 132数值可选, 0代表640*640正方形头像),
	// 用户没有头像时该项为空
	HeadImageURL string `json:"headimgurl"`

	// 用户关注时间, 为时间戳. 如果用户曾多次关注, 则取最后关注时间
	SubscribeTime int64 `json:"subscribe_time"`

	// 只有在用户将公众号绑定到微信开放平台帐号后, 才会出现该字段.
	UnionId string `json:"unionid"`

	Remark  string `json:"remark"`  // 公众号运营者对粉丝的备注, 公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	GroupId int64  `json:"groupid"` // 用户所在的分组ID
}

var ErrNoHeadImage = errors.New("没有头像")

// 获取用户图像的大小, 如果用户没有图像则返回 ErrNoHeadImage 错误.
func (info *UserInfo) HeadImageSize() (size int, err error) {
	HeadImageURL := info.HeadImageURL
	if HeadImageURL == "" {
		err = ErrNoHeadImage
		return
	}

	lastSlashIndex := strings.LastIndex(HeadImageURL, "/")
	if lastSlashIndex == -1 {
		err = fmt.Errorf("invalid HeadImageURL: %s", HeadImageURL)
		return
	}
	HeadImageIndex := lastSlashIndex + 1
	if HeadImageIndex == len(HeadImageURL) {
		err = fmt.Errorf("invalid HeadImageURL: %s", HeadImageURL)
		return
	}

	sizeStr := HeadImageURL[HeadImageIndex:]

	size, err = strconv.Atoi(sizeStr)
	if err != nil {
		err = fmt.Errorf("invalid HeadImageURL: %s", HeadImageURL)
		return
	}

	if size == 0 {
		size = 640
	}
	return
}

// 获取用户基本信息.
//  注意:
//  1. 需要判断返回的 UserInfo.IsSubscriber 是否等于 1 还是 0
//  2. lang 可以是 zh_CN, zh_TW, en, 如果留空 "" 则默认为 zh_CN
func (clt *Client) UserInfo(openId string, lang string) (userinfo *UserInfo, err error) {
	if openId == "" {
		err = errors.New("empty openId")
		return
	}

	switch lang {
	case "":
		lang = Language_zh_CN
	case Language_zh_CN, Language_zh_TW, Language_en:
	default:
		lang = Language_zh_CN
	}

	var result struct {
		mp.Error
		UserInfo
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/user/info?openid=" + url.QueryEscape(openId) +
		"&lang=" + url.QueryEscape(lang) + "&access_token="
	if err = ((*mp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	userinfo = &result.UserInfo
	return
}

type UserInfoBatchGetRequestItem struct {
	OpenId   string `json:"openid"`
	Language string `json:"lang,omitempty"`
}

// 创建 []UserInfoBatchGetRequestItem
//  lang 的取值可以为 "", Language_zh_CN, Language_zh_TW, Language_en
func NewUserInfoBatchGetRequest(openIdList []string, lang string) (ret []UserInfoBatchGetRequestItem) {
	switch lang {
	case "", Language_zh_CN, Language_zh_TW, Language_en:
	default:
		lang = ""
	}

	ret = make([]UserInfoBatchGetRequestItem, len(openIdList))
	for i := 0; i < len(openIdList); i++ {
		ret[i].OpenId = openIdList[i]
		ret[i].Language = lang
	}
	return
}

// 批量获取用户基本信息
//  注意: 需要对返回的 UserInfoList 的每个 UserInfo.IsSubscriber 做判断
func (clt *Client) UserInfoBatchGet(req []UserInfoBatchGetRequestItem) (UserInfoList []UserInfo, err error) {
	if len(req) <= 0 {
		err = errors.New("empty request")
		return
	}

	var request = struct {
		UserList []UserInfoBatchGetRequestItem `json:"user_list,omitempty"`
	}{
		UserList: req,
	}

	var result struct {
		mp.Error
		UserInfoList []UserInfo `json:"user_info_list"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/user/info/batchget?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	UserInfoList = result.UserInfoList
	return
}

// 开发者可以通过该接口对指定用户设置备注名.
func (clt *Client) UserUpdateRemark(openId, remark string) (err error) {
	var request = struct {
		OpenId string `json:"openid"`
		Remark string `json:"remark"`
	}{
		OpenId: openId,
		Remark: remark,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/user/info/updateremark?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 获取关注者列表返回的数据结构
type UserListResult struct {
	TotalCount int `json:"total"` // 关注该公众账号的总用户数
	GotCount   int `json:"count"` // 拉取的 OPENID 个数, 最大值为10000

	Data struct {
		OpenIdList []string `json:"openid,omitempty"`
	} `json:"data"` // 列表数据, OPENID 的列表

	// 拉取列表的后一个用户的OPENID, 如果 next_openid == "" 则表示没有了用户数据
	NextOpenId string `json:"next_openid"`
}

// 获取关注者列表.
//  NOTE:
//  1. 每次最多能获取 10000 个用户, 可以多次指定 NextOpenId 来获取以满足需求, 如果 NextOpenId == "" 则表示从头获取
//  2. 目前微信返回的数据并不包括 NextOpenId 本身, 是从 NextOpenId 下一个用户开始的, 和微信文档描述不一样!!!
func (clt *Client) UserList(NextOpenId string) (rslt *UserListResult, err error) {
	var result struct {
		mp.Error
		UserListResult
	}

	var incompleteURL string
	if NextOpenId == "" {
		incompleteURL = "https://api.weixin.qq.com/cgi-bin/user/get?access_token="
	} else {
		incompleteURL = "https://api.weixin.qq.com/cgi-bin/user/get?next_openid=" + url.QueryEscape(NextOpenId) + "&access_token="
	}

	if err = ((*mp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}

	rslt = &result.UserListResult
	return
}
