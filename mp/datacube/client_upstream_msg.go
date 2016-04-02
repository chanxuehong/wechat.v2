// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package datacube

import (
	"errors"

	"github.com/chanxuehong/wechat/mp"
)

// 消息发送概况数据
type UpstreamMsgData struct {
	RefDate    string `json:"ref_date"`    // 数据的日期, YYYY-MM-DD 格式
	UserSource int    `json:"user_source"` // 返回的 json 有这个字段, 文档中没有, 都是 0 值, 可能没有实际意义!!!

	// 消息类型, 代表含义如下:
	// 1代表文字
	// 2代表图片
	// 3代表语音
	// 4代表视频
	// 6代表第三方应用消息(链接消息)
	MsgType  int `json:"msg_type"`
	MsgUser  int `json:"msg_user"`  // 上行发送了(向公众号发送了)消息的用户数
	MsgCount int `json:"msg_count"` // 上行发送了消息的消息总数
}

// 获取消息发送概况数据.
func (clt *Client) GetUpstreamMsg(req *Request) (list []UpstreamMsgData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UpstreamMsgData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsg?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, req, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息分送分时数据
type UpstreamMsgHourData struct {
	RefHour int `json:"ref_hour"` // 数据的小时, 包括从000到2300, 分别代表的是[000,100)到[2300,2400), 即每日的第1小时和最后1小时
	UpstreamMsgData
}

// 获取消息分送分时数据.
func (clt *Client) GetUpstreamMsgHour(req *Request) (list []UpstreamMsgHourData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UpstreamMsgHourData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsghour?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, req, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息发送周数据
type UpstreamMsgWeekData UpstreamMsgData

// 获取消息发送周数据.
func (clt *Client) GetUpstreamMsgWeek(req *Request) (list []UpstreamMsgWeekData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UpstreamMsgWeekData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsgweek?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, req, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息发送月数据
type UpstreamMsgMonthData UpstreamMsgData

// 获取消息发送月数据.
func (clt *Client) GetUpstreamMsgMonth(req *Request) (list []UpstreamMsgMonthData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UpstreamMsgMonthData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsgmonth?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, req, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息发送分布数据
type UpstreamMsgDistData struct {
	RefDate       string `json:"ref_date"`       // 数据的日期, YYYY-MM-DD 格式
	UserSource    int    `json:"user_source"`    // 返回的 json 有这个字段, 文档中没有, 都是 0 值, 可能没有实际意义!!!
	CountInterval int    `json:"count_interval"` // 当日发送消息量分布的区间, 0代表 "0", 1代表"1-5", 2代表"6-10", 3代表"10次以上"
	MsgUser       int    `json:"msg_user"`       // 上行发送了(向公众号发送了)消息的用户数
}

// 获取消息发送分布数据.
func (clt *Client) GetUpstreamMsgDist(req *Request) (list []UpstreamMsgDistData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UpstreamMsgDistData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsgdist?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, req, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息发送分布周数据
type UpstreamMsgDistWeekData UpstreamMsgDistData

// 获取消息发送分布周数据.
func (clt *Client) GetUpstreamMsgDistWeek(req *Request) (list []UpstreamMsgDistWeekData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UpstreamMsgDistWeekData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsgdistweek?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, req, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息发送分布月数据
type UpstreamMsgDistMonthData UpstreamMsgDistData

// 获取消息发送分布月数据.
func (clt *Client) GetUpstreamMsgDistMonth(req *Request) (list []UpstreamMsgDistMonthData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UpstreamMsgDistMonthData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsgdistmonth?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, req, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}
