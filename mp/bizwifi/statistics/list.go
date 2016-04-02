// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package statistics

import (
	"github.com/chanxuehong/wechat/mp"
)

type Statistics struct {
	ShopId     int64 `json:"shop_id"`     // 门店ID，-1为总统计
	StatisTime int64 `json:"statis_time"` // 统计时间，单位为毫秒
	TotalUser  int   `json:"total_user"`  // 微信连wifi成功人数
	HomepageUV int   `json:"homepage_uv"` // 商家主页访问人数
	NewFans    int   `json:"new_fans"`    // 新增公众号关注人数
	TotalFans  int   `json:"total_fans"`  // 累计公众号关注人数
}

// 数据统计
//  shopId     按门店ID搜索，-1为总统计
//  beginDate: 起始日期时间，格式yyyy-mm-dd，最长时间跨度为30天
//  endDate:   结束日期时间戳，格式yyyy-mm-dd，最长时间跨度为30天
func List(clt *mp.Client, shopId int64, beginDate, endDate string) (data []Statistics, err error) {
	request := struct {
		ShopId    int64  `json:"shop_id"`
		BeginDate string `json:"begin_date"`
		EndDate   string `json:"end_date"`
	}{
		ShopId:    shopId,
		BeginDate: beginDate,
		EndDate:   endDate,
	}

	var result struct {
		mp.Error
		Data []Statistics `json:"data"`
	}

	incompleteURL := "https://api.weixin.qq.com/bizwifi/statistics/list?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}

	data = result.Data
	return
}
