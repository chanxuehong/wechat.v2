// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package dkf

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/chanxuehong/wechat/mp"
)

// 客服基本信息
type KfInfo struct {
	Id           string `json:"kf_id"`         // 客服工号
	Account      string `json:"kf_account"`    // 完整客服账号, 格式为: 账号前缀@公众号微信号
	Nickname     string `json:"kf_nick"`       // 客服昵称
	HeadImageURL string `json:"kf_headimgurl"` // 客服头像
}

var ErrNoHeadImage = errors.New("没有头像")

// 获取客服图像的大小, 如果客服没有图像则返回 ErrNoHeadImage 错误.
func (info *KfInfo) HeadImageSize() (size int, err error) {
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

// 获取客服基本信息.
func KfList(clt *mp.Client) (kfList []KfInfo, err error) {
	var result struct {
		mp.Error
		KfList []KfInfo `json:"kf_list"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/customservice/getkflist?access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	kfList = result.KfList
	return
}

const (
	OnlineKfInfoStatusPC          = 1
	OnlineKfInfoStatusMobile      = 2
	OnlineKfInfoStatusPCAndMobile = 3
)

// 在线客服接待信息
type OnlineKfInfo struct {
	Id                  string `json:"kf_id"`         // 客服工号
	Account             string `json:"kf_account"`    // 完整客服账号, 格式为: 账号前缀@公众号微信号
	Status              int    `json:"status"`        // 客服在线状态 1: pc在线, 2: 手机在线. 若pc和手机同时在线则为 1+2=3
	AutoAcceptThreshold int    `json:"auto_accept"`   // 客服设置的最大自动接入数
	AcceptingNumber     int    `json:"accepted_case"` // 客服当前正在接待的会话数
}

// 获取在线客服接待信息.
func OnlineKfList(clt *mp.Client) (kfList []OnlineKfInfo, err error) {
	var result struct {
		mp.Error
		OnlineKfInfoList []OnlineKfInfo `json:"kf_online_list"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/customservice/getonlinekflist?access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	kfList = result.OnlineKfInfoList
	return
}
