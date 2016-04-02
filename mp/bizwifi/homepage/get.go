// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package homepage

import (
	"github.com/chanxuehong/wechat/mp"
)

type Homepage struct {
	ShopId     int64  `json:"shop_id"`     // 门店ID
	TemplateId int64  `json:"template_id"` // 模板类型
	URL        string `json:"url"`         // 商家主页链接
}

func Get(clt *mp.Client, shopId int64) (homepage *Homepage, err error) {
	request := struct {
		ShopId int64 `json:"shop_id"`
	}{
		ShopId: shopId,
	}

	var result struct {
		mp.Error
		Homepage `json:"data"`
	}

	incompleteURL := "https://api.weixin.qq.com/bizwifi/homepage/get?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}

	homepage = &result.Homepage
	return
}
