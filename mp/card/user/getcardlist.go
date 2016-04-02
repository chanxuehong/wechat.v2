// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package user

import (
	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/card/code"
)

// 获取用户已领取卡券接口
//  openid: 需要查询的用户openid
//  cardid: 卡券ID。不填写时默认查询当前appid下的卡券。
func GetCardList(clt *mp.Client, openid, cardid string) (list []code.CardItemIdentifier, err error) {
	request := struct {
		OpenId string `json:"openid"`
		CardId string `json:"card_id,omitempty"`
	}{
		OpenId: openid,
		CardId: cardid,
	}

	var result struct {
		mp.Error
		CardList []code.CardItemIdentifier `json:"card_list"`
	}

	incompleteURL := "https://api.weixin.qq.com/card/user/getcardlist?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.CardList
	return
}
