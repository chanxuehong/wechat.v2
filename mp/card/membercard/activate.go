// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package membercard

import (
	"github.com/chanxuehong/wechat/mp"
)

type ActivateParameters struct {
	Code   string `json:"code"`              // 必填, 创建会员卡时获取的初始code。
	CardId string `json:"card_id,omitempty"` // 可选; 卡券ID. 自定义code 的会员卡必填card_id, 非自定义code 的会员卡不必填.

	MembershipNumber string `json:"membership_number,omitempty"` // 必填, 会员卡编号，由开发者填入，作为序列号显示在用户的卡包里。可与Code码保持等值。

	ActivateBeginTime int64 `json:"activate_begin_time,omitempty"` // 可选; 激活后的有效起始时间。若不填写默认以创建时的 data_info 为准。Unix时间戳格式。
	ActivateEndTime   int64 `json:"activate_end_time,omitempty"`   // 可选; 激活后的有效截至时间。若不填写默认以创建时的 data_info 为准。Unix时间戳格式。

	InitBonus   *int `json:"init_bonus,omitempty"`   // 可选; 初始积分, 不填为0
	InitBalance *int `json:"init_balance,omitempty"` // 可选; 初始余额, 不填为0

	InitCustomFieldValue1 string `json:"init_custom_field_value1,omitempty"` // 可选, 创建时字段custom_field1定义类型的初始值，限制为4个汉字，12字节。
	InitCustomFieldValue2 string `json:"init_custom_field_value2,omitempty"` // 可选, 创建时字段custom_field2定义类型的初始值，限制为4个汉字，12字节。
	InitCustomFieldValue3 string `json:"init_custom_field_value3,omitempty"` // 可选, 创建时字段custom_field3定义类型的初始值，限制为4个汉字，12字节。
}

// 激活/绑定会员卡
func Activate(clt *mp.Client, para *ActivateParameters) (err error) {
	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/card/membercard/activate?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
