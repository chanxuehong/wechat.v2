// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package datacube

import (
	"errors"

	"github.com/chanxuehong/wechat/mp"
)

// 用户增减数据
type UserSummaryData struct {
	RefDate string `json:"ref_date"` // 数据的日期, YYYY-MM-DD 格式

	// 用户的渠道, 数值代表的含义如下:
	// 0  代表其他
	// 1  xxx(文档没有说明)
	// 2  xxx(文档没有说明)
	// 3  代表扫二维码
	// 4  xxx(文档没有说明)
	// 5  xxx(文档没有说明)
	// 17 代表名片分享
	// 35 代表搜号码(即微信添加朋友页的搜索)
	// 39 代表查询微信公众帐号
	// 43 代表图文页右上角菜单
	UserSource int `json:"user_source"`

	NewUser    int `json:"new_user"`    // 新增的用户数量
	CancelUser int `json:"cancel_user"` // 取消关注的用户数量, new_user 减去 cancel_user即为净增用户数量
}

// 获取用户增减数据.
func (clt *Client) GetUserSummary(req *Request) (list []UserSummaryData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UserSummaryData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getusersummary?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, req, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 累计用户数据
type UserCumulateData struct {
	RefDate      string `json:"ref_date"`      // 数据的日期, YYYY-MM-DD 格式
	UserSource   int    `json:"user_source"`   // 返回的 json 有这个字段, 文档中没有, 都是 0 值, 可能没有实际意义!!!
	CumulateUser int    `json:"cumulate_user"` // 总用户量
}

// 获取累计用户数据.
func (clt *Client) GetUserCumulate(req *Request) (list []UserCumulateData, err error) {
	if req == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		mp.Error
		List []UserCumulateData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getusercumulate?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, req, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}
