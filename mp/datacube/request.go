// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package datacube

import (
	"time"
)

// 获取统计数据通用的请求结构.
type Request struct {
	// 获取数据的起始日期, YYYY-MM-DD 格式;
	// begin_date 和 end_date 的差值需小于"最大时间跨度"(比如最大时间跨度为1时,
	// begin_date 和 end_date 的差值只能为0, 才能小于1), 否则会报错.
	BeginDate string `json:"begin_date,omitempty"`

	// 获取数据的结束日期, YYYY-MM-DD 格式;
	// end_date 允许设置的最大值为昨日.
	EndDate string `json:"end_date,omitempty"`
}

// NewRequest 创建一个 Request.
//  请注意 BeginDate, EndDate 的 Location.
func NewRequest(BeginDate, EndDate time.Time) *Request {
	return &Request{
		BeginDate: BeginDate.Format("2006-01-02"),
		EndDate:   EndDate.Format("2006-01-02"),
	}
}
