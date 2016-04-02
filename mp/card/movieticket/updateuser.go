// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package movieticket

import (
	"github.com/chanxuehong/wechat/mp"
)

type UpdateUserParameters struct {
	Code   string `json:"code"`              // 必须; 卡券Code码。
	CardId string `json:"card_id,omitempty"` // 可选; 要更新门票序列号所述的card_id，生成券时use_custom_code填写true时必填。

	TicketClass   string `json:"ticket_class,omitempty"`   // 必须; 电影票的类别，如2D、3D。
	ShowTime      int64  `json:"show_time,omitempty"`      // 必须; 电影的放映时间，Unix时间戳格式。
	Duration      int    `json:"duration,omitempty"`       // 必须; 放映时长，填写整数。
	ScreeningRoom string `json:"screening_room,omitempty"` // 可选; 该场电影的影厅信息。
	SeatNumber    string `json:"seat_number,omitempty"`    // 可选; 座位号。
}

// 更新电影票
func UpdateUser(clt *mp.Client, para *UpdateUserParameters) (err error) {
	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/card/movieticket/updateuser?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
