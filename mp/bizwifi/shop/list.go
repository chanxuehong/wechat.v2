// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package shop

import (
	"errors"

	"github.com/chanxuehong/wechat/mp"
)

type Shop struct {
	Id           int64  `json:"shop_id"`       // 门店ID
	Name         string `json:"shop_name"`     // 门店名称
	SSID         string `json:"ssid"`          // 无线网络设备的ssid，未添加设备为空
	ProtocolType int    `json:"protocol_type"` // 门店内设备的设备类型，0-未添加设备，1-专业型设备，4-通用型设备
}

type ListResult struct {
	PageIndex int `json:"pageindex"` // 分页下标
	PageCount int `json:"pagecount"` // 分页页数

	TotalCount int `json:"totalcount"` // 总数
	ItemCount  int `json:"itemcount"`  // 当前页列表大小

	Records []Shop `json:"records"` // 当前页列表数组
}

// 获取WiFi门店列表.
//  pageIndex: 分页下标，默认从1开始
//  pageSize:  每页的个数，默认10个，最大20个
func List(clt *mp.Client, pageIndex, pageSize int) (rslt *ListResult, err error) {
	if pageIndex < 1 {
		err = errors.New("Incorrect pageIndex")
		return
	}
	if pageSize < 1 {
		err = errors.New("Incorrect pageSize")
		return
	}

	request := struct {
		PageIndex int `json:"pageindex"`
		PageSize  int `json:"pagesize"`
	}{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}

	var result struct {
		mp.Error
		ListResult `json:"data"`
	}

	incompleteURL := "https://api.weixin.qq.com/bizwifi/shop/list?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}

	result.ListResult.ItemCount = len(result.ListResult.Records)
	rslt = &result.ListResult
	return
}

// ShopIterator
//
//  iter, err := NewShopIterator(clt, pageIndex, pageSize)
//  if err != nil {
//      // TODO: 增加你的代码
//  }
//
//  for iter.HasNext() {
//      items, err := iter.NextPage()
//      if err != nil {
//          // TODO: 增加你的代码
//      }
//      // TODO: 增加你的代码
//  }
type ShopIterator struct {
	clt *mp.Client

	pageSize      int
	nextPageIndex int

	lastListResult    *ListResult // 最近一次获取的数据
	nextPageHasCalled bool        // NextPage() 是否调用过
}

func (iter *ShopIterator) TotalCount() int {
	return iter.lastListResult.TotalCount
}

func (iter *ShopIterator) HasNext() bool {
	if !iter.nextPageHasCalled { // 第一次调用需要特殊对待
		return iter.lastListResult.ItemCount > 0 ||
			iter.nextPageIndex <= iter.lastListResult.PageCount
	}

	return iter.nextPageIndex <= iter.lastListResult.PageCount
}

func (iter *ShopIterator) NextPage() (records []Shop, err error) {
	if !iter.nextPageHasCalled { // 第一次调用需要特殊对待
		iter.nextPageHasCalled = true

		records = iter.lastListResult.Records
		return
	}

	rslt, err := List(iter.clt, iter.nextPageIndex, iter.pageSize)
	if err != nil {
		return
	}

	iter.nextPageIndex++
	iter.lastListResult = rslt

	records = rslt.Records
	return
}

func NewShopIterator(clt *mp.Client, pageIndex, pageSize int) (iter *ShopIterator, err error) {
	// 逻辑上相当于第一次调用 ShopIterator.NextPage, 因为第一次调用 ShopIterator.HasNext 需要数据支撑, 所以提前获取了数据

	rslt, err := List(clt, pageIndex, pageSize)
	if err != nil {
		return
	}

	iter = &ShopIterator{
		clt: clt,

		pageSize:      pageSize,
		nextPageIndex: pageIndex + 1,

		lastListResult:    rslt,
		nextPageHasCalled: false,
	}
	return
}
