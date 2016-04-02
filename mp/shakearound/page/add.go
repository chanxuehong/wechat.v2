// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com), magicshui(shuiyuzhe@gmail.com), Harry Rong(harrykobe@gmail.com)

package page

import (
	"github.com/chanxuehong/wechat/mp"
)

type AddParameters struct {
	Title       string `json:"title"`             // 必须, 在摇一摇页面展示的主标题，不超过6个字
	Description string `json:"description"`       // 必须, 在摇一摇页面展示的副标题，不超过7个字
	PageURL     string `json:"page_url"`          // 必须, 跳转链接
	IconURL     string `json:"icon_url"`          // 必须, 在摇一摇页面展示的图片。图片需先上传至微信侧服务器，用“素材管理-上传图片素材”接口上传图片，返回的图片URL再配置在此处
	Comment     string `json:"comment,omitempty"` // 可选, 页面的备注信息，不超过15个字
}

// 新增页面
func Add(clt *mp.Client, para *AddParameters) (pageId int64, err error) {
	var result struct {
		mp.Error
		Data struct {
			PageId int64 `json:"page_id"`
		} `json:"data"`
	}

	incompleteURL := "https://api.weixin.qq.com/shakearound/page/add?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	pageId = result.Data.PageId
	return
}
