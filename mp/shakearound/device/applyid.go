// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com), magicshui(shuiyuzhe@gmail.com), Harry Rong(harrykobe@gmail.com)

package device

import (
	"github.com/chanxuehong/wechat/mp"
)

type ApplyIdParameters struct {
	Quantity    int    `json:"quantity"`          // 必须, 申请的设备ID的数量，单次新增设备超过500个，需走人工审核流程
	ApplyReason string `json:"apply_reason"`      // 必须, 申请理由，不超过100个字
	Comment     string `json:"comment,omitempty"` // 可选, 备注，不超过15个汉字或30个英文字母
	PoiId       *int64 `json:"poi_id,omitempty"`  // 可选, 设备关联的门店ID，关联门店后，在门店1KM的范围内有优先摇出信息的机会。
}

type ApplyIdResult struct {
	ApplyId      int64  `json:"apply_id"`      // 申请的批次ID，可用在“查询设备列表”接口按批次查询本次申请成功的设备ID。
	AuditStatus  int    `json:"audit_status"`  // 审核状态。0：审核未通过、1：审核中、2：审核已通过；若单次申请的设备ID数量小于等于500个，系统会进行快速审核；若单次申请的设备ID数量大于500个，会在三个工作日内完成审核
	AuditComment string `json:"audit_comment"` // 审核备注，包括审核不通过的原因
}

// 申请设备ID
func ApplyId(clt *mp.Client, para *ApplyIdParameters) (rslt *ApplyIdResult, err error) {
	var result struct {
		mp.Error
		ApplyIdResult `json:"data"`
	}

	incompleteURL := "https://api.weixin.qq.com/shakearound/device/applyid?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}

	rslt = &result.ApplyIdResult
	return
}
