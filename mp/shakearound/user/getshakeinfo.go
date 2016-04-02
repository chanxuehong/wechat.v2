// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com), magicshui(shuiyuzhe@gmail.com), Harry Rong(harrykobe@gmail.com)

package user

import (
	"github.com/chanxuehong/wechat/mp"
)

type BeaconInfo struct {
	Distance float64 `json:"distance"` // Beacon信号与手机的距离，单位为米
	UUID     string  `json:"uuid"`
	Major    int     `json:"major"`
	Minor    int     `json:"minor"`
}

type Shakeinfo struct {
	PageId     int64      `json:"page_id"`     // 摇周边页面唯一ID
	BeaconInfo BeaconInfo `json:"beacon_info"` // 设备信息，包括UUID、major、minor，以及距离
	Openid     string     `json:"openid"`      // 商户AppID下用户的唯一标识
	PoiId      *int64     `json:"poi_id"`      // 门店ID，有的话则返回，反之不会在JSON格式内
}

// 获取摇周边的设备及用户信息
//  ticket:  摇周边业务的ticket，可在摇到的URL中得到，ticket生效时间为30分钟，每一次摇都会重新生成新的ticket
//  needPoi: 是否需要返回门店poi_id
func GetShakeInfo(clt *mp.Client, ticket string, needPoi bool) (info *Shakeinfo, err error) {
	request := struct {
		Ticket  string `json:"ticket"`
		NeedPoi int    `json:"need_poi,omitempty"`
	}{
		Ticket: ticket,
	}

	if needPoi {
		request.NeedPoi = 1
	}

	var result struct {
		mp.Error
		Shakeinfo `json:"data"`
	}

	incompleteURL := "https://api.weixin.qq.com/shakearound/user/getshakeinfo?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.Shakeinfo
	return
}
