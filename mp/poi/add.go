// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package poi

import (
	"errors"

	"github.com/chanxuehong/wechat/mp"
)

type Photo struct {
	PhotoURL string `json:"photo_url"`
}

type PoiAddParameters struct {
	BaseInfo struct {
		Sid          string   `json:"sid,omitempty"`           // 可选, 商户自己的id, 用于后续审核通过收到poi_id 的通知时, 做对应关系. 请商户自己保证唯一识别性
		BusinessName string   `json:"business_name,omitempty"` // 必须, 门店名称(仅为商户名, 如: 国美, 麦当劳, 不应包含地区, 店号等信息, 错误示例: 北京国美)
		BranchName   string   `json:"branch_name,omitempty"`   // 必须, 分店名称(不应包含地区信息, 不应与门店名重复, 错误示例: 北京王府井店)
		Province     string   `json:"province,omitempty"`      // 必须, 门店所在的省份(直辖市填城市名,如: 北京市)
		City         string   `json:"city,omitempty"`          // 必须, 门店所在的城市
		District     string   `json:"district,omitempty"`      // 必须, 门店所在地区
		Address      string   `json:"address,omitempty"`       // 必须, 门店所在的详细街道地址(不要填写省市信息)
		Telephone    string   `json:"telephone,omitempty"`     // 必须, 门店的电话(纯数字, 区号, 分机号均由"-"隔开)
		Categories   []string `json:"categories,omitempty"`    // 必须, 门店的类型(详细分类参见分类附表, 不同级分类用","隔开, 如: 美食, 川菜, 火锅)
		OffsetType   int      `json:"offset_type"`             // 必须, 坐标类型, 1 为火星坐标(目前只能选1)
		Longitude    float64  `json:"longitude"`               // 必须, 门店所在地理位置的经度
		Latitude     float64  `json:"latitude"`                // 必须, 门店所在地理位置的纬度(经纬度均为火星坐标, 最好选用腾讯地图标记的坐标)
		PhotoList    []Photo  `json:"photo_list,omitempty"`    // 必须, 图片列表, url 形式, 可以有多张图片, 尺寸为640*340px. 必须为上一接口生成的url
		Recommend    string   `json:"recommend,omitempty"`     // 可选, 推荐品, 餐厅可为推荐菜; 酒店为推荐; 景点为推荐游玩景点等, 针对自己行业的推荐内容
		Special      string   `json:"special,omitempty"`       // 必须, 特色服务, 如免费wifi, 免费停车, 送货上门等商户能提供的特色功能或服务
		Introduction string   `json:"introduction,omitempty"`  // 可选, 商户简介, 主要介绍商户信息等
		OpenTime     string   `json:"open_time,omitempty"`     // 必须, 营业时间, 24 小时制表示, 用"-"连接, 如8:00-20:00
		AvgPrice     int      `json:"avg_price,omitempty"`     // 可选, 人均价格, 大于0 的整数
	} `json:"base_info"`
}

// 创建门店.
func (clt *Client) PoiAdd(para *PoiAddParameters) (err error) {
	if para == nil {
		return errors.New("nil PoiAddParameters")
	}

	var request = struct {
		*PoiAddParameters `json:"business,omitempty"`
	}{
		PoiAddParameters: para,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/poi/addpoi?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
