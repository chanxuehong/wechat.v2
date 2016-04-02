// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package code

import (
	"github.com/chanxuehong/wechat/mp"
)

// 设置卡券失效接口.
func Unavailable(clt *mp.Client, id *CardItemIdentifier) (err error) {
	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/card/code/unavailable?access_token="
	if err = clt.PostJSON(incompleteURL, id, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
