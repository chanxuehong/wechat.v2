// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package template

import (
	"errors"
	"net/http"

	"github.com/chanxuehong/wechat/mp"
)

type Client mp.Client

func NewClient(srv mp.AccessTokenServer, clt *http.Client) *Client {
	return (*Client)(mp.NewClient(srv, clt))
}

// 设置所属行业.
//  目前 industryId 的个数只能为 2.
func (clt *Client) SetIndustry(industryId ...int64) (err error) {
	if len(industryId) < 2 {
		return errors.New("industryId 的个数不能小于 2")
	}

	var request = struct {
		IndustryId1 int64 `json:"industry_id1"`
		IndustryId2 int64 `json:"industry_id2"`
	}{
		IndustryId1: industryId[0],
		IndustryId2: industryId[1],
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/template/api_set_industry?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 从行业模板库选择模板添加到账号后台, 并返回模板id.
//  templateIdShort: 模板库中模板的编号, 有"TM**"和"OPENTMTM**"等形式.
func (clt *Client) AddTemplate(templateIdShort string) (templateId string, err error) {
	var request = struct {
		TemplateIdShort string `json:"template_id_short"`
	}{
		TemplateIdShort: templateIdShort,
	}

	var result struct {
		mp.Error
		TemplateId string `json:"template_id"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/template/api_add_template?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	templateId = result.TemplateId
	return
}

// 发送模板消息
func (clt *Client) Send(msg *TemplateMessage) (msgid int64, err error) {
	if msg == nil {
		err = errors.New("nil TemplateMessage")
		return
	}

	var result struct {
		mp.Error
		MsgId int64 `json:"msgid"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, msg, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	msgid = result.MsgId
	return
}
