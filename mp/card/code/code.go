// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package code

// 某一张特定卡券的标识
type CardItemIdentifier struct {
	Code   string `json:"code"`              // 卡券的Code码
	CardId string `json:"card_id,omitempty"` // 卡券ID。创建卡券时use_custom_code填写true时必填。非自定义Code不必填写。
}

// 某一张特定卡券的信息
type CardItem struct {
	Code   string `json:"code"`   // 卡券的Code码
	OpenId string `json:"openid"` // 用户openid
	Card   struct {
		CardId    string `json:"card_id"`    // 卡券ID
		BeginTime int64  `json:"begin_time"` // 起始使用时间
		EndTime   int64  `json:"end_time"`   // 结束时间
	} `json:"card"`
}
