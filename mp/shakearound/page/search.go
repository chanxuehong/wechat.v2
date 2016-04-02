// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com), magicshui(shuiyuzhe@gmail.com), Harry Rong(harrykobe@gmail.com)

package page

import (
	"errors"

	"github.com/chanxuehong/wechat/internal/util"
	"github.com/chanxuehong/wechat/mp"
)

type SearchQuery struct {
	Type    int     `json:"type"`               // 查询类型。1： 查询页面id列表中的页面信息；2：分页查询所有页面信息
	PageIds []int64 `json:"page_ids,omitempty"` // 指定页面的id列表；当type为1时，此项为必填
	Begin   *int    `json:"begin,omitempty"`    // 页面列表的起始索引值；当type为2时，此项为必填
	Count   *int    `json:"count,omitempty"`    // 待查询的页面数量，不能超过50个；当type为2时，此项为必填
}

func NewSearchQuery1(pageIds []int64) *SearchQuery {
	return &SearchQuery{
		Type:    1,
		PageIds: pageIds,
	}
}

func NewSearchQuery2(begin, count int) *SearchQuery {
	return &SearchQuery{
		Type:  2,
		Begin: util.Int(begin),
		Count: util.Int(count),
	}
}

type SearchResult struct {
	TotalCount int    `json:"total_count"` // 商户名下的页面总数
	ItemCount  int    `json:"item_count"`  // 查询的页面数量
	Pages      []Page `json:"pages"`       // 查询的页面信息列表
}

type Page struct {
	PageId      int64  `json:"page_id"`     // 摇周边页面唯一ID
	Title       string `json:"title"`       // 在摇一摇页面展示的主标题
	Description string `json:"description"` // 在摇一摇页面展示的副标题
	PageURL     string `json:"page_url"`    // 跳转链接
	IconURL     string `json:"icon_url"`    // 在摇一摇页面展示的图片
	Comment     string `json:"comment"`     // 页面的备注信息
}

// 查询页面列表.
func Search(clt *mp.Client, query *SearchQuery) (rslt *SearchResult, err error) {
	var result struct {
		mp.Error
		SearchResult `json:"data"`
	}

	incompleteURL := "https://api.weixin.qq.com/shakearound/page/search?access_token="
	if err = clt.PostJSON(incompleteURL, query, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}

	result.SearchResult.ItemCount = len(result.SearchResult.Pages)
	rslt = &result.SearchResult
	return
}

// PageIterator
//
//  iter, err := NewPageIterator(*mp.Client, *SearchQuery)
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
type PageIterator struct {
	clt *mp.Client

	nextQuery *SearchQuery // 下一次查询参数

	lastSearchResult  *SearchResult // 最近一次获取的数据
	nextPageHasCalled bool          // NextPage() 是否调用过
}

func (iter *PageIterator) TotalCount() int {
	return iter.lastSearchResult.TotalCount
}

func (iter *PageIterator) HasNext() bool {
	if !iter.nextPageHasCalled { // 第一次调用需要特殊对待
		return iter.lastSearchResult.ItemCount > 0 ||
			*iter.nextQuery.Begin < iter.lastSearchResult.TotalCount
	}

	return *iter.nextQuery.Begin < iter.lastSearchResult.TotalCount
}

func (iter *PageIterator) NextPage() (pages []Page, err error) {
	if !iter.nextPageHasCalled { // 第一次调用需要特殊对待
		iter.nextPageHasCalled = true

		pages = iter.lastSearchResult.Pages
		return
	}

	rslt, err := Search(iter.clt, iter.nextQuery)
	if err != nil {
		return
	}

	*iter.nextQuery.Begin += rslt.ItemCount
	iter.lastSearchResult = rslt

	pages = rslt.Pages
	return
}

func NewPageIterator(clt *mp.Client, query *SearchQuery) (iter *PageIterator, err error) {
	if query.Type != 2 {
		err = errors.New("Unsupported SearchQuery.Type")
		return
	}

	// 逻辑上相当于第一次调用 PageIterator.NextPage, 因为第一次调用 PageIterator.HasNext 需要数据支撑, 所以提前获取了数据

	rslt, err := Search(clt, query)
	if err != nil {
		return
	}

	*query.Begin += rslt.ItemCount

	iter = &PageIterator{
		clt: clt,

		nextQuery: query,

		lastSearchResult:  rslt,
		nextPageHasCalled: false,
	}
	return
}
