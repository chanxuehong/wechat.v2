// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package poi

import (
	"errors"

	"github.com/chanxuehong/wechat/mp"
)

// 特别注意：
// PoiUpdateParameters 字段，若有填写内容则为覆盖更新，若无内容则视为不修改，维持原有内容。
// photo_list 字段为全列表覆盖，若需要增加图片，需将之前图片同样放入list 中，
// 在其后增加新增图片。如：已有A、B、C 三张图片，又要增加D、E 两张图，则需要调用该接口，
// photo_list 传入A、B、C、D、E 五张图片的链接。
type PoiUpdateParameters struct {
	BaseInfo struct {
		PoiId int64 `json:"poi_id,string"`

		Telephone    string  `json:"telephone,omitempty"`    // 门店的电话(纯数字, 区号, 分机号均由"-"隔开)
		PhotoList    []Photo `json:"photo_list,omitempty"`   // 图片列表, url 形式, 可以有多张图片, 尺寸为640*340px. 必须为上一接口生成的url
		Recommend    string  `json:"recommend,omitempty"`    // 推荐品, 餐厅可为推荐菜; 酒店为推荐; 景点为推荐游玩景点等, 针对自己行业的推荐内容
		Special      string  `json:"special,omitempty"`      // 特色服务, 如免费wifi, 免费停车, 送货上门等商户能提供的特色功能或服务
		Introduction string  `json:"introduction,omitempty"` // 商户简介, 主要介绍商户信息等
		OpenTime     string  `json:"open_time,omitempty"`    // 营业时间, 24 小时制表示, 用"-"连接, 如8:00-20:00
		AvgPrice     int     `json:"avg_price,omitempty"`    // 人均价格, 大于0 的整数
	} `json:"base_info"`
}

// 修改门店服务信息.
//  商户可以通过该接口, 修改门店的服务信息, 包括: 图片列表, 营业时间, 推荐, 特色服务, 简
//  介, 人均价格, 电话7 个字段. 目前基础字段包括(名称, 坐标, 地址等不可修改)
func (clt *Client) PoiUpdate(para *PoiUpdateParameters) (err error) {
	if para == nil {
		return errors.New("nil PoiUpdateParameters")
	}

	var request = struct {
		*PoiUpdateParameters `json:"business,omitempty"`
	}{
		PoiUpdateParameters: para,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/poi/updatepoi?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
