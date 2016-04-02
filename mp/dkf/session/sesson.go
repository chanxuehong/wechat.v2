// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package session

import (
	"net/url"

	"github.com/chanxuehong/wechat/mp"
)

// 创建会话
//  openId:    必须, 客户openid
//  kfAccount: 必须, 完整客服账号，格式为：账号前缀@公众号微信号
//  text:      可选, 附加信息，文本会展示在客服人员的多客服客户端
func CreateSession(clt *mp.Client, openId, kfAccount, text string) (err error) {
	request := struct {
		KfAccount string `json:"kf_account"`
		OpenId    string `json:"openid"`
		Text      string `json:"text,omitempty"`
	}{
		KfAccount: kfAccount,
		OpenId:    openId,
		Text:      text,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/customservice/kfsession/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 关闭会话
//  openId:    必须, 客户openid
//  kfAccount: 必须, 完整客服账号，格式为：账号前缀@公众号微信号
//  text:      可选, 附加信息，文本会展示在客服人员的多客服客户端
func CloseSession(clt *mp.Client, openId, kfAccount, text string) (err error) {
	request := struct {
		KfAccount string `json:"kf_account"`
		OpenId    string `json:"openid"`
		Text      string `json:"text,omitempty"`
	}{
		KfAccount: kfAccount,
		OpenId:    openId,
		Text:      text,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/customservice/kfsession/close?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

type Session struct {
	OpenId     string `json:"openid"`
	KfAccount  string `json:"kf_account"`
	CreateTime int64  `json:"createtime"`
}

// 获取客户的会话
func GetSession(clt *mp.Client, openId string) (ss *Session, err error) {
	var result struct {
		mp.Error
		Session
	}

	incompleteURL := "https://api.weixin.qq.com/customservice/kfsession/getsession?openid=" +
		url.QueryEscape(openId) + "&access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	result.Session.OpenId = openId
	ss = &result.Session
	return
}

// 获取客服的会话列表
//  开发者可以通过本接口获取某个客服正在接待的会话列表。
func GetSessionList(clt *mp.Client, kfAccount string) (list []Session, err error) {
	var result struct {
		mp.Error
		SessionList []Session `json:"sessionlist"`
	}

	// TODO
	//	incompleteURL := "https://api.weixin.qq.com/customservice/kfsession/getsessionlist?kf_account=" +
	//		url.QueryEscape(kfAccount) + "&access_token="
	incompleteURL := "https://api.weixin.qq.com/customservice/kfsession/getsessionlist?kf_account=" +
		kfAccount + "&access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	for i, l := 0, result.SessionList; i < len(l); i++ {
		l[i].KfAccount = kfAccount
	}
	list = result.SessionList
	return
}

// 获取未接入会话列表
//  开发者可以通过本接口获取当前正在等待队列中的会话列表，此接口最多返回最早进入队列的100个未接入会话。
func GetWaitSessionList(clt *mp.Client) (list []Session, totalCount int, err error) {
	var result struct {
		mp.Error
		TotalCount  int       `json:"count"`
		SessionList []Session `json:"waitcaselist"`
	}

	incompleteURL := "https://api.weixin.qq.com/customservice/kfsession/getwaitcase?access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.SessionList
	totalCount = result.TotalCount
	return
}
