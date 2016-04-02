// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package qrcode

import (
	"github.com/chanxuehong/wechat/mp"
)

// 获取物料二维码
//  shopId: 门店ID
//  imgId:  物料样式编号：
//          0-二维码，可用于自由设计宣传材料；
//          1-桌贴（二维码），100mm×100mm(宽×高)，可直接张贴
func Get(clt *mp.Client, shopId int64, imgId int) (qrcodeURL string, err error) {
	request := struct {
		ShopId int64 `json:"shop_id"`
		ImgId  int   `json:"img_id"`
	}{
		ShopId: shopId,
		ImgId:  imgId,
	}

	var result struct {
		mp.Error
		Data struct {
			QRCodeURL string `json:"qrcode_url"`
		} `json:"data"`
	}

	incompleteURL := "https://api.weixin.qq.com/bizwifi/qrcode/get?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}

	qrcodeURL = result.Data.QRCodeURL
	return
}
