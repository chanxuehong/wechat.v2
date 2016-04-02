// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com), magicshui(shuiyuzhe@gmail.com), Harry Rong(harrykobe@gmail.com)

package relation

import (
	"errors"

	"github.com/chanxuehong/wechat/internal/util"
	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/shakearound/device"
)

type SearchQuery struct {
	Type             int                      `json:"type"`                        // 查询方式。1： 查询设备的关联关系；2：查询页面的关联关系
	DeviceIdentifier *device.DeviceIdentifier `json:"device_identifier,omitempty"` // 指定的设备；当type为1时，此项为必填
	PageId           *int64                   `json:"page_id,omitempty"`           // 指定的页面id；当type为2时，此项为必填
	Begin            *int                     `json:"begin,omitempty"`             // 关联关系列表的起始索引值；当type为2时，此项为必填
	Count            *int                     `json:"count,omitempty"`             // 待查询的关联关系数量，不能超过50个；当type为2时，此项为必填
}

func NewSearchQuery1(deviceIdentifier *device.DeviceIdentifier) *SearchQuery {
	return &SearchQuery{
		Type:             1,
		DeviceIdentifier: deviceIdentifier,
	}
}

func NewSearchQuery1X(deviceIdentifier *device.DeviceIdentifier, begin, count int) *SearchQuery {
	return &SearchQuery{
		Type:             1,
		DeviceIdentifier: deviceIdentifier,
		Begin:            util.Int(begin),
		Count:            util.Int(count),
	}
}

func NewSearchQuery2(pageId int64, begin, count int) *SearchQuery {
	return &SearchQuery{
		Type:   2,
		PageId: util.Int64(pageId),
		Begin:  util.Int(begin),
		Count:  util.Int(count),
	}
}

type SearchResult struct {
	TotalCount int        `json:"total_count"` // 设备或页面的关联关系总数
	ItemCount  int        `json:"item_count"`  // 查询的关联关系数量
	Relations  []Relation `json:"relations"`   // 查询的关联关系列表
}

type Relation struct {
	device.DeviceBase
	PageId int64 `json:"page_id"`
}

// 查询设备与页面的关联关系.
func Search(clt *mp.Client, query *SearchQuery) (rslt *SearchResult, err error) {
	var result struct {
		mp.Error
		SearchResult `json:"data"`
	}

	incompleteURL := "https://api.weixin.qq.com/shakearound/relation/search?access_token="
	if err = clt.PostJSON(incompleteURL, query, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}

	result.SearchResult.ItemCount = len(result.SearchResult.Relations)
	rslt = &result.SearchResult
	return
}

// RelationIterator
//
//  iter, err := NewRelationIterator(*mp.Client, *SearchQuery)
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
type RelationIterator struct {
	clt *mp.Client

	nextQuery *SearchQuery // 下一次查询参数

	lastSearchResult  *SearchResult // 最近一次获取的数据
	nextPageHasCalled bool          // NextPage() 是否调用过
}

func (iter *RelationIterator) TotalCount() int {
	return iter.lastSearchResult.TotalCount
}

func (iter *RelationIterator) HasNext() bool {
	if !iter.nextPageHasCalled { // 第一次调用需要特殊对待
		return iter.lastSearchResult.ItemCount > 0 ||
			*iter.nextQuery.Begin < iter.lastSearchResult.TotalCount
	}

	return *iter.nextQuery.Begin < iter.lastSearchResult.TotalCount
}

func (iter *RelationIterator) NextPage() (relations []Relation, err error) {
	if !iter.nextPageHasCalled { // 第一次调用需要特殊对待
		iter.nextPageHasCalled = true

		relations = iter.lastSearchResult.Relations
		return
	}

	rslt, err := Search(iter.clt, iter.nextQuery)
	if err != nil {
		return
	}

	*iter.nextQuery.Begin += rslt.ItemCount
	iter.lastSearchResult = rslt

	relations = rslt.Relations
	return
}

func NewRelationIterator(clt *mp.Client, query *SearchQuery) (iter *RelationIterator, err error) {
	if query.Begin == nil {
		err = errors.New("nil SearchQuery.Begin")
		return
	}
	if query.Count == nil {
		err = errors.New("nil SearchQuery.Count")
		return
	}

	// 逻辑上相当于第一次调用 RelationIterator.NextPage, 因为第一次调用 RelationIterator.HasNext 需要数据支撑, 所以提前获取了数据

	rslt, err := Search(clt, query)
	if err != nil {
		return
	}

	*query.Begin += rslt.ItemCount

	iter = &RelationIterator{
		clt: clt,

		nextQuery: query,

		lastSearchResult:  rslt,
		nextPageHasCalled: false,
	}
	return
}
