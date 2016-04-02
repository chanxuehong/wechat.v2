// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package record

import (
	"errors"

	"github.com/chanxuehong/wechat/mp"
)

type Record struct {
	Worker    string `json:"worker"`   // 客服账号
	OpenId    string `json:"openid"`   // 用户的标识, 对当前公众号唯一
	OperCode  int    `json:"opercode"` // 操作ID(会话状态)
	Timestamp int64  `json:"time"`     // 操作时间, UNIX时间戳
	Text      string `json:"text"`     // 聊天记录
}

const (
	RecordPageSizeLimit = 50
)

type GetRecordRequest struct {
	StartTime int64  `json:"starttime"`        // 查询开始时间, UNIX时间戳
	EndTime   int64  `json:"endtime"`          // 查询结束时间, UNIX时间戳, 每次查询不能跨日查询
	PageIndex int    `json:"pageindex"`        // 查询第几页, 从1开始
	PageSize  int    `json:"pagesize"`         // 每页大小, 每页最多拉取50条
	OpenId    string `json:"openid,omitempty"` // 普通用户的标识, 对当前公众号唯一
}

// 获取客服聊天记录
func GetRecord(clt *mp.Client, request *GetRecordRequest) (recordList []Record, err error) {
	if request == nil {
		err = errors.New("nil GetRecordRequest")
		return
	}

	var result struct {
		mp.Error
		RecordList []Record `json:"recordlist"`
	}

	incompleteURL := "https://api.weixin.qq.com/customservice/msgrecord/getrecord?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	recordList = result.RecordList
	return
}
