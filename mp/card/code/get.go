// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package code

import (
	"github.com/chanxuehong/wechat/mp"
)

// 查询code.
func Get(clt *mp.Client, id *CardItemIdentifier) (info *CardItem, err error) {
	var result struct {
		mp.Error
		CardItem
	}

	incompleteURL := "https://api.weixin.qq.com/card/code/get?access_token="
	if err = clt.PostJSON(incompleteURL, id, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	result.CardItem.Code = id.Code
	info = &result.CardItem
	return
}
