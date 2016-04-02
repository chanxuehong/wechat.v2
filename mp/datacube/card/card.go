// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

// 卡券数据统计接口
package card

// 请求数据结构
type Request struct {
	BeginDate  string `json:"begin_date"`        // 查询数据的起始时间, YYYY-MM-DD 格式;
	EndDate    string `json:"end_date"`          // 查询数据的截至时间, YYYY-MM-DD 格式;
	CondSource int    `json:"cond_source"`       // 卡券来源，0为公众平台创建的卡券数据、1是API创建的卡券数据
	CardId     string `json:"card_id,omitempty"` // 可选; 卡券ID。填写后，指定拉出该卡券的相关数据。
}
