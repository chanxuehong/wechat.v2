// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package suite

import (
	"github.com/chanxuehong/wechat/corp"
)

type AgentInfo struct {
	AgentId        int64  `json:"agentid"`
	Name           string `json:"name"`
	SquareLogoURL  string `json:"square_logo_url"`
	RoundLogoURL   string `json:"round_logo_url"`
	Description    string `json:"description"`
	AllowUserInfos struct {
		UserList []AgentInfoUser `json:"user,omitempty"`
	} `json:"allow_userinfos"`
	AllowParties struct {
		PartyIdList []int64 `json:"partyid,omitempty"`
	} `json:"allow_partys"`
	AllowTags struct {
		TagIdList []int64 `json:"tagid,omitempty"`
	} `json:"allow_tags"`
	Closed             int    `json:"close"`
	RedirectDomain     string `json:"redirect_domain"`
	ReportLocationFlag int    `json:"report_location_flag"`
	IsReportUser       int    `json:"isreportuser"`
	IsReportEnter      int    `json:"isreportenter"`
}

type AgentInfoUser struct {
	UserId string `json:"userid"`
	Status int    `json:"status"`
}

// 获取企业号应用
func (clt *Client) GetAgent(authCorpId, permanentCode string, agentId int64) (info *AgentInfo, err error) {
	request := struct {
		SuiteId       string `json:"suite_id"`
		AuthCorpId    string `json:"auth_corpid"`
		PermanentCode string `json:"permanent_code"`
		AgentId       int64  `json:"agentid"`
	}{
		SuiteId:       clt.SuiteId,
		AuthCorpId:    authCorpId,
		PermanentCode: permanentCode,
		AgentId:       agentId,
	}

	var result struct {
		corp.Error
		AgentInfo
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/service/get_agent?suite_access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.AgentInfo
	return
}
