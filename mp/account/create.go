// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package account

import (
	"errors"

	"github.com/chanxuehong/wechat/mp"
)

const (
	TemporaryQRCodeExpireSecondsLimit = 604800 // 临时二维码 expire_seconds 限制
	PermanentQRCodeSceneIdLimit       = 100000 // 永久二维码 scene_id 限制
)

// 永久二维码
type PermanentQRCode struct {
	// 下面两个字段同时只有一个有效, 非zero值表示有效.
	SceneId     uint32 `json:"scene_id,omitempty"`  // 场景值ID, 临时二维码时为32位非0整型, 永久二维码时最大值为100000(目前参数只支持1--100000)
	SceneString string `json:"scene_str,omitempty"` // 场景值ID(字符串形式的ID), 字符串类型, 长度限制为1到64, 仅永久二维码支持此字段

	Ticket string `json:"ticket"` // 获取的二维码ticket, 凭借此ticket可以在有效时间内换取二维码.
	URL    string `json:"url"`    // 二维码图片解析后的地址, 开发者可根据该地址自行生成需要的二维码图片
}

// 临时二维码
type TemporaryQRCode struct {
	ExpireSeconds int `json:"expire_seconds,omitempty"` // 二维码的有效时间, 以秒为单位. 最大不超过 604800.
	PermanentQRCode
}

// 创建临时二维码
//  SceneId:       场景值ID, 为32位非0整型
//  ExpireSeconds: 二维码有效时间, 以秒为单位.  最大不超过 604800.
func (clt *Client) CreateTemporaryQRCode(SceneId uint32, ExpireSeconds int) (qrcode *TemporaryQRCode, err error) {
	if SceneId == 0 {
		err = errors.New("SceneId should be greater than 0")
		return
	}
	if ExpireSeconds <= 0 {
		err = errors.New("ExpireSeconds should be greater than 0")
		return
	}
	var request struct {
		ExpireSeconds int    `json:"expire_seconds"`
		ActionName    string `json:"action_name"`
		ActionInfo    struct {
			Scene struct {
				SceneId uint32 `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ExpireSeconds = ExpireSeconds
	request.ActionName = "QR_SCENE"
	request.ActionInfo.Scene.SceneId = SceneId

	var result struct {
		mp.Error
		TemporaryQRCode
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	result.TemporaryQRCode.SceneId = SceneId
	qrcode = &result.TemporaryQRCode
	return
}

// 创建永久二维码
//  SceneId: 场景值ID, 目前参数只支持1--100000
func (clt *Client) CreatePermanentQRCode(SceneId uint32) (qrcode *PermanentQRCode, err error) {
	if SceneId == 0 {
		err = errors.New("SceneId should be greater than 0")
		return
	}
	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneId uint32 `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ActionName = "QR_LIMIT_SCENE"
	request.ActionInfo.Scene.SceneId = SceneId

	var result struct {
		mp.Error
		PermanentQRCode
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	result.PermanentQRCode.SceneId = SceneId
	qrcode = &result.PermanentQRCode
	return
}

// 创建永久二维码
//  SceneString: 场景值ID(字符串形式的ID), 字符串类型, 长度限制为1到64
func (clt *Client) CreatePermanentQRCodeWithSceneString(SceneString string) (qrcode *PermanentQRCode, err error) {
	if SceneString == "" {
		err = errors.New("SceneString should not be empty")
		return
	}
	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneString string `json:"scene_str"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ActionName = "QR_LIMIT_STR_SCENE"
	request.ActionInfo.Scene.SceneString = SceneString

	var result struct {
		mp.Error
		PermanentQRCode
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	result.PermanentQRCode.SceneString = SceneString
	qrcode = &result.PermanentQRCode
	return
}
