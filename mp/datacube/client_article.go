// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package datacube

import (
	"errors"

	"github.com/chanxuehong/wechat/mp"
)

type ArticleBaseData struct {
	IntPageReadUser  int `json:"int_page_read_user"`  // 图文页(点击群发图文卡片进入的页面)的阅读人数
	IntPageReadCount int `json:"int_page_read_count"` // 图文页的阅读次数
	OriPageReadUser  int `json:"ori_page_read_user"`  // 原文页(点击图文页"阅读原文"进入的页面)的阅读人数, 无原文页时此处数据为0
	OriPageReadCount int `json:"ori_page_read_count"` // 原文页的阅读次数
	ShareUser        int `json:"share_user"`          // 分享的人数
	ShareCount       int `json:"share_count"`         // 分享的次数
	AddToFavUser     int `json:"add_to_fav_user"`     // 收藏的人数
	AddToFavCount    int `json:"add_to_fav_count"`    // 收藏的次数
}

// 图文群发每日数据
type ArticleSummaryData struct {
	RefDate    string `json:"ref_date"`    // 数据的日期, YYYY-MM-DD 格式
	UserSource int    `json:"user_source"` // 返回的 json 有这个字段, 文档中没有, 都是 0 值, 可能没有实际意义!!!

	// 这里的msgid实际上是由msgid(图文消息id)和index(消息次序索引)组成,
	// 例如12003_3,  其中12003是msgid, 即一次群发的id消息的;
	// 3为index, 假设该次群发的图文消息共5个文章(因为可能为多图文),  3表示5个中的第3个
	MsgId string `json:"msgid"`
	Title string `json:"title"` // 图文消息的标题
	ArticleBaseData
}

// 获取图文群发每日数据.
func (clt *Client) GetArticleSummary(req *Request) (list []ArticleSummaryData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []ArticleSummaryData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getarticlesummary?access_token="
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

// 图文群发总数据
type ArticleTotalData struct {
	RefDate    string `json:"ref_date"`    // 数据的日期, YYYY-MM-DD 格式
	UserSource int    `json:"user_source"` // 返回的 json 有这个字段, 文档中没有, 都是 0 值, 可能没有实际意义!!!
	MsgId      string `json:"msgid"`       // 同 ArticleSummaryData.MsgId
	Title      string `json:"title"`
	Details    []struct {
		StatDate   string `json:"stat_date"`   // 统计的日期, 在getarticletotal接口中, ref_date指的是文章群发出日期,  而stat_date是数据统计日期
		TargetUser int    `json:"target_user"` // 送达人数, 一般约等于总粉丝数(需排除黑名单或其他异常情况下无法收到消息的粉丝)
		ArticleBaseData
	} `json:"details"`
}

// 获取图文群发总数据.
func (clt *Client) GetArticleTotal(req *Request) (list []ArticleTotalData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []ArticleTotalData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getarticletotal?access_token="
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

// 图文统计数据
type UserReadData struct {
	RefDate    string `json:"ref_date"` // 数据的日期, YYYY-MM-DD 格式
	UserSource int    `json:"user_source"`
	ArticleBaseData
}

// 获取图文统计数据.
func (clt *Client) GetUserRead(req *Request) (list []UserReadData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UserReadData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getuserread?access_token="
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

// 图文统计分时数据
type UserReadHourData struct {
	RefHour         int   `json:"ref_hour"`          // 数据的小时, 包括从000到2300, 分别代表的是[000,100)到[2300,2400), 即每日的第1小时和最后1小时
	TotalOnlineTime int64 `json:"total_online_time"` // 返回的 json 有这个字段, 文档中没有, 都是 0 值, 可能没有实际意义!!!
	UserReadData
}

// 获取图文统计分时数据.
func (clt *Client) GetUserReadHour(req *Request) (list []UserReadHourData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UserReadHourData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getuserreadhour?access_token="
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

// 图文分享转发数据
type UserShareData struct {
	RefDate    string `json:"ref_date"`    // 数据的日期, YYYY-MM-DD 格式
	UserSource int    `json:"user_source"` // 返回的 json 有这个字段, 文档中没有, 都是 0 值, 可能没有实际意义!!!
	ShareScene int    `json:"share_scene"` // 分享的场景, 1代表好友转发 2代表朋友圈 3代表腾讯微博 255代表其他
	ShareCount int    `json:"share_count"` // 分享的次数
	ShareUser  int    `json:"share_user"`  // 分享的人数
}

// 获取图文分享转发数据.
func (clt *Client) GetUserShare(req *Request) (list []UserShareData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UserShareData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getusershare?access_token="
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

// 图文分享转发分时数据
type UserShareHourData struct {
	RefHour int `json:"ref_hour"` // 数据的小时, 包括从000到2300, 分别代表的是[000,100)到[2300,2400), 即每日的第1小时和最后1小时
	UserShareData
}

// 获取图文分享转发分时数据.
func (clt *Client) GetUserShareHour(req *Request) (list []UserShareHourData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UserShareHourData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getusersharehour?access_token="
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
