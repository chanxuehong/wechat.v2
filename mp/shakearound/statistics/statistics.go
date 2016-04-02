// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com), magicshui(shuiyuzhe@gmail.com), Harry Rong(harrykobe@gmail.com)

package statistics

import (
	"github.com/chanxuehong/wechat/mp/shakearound/device"
)

type StatisticsBase struct {
	Ftime   int64 `json:"ftime"`    // 当天0点对应的时间戳
	ClickPV int   `json:"click_pv"` // 点击摇周边消息的次数
	ClickUV int   `json:"click_uv"` // 点击摇周边消息的人数
	ShakePV int   `json:"shake_pv"` // 摇周边的次数
	ShakeUV int   `json:"shake_uv"` // 摇周边的人数
}

type DeviceStatistics struct {
	device.DeviceBase
	StatisticsBase
}

type PageStatistics struct {
	PageId int64 `json:"page_id"`
	StatisticsBase
}
