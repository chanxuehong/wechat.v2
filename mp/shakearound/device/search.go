// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com), magicshui(shuiyuzhe@gmail.com), Harry Rong(harrykobe@gmail.com)

package device

import (
	"errors"

	"github.com/chanxuehong/wechat/internal/util"
	"github.com/chanxuehong/wechat/mp"
)

type SearchQuery struct {
	Type              int                 `json:"type"`                         // 查询类型。1：查询设备id列表中的设备；2：分页查询所有设备信息；3：分页查询某次申请的所有设备信息
	DeviceIdentifiers []*DeviceIdentifier `json:"device_identifiers,omitempty"` // 指定的设备 ； 当type为1时，此项为必填
	ApplyId           *int64              `json:"apply_id,omitempty"`           // 批次ID，申请设备ID时所返回的批次ID；当type为3时，此项为必填
	Begin             *int                `json:"begin,omitempty"`              // 设备列表的起始索引值
	Count             *int                `json:"count,omitempty"`              // 待查询的设备数量，不能超过50个
}

func NewSearchQuery1(deviceIdentifiers []*DeviceIdentifier) *SearchQuery {
	return &SearchQuery{
		Type:              1,
		DeviceIdentifiers: deviceIdentifiers,
	}
}

func NewSearchQuery2(begin, count int) *SearchQuery {
	return &SearchQuery{
		Type:  2,
		Begin: util.Int(begin),
		Count: util.Int(count),
	}
}

func NewSearchQuery3(applyId int64, begin, count int) *SearchQuery {
	return &SearchQuery{
		Type:    3,
		ApplyId: util.Int64(applyId),
		Begin:   util.Int(begin),
		Count:   util.Int(count),
	}
}

type SearchResult struct {
	TotalCount int      `json:"total_count"` // 商户名下的设备总量
	ItemCount  int      `json:"item_count"`  // 查询的设备数量
	Devices    []Device `json:"devices"`     // 查询的设备信息列表
}

type DeviceBase struct {
	DeviceId int64  `json:"device_id"`
	UUID     string `json:"uuid"`
	Major    int    `json:"major"`
	Minor    int    `json:"minor"`
}

type Device struct {
	DeviceBase
	Comment string `json:"comment"`  // 设备的备注信息
	PageIds string `json:"page_ids"` // 与此设备关联的页面ID列表，用逗号隔开
	Status  int    `json:"status"`   // 激活状态，0：未激活，1：已激活（但不活跃），2：活跃
	PoiId   int64  `json:"poi_id"`   // 设备关联的门店ID，关联门店后，在门店1KM的范围内有优先摇出信息的机会。门店相关信息具体可查看门店相关的接口文档
}

// 查询设备列表.
func Search(clt *mp.Client, query *SearchQuery) (rslt *SearchResult, err error) {
	var result struct {
		mp.Error
		SearchResult `json:"data"`
	}

	incompleteURL := "https://api.weixin.qq.com/shakearound/device/search?access_token="
	if err = clt.PostJSON(incompleteURL, query, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}

	result.SearchResult.ItemCount = len(result.SearchResult.Devices)
	rslt = &result.SearchResult
	return
}

// DeviceIterator
//
//  iter, err := NewDeviceIterator(*mp.Client, *SearchQuery)
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
type DeviceIterator struct {
	clt *mp.Client

	nextQuery *SearchQuery // 下一次查询参数

	lastSearchResult  *SearchResult // 最近一次获取的数据
	nextPageHasCalled bool          // NextPage() 是否调用过
}

func (iter *DeviceIterator) TotalCount() int {
	return iter.lastSearchResult.TotalCount
}

func (iter *DeviceIterator) HasNext() bool {
	if !iter.nextPageHasCalled { // 第一次调用需要特殊对待
		return iter.lastSearchResult.ItemCount > 0 ||
			*iter.nextQuery.Begin < iter.lastSearchResult.TotalCount
	}

	return *iter.nextQuery.Begin < iter.lastSearchResult.TotalCount
}

func (iter *DeviceIterator) NextPage() (devices []Device, err error) {
	if !iter.nextPageHasCalled { // 第一次调用需要特殊对待
		iter.nextPageHasCalled = true

		devices = iter.lastSearchResult.Devices
		return
	}

	rslt, err := Search(iter.clt, iter.nextQuery)
	if err != nil {
		return
	}

	*iter.nextQuery.Begin += rslt.ItemCount
	iter.lastSearchResult = rslt

	devices = rslt.Devices
	return
}

func NewDeviceIterator(clt *mp.Client, query *SearchQuery) (iter *DeviceIterator, err error) {
	if query.Type != 2 {
		err = errors.New("Unsupported SearchQuery.Type")
		return
	}

	// 逻辑上相当于第一次调用 DeviceIterator.NextPage, 因为第一次调用 DeviceIterator.HasNext 需要数据支撑, 所以提前获取了数据

	rslt, err := Search(clt, query)
	if err != nil {
		return
	}

	*query.Begin += rslt.ItemCount

	iter = &DeviceIterator{
		clt: clt,

		nextQuery: query,

		lastSearchResult:  rslt,
		nextPageHasCalled: false,
	}
	return
}
