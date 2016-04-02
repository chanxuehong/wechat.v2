// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package material

import (
	"fmt"

	"github.com/chanxuehong/wechat/mp"
)

// 删除永久素材.
func (clt *Client) DeleteMaterial(mediaId string) (err error) {
	var request = struct {
		MediaId string `json:"media_id"`
	}{
		MediaId: mediaId,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/material/del_material?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 公众号永久素材总数信息
type MaterialCountInfo struct {
	VoiceCount int `json:"voice_count"`
	VideoCount int `json:"video_count"`
	ImageCount int `json:"image_count"`
	NewsCount  int `json:"news_count"`
}

// 获取素材总数.
func (clt *Client) GetMaterialCount() (info *MaterialCountInfo, err error) {
	var result struct {
		mp.Error
		MaterialCountInfo
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/material/get_materialcount?access_token="
	if err = ((*mp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.MaterialCountInfo
	return
}

type MaterialInfo struct {
	MediaId    string `json:"media_id"`    // 素材id
	Name       string `json:"name"`        // 文件名称
	UpdateTime int64  `json:"update_time"` // 最后更新时间
	URL        string `json:"url"`         // 当获取的列表是图片素材列表时, 该字段是图片的URL
}

type BatchGetMaterialResult struct {
	TotalCount int            `json:"total_count"` // 该类型的素材的总数
	ItemCount  int            `json:"item_count"`  // 本次调用获取的素材的数量
	Items      []MaterialInfo `json:"item"`        // 本次调用获取的素材列表
}

// 获取素材列表.
//
//  MaterialType: 素材的类型, 图片(image), 视频(video), 语音 (voice)
//  offset:       从全部素材的该偏移位置开始返回, 0表示从第一个素材
//  count:        返回素材的数量, 取值在1到20之间
func (clt *Client) BatchGetMaterial(MaterialType string, offset, count int) (rslt *BatchGetMaterialResult, err error) {
	switch MaterialType {
	case MaterialTypeImage, MaterialTypeVideo, MaterialTypeVoice:
	default:
		err = fmt.Errorf("Incorrect MaterialType: %s", MaterialType)
		return
	}

	if offset < 0 {
		err = fmt.Errorf("Incorrect offset: %d", offset)
		return
	}
	if count <= 0 {
		err = fmt.Errorf("Incorrect count: %d", count)
		return
	}

	var request = struct {
		MaterialType string `json:"type"`
		Offset       int    `json:"offset"`
		Count        int    `json:"count"`
	}{
		MaterialType: MaterialType,
		Offset:       offset,
		Count:        count,
	}

	var result struct {
		mp.Error
		BatchGetMaterialResult
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}

	rslt = &result.BatchGetMaterialResult
	return
}

// MaterialIterator
//
//  iter, err := Client.MaterialIterator(MaterialTypeImage, 0, 10)
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
type MaterialIterator struct {
	clt *Client // 关联的微信 Client

	materialType string // image, video, voice
	nextOffset   int    // 下一次获取数据时的 offset
	count        int    // 步长

	lastBatchGetMaterialResult *BatchGetMaterialResult // 最近一次获取的数据
	nextPageHasCalled          bool                    // NextPage() 是否调用过
}

func (iter *MaterialIterator) TotalCount() int {
	return iter.lastBatchGetMaterialResult.TotalCount
}

func (iter *MaterialIterator) HasNext() bool {
	if !iter.nextPageHasCalled { // 第一次调用需要特殊对待
		return iter.lastBatchGetMaterialResult.ItemCount > 0 ||
			iter.nextOffset < iter.lastBatchGetMaterialResult.TotalCount
	}

	return iter.nextOffset < iter.lastBatchGetMaterialResult.TotalCount
}

func (iter *MaterialIterator) NextPage() (items []MaterialInfo, err error) {
	if !iter.nextPageHasCalled { // 第一次调用需要特殊对待
		iter.nextPageHasCalled = true

		items = iter.lastBatchGetMaterialResult.Items
		return
	}

	rslt, err := iter.clt.BatchGetMaterial(iter.materialType, iter.nextOffset, iter.count)
	if err != nil {
		return
	}

	iter.nextOffset += rslt.ItemCount
	iter.lastBatchGetMaterialResult = rslt

	items = rslt.Items
	return
}

func (clt *Client) MaterialIterator(MaterialType string, offset, count int) (iter *MaterialIterator, err error) {
	// 逻辑上相当于第一次调用 MaterialIterator.NextPage, 因为第一次调用 MaterialIterator.HasNext 需要数据支撑, 所以提前获取了数据

	rslt, err := clt.BatchGetMaterial(MaterialType, offset, count)
	if err != nil {
		return
	}

	iter = &MaterialIterator{
		clt: clt,

		materialType: MaterialType,
		nextOffset:   offset + rslt.ItemCount,
		count:        count,

		lastBatchGetMaterialResult: rslt,
		nextPageHasCalled:          false,
	}
	return
}
