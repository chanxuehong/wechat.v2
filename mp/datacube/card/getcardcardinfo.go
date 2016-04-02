// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package card

import (
	"github.com/chanxuehong/wechat/mp"
)

// 免费券数据
type CardData struct {
	RefDate      string `json:"ref_date"`     // 日期信息, YYYY-MM-DD
	CardId       string `json:"card_id"`      // 卡券ID
	CardType     int    `json:"card_type"`    // cardtype:0：折扣券，1：代金券，2：礼品券，3：优惠券，4：团购券（暂不支持拉取特殊票券类型数据，电影票、飞机票、会议门票、景区门票）
	IsPay        int    `json:"is_pay"`       // 是否付费券。0为非付费券，1为付费券
	ViewCount    int    `json:"view_cnt"`     // 浏览次数
	ViewUser     int    `json:"view_user"`    // 浏览人数
	ReceiveCount int    `json:"receive_cnt"`  // 领取次数
	ReceiveUser  int    `json:"receive_user"` // 领取人数
	VerifyCount  int    `json:"verify_cnt"`   // 使用次数
	VerifyUser   int    `json:"verify_user"`  // 使用人数
	GivenCount   int    `json:"given_cnt"`    // 转赠次数
	GivenUser    int    `json:"given_user"`   // 转赠人数
	ExpireCount  int    `json:"expire_cnt"`   // 过期次数
	ExpireUser   int    `json:"expire_user"`  // 过期人数
}

// 获取免费券数据接口
func GetCardInfo(clt *mp.Client, req *Request) (list []CardData, err error) {
	var result struct {
		mp.Error
		List []CardData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getcardcardinfo?access_token="
	if err = clt.PostJSON(incompleteURL, req, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}
