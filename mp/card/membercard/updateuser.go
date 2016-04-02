// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package membercard

import (
	"github.com/chanxuehong/wechat/mp"
)

type UpdateUserParameters struct {
	Code   string `json:"code"`              // 必须; 要消耗的序列号.
	CardId string `json:"card_id,omitempty"` // 可选; 要消耗序列号所述的card_id. 自定义code 的会员卡必填

	AddBonus      int    `json:"add_bonus,omitempty"`      // 必须; 需要变更的积分，扣除积分用“-“表示。
	RecordBonus   string `json:"record_bonus,omitempty"`   // 可选; 商家自定义积分消耗记录，不超过14个汉字。
	AddBalance    int    `json:"add_balance,omitempty"`    // 可选; 需要变更的余额，扣除金额用“-”表示。单位为分。
	RecordBalance string `json:"record_balance,omitempty"` // 可选; 商家自定义金额消耗记录，不超过14个汉字。

	CustomFieldValue1 string `json:"custom_field_value1,omitempty"` // 可选, 创建时字段custom_field1定义类型的初始值，限制为4个汉字，12字节。
	CustomFieldValue2 string `json:"custom_field_value2,omitempty"` // 可选, 创建时字段custom_field2定义类型的初始值，限制为4个汉字，12字节。
	CustomFieldValue3 string `json:"custom_field_value3,omitempty"` // 可选, 创建时字段custom_field3定义类型的初始值，限制为4个汉字，12字节。
}

type UpdateUserResult struct {
	ResultBonus   int    `json:"result_bonus"`   // 当前用户积分总额。
	ResultBalance int    `json:"result_balance"` // 当前用户预存总金额。
	OpenId        string `json:"openid"`         // 用户openid。
}

// 更新会员信息
func UpdateUser(clt *mp.Client, para *UpdateUserParameters) (rslt *UpdateUserResult, err error) {
	var result struct {
		mp.Error
		UpdateUserResult
	}

	incompleteURL := "https://api.weixin.qq.com/card/membercard/updateuser?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	rslt = &result.UpdateUserResult
	return
}
