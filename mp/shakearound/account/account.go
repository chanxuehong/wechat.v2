// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com), magicshui(shuiyuzhe@gmail.com), Harry Rong(harrykobe@gmail.com)

package account

import (
	"github.com/chanxuehong/wechat/mp"
)

type RegisterParameters struct {
	Name                  string   `json:"name"`                    // 必须, 联系人姓名
	PhoneNumber           string   `json:"phone_number"`            // 必须, 联系人电话
	Email                 string   `json:"email"`                   // 必须, 联系人邮箱
	IndustryId            string   `json:"industry_id"`             // 必须, 平台定义的行业代号，具体请查看链接行业代号
	QualificationCertURLs []string `json:"qualification_cert_urls"` // 必须, 相关资质文件的图片url，图片需先上传至微信侧服务器，用“素材管理-上传图片素材”接口上传图片，返回的图片URL再配置在此处；当不需要资质文件时，数组内可以不填写url
	ApplyReason           string   `json:"apply_reason,omitempty"`  // 可选, 申请理由
}

// 申请开通功能
func Register(clt *mp.Client, para *RegisterParameters) (err error) {
	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/shakearound/account/register?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

type AuditStatus struct {
	ApplyTime    int64  `json:"apply_time"`    // 提交申请的时间戳
	AuditStatus  int    `json:"audit_status"`  // 审核状态。0：审核未通过、1：审核中、2：审核已通过；审核会在三个工作日内完成
	AuditComment string `json:"audit_comment"` // 审核备注，包括审核不通过的原因
	AuditTime    int64  `json:"audit_time"`    // 确定审核结果的时间戳；若状态为审核中，则该时间值为0
}

// 查询审核状态
func GetAuditStatus(clt *mp.Client) (status *AuditStatus, err error) {
	var result struct {
		mp.Error
		AuditStatus `json:"data"`
	}

	incompleteURL := "https://api.weixin.qq.com/shakearound/account/auditstatus?access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	status = &result.AuditStatus
	return
}
