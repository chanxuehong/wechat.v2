// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package suite

import (
	"errors"

	"github.com/chanxuehong/wechat/corp"
)

type SetAgentParameters struct {
	AgentId            int64  `json:"agentid"`
	ReportLocationFlag *int   `json:"report_location_flag,omitempty"`
	LogoMediaId        string `json:"logo_mediaid,omitempty"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	RedirectDomain     string `json:"redirect_domain,omitempty"`
	IsReportUser       *int   `json:"isreportuser,omitempty"`
	IsReportEnter      *int   `json:"isreportenter,omitempty"`
}

func (clt *Client) SetAgent(authCorpId, permanentCode string, para *SetAgentParameters) (err error) {
	if para == nil {
		return errors.New("nil AgentParameters")
	}

	request := struct {
		SuiteId       string              `json:"suite_id"`
		AuthCorpId    string              `json:"auth_corpid"`
		PermanentCode string              `json:"permanent_code"`
		Agent         *SetAgentParameters `json:"agent,omitempty"`
	}{
		SuiteId:       clt.SuiteId,
		AuthCorpId:    authCorpId,
		PermanentCode: permanentCode,
		Agent:         para,
	}

	var result corp.Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/service/set_agent?suite_access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result
		return
	}
	return
}
