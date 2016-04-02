// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package corp

import "fmt"

const (
	ErrCodeOK                      = 0
	ErrCodeAccessTokenExpired      = 42001 // access_token 过期(无效)返回这个错误
	ErrCodeSuiteAccessTokenExpired = 42009 // suite_access_token 过期(无效)返回这个错误
)

type Error struct {
	// NOTE: StructField 固定这个顺序, RETRY 依赖这个顺序
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("errcode: %d, errmsg: %s", e.ErrCode, e.ErrMsg)
}
