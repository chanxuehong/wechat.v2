// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package addresslist

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/chanxuehong/wechat/corp"
)

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// 创建成员的参数
type UserCreateParameters struct {
	UserId     string  `json:"userid,omitempty"`     // 必须;  员工UserID. 对应管理端的帐号, 企业内必须唯一. 长度为1~64个字符
	Name       string  `json:"name,omitempty"`       // 必须;  成员名称. 长度为1~64个字符
	Department []int64 `json:"department,omitempty"` // 非必须; 成员所属部门id列表. 注意, 每个部门的直属员工上限为1000个
	Position   string  `json:"position,omitempty"`   // 非必须; 职位信息. 长度为0~64个字符
	Mobile     string  `json:"mobile,omitempty"`     // 非必须; 手机号码. 企业内必须唯一, mobile/weixinid/email三者不能同时为空
	Email      string  `json:"email,omitempty"`      // 非必须; 邮箱. 长度为0~64个字符. 企业内必须唯一
	WeixinId   string  `json:"weixinid,omitempty"`   // 非必须; 微信号. 企业内必须唯一. (注意: 是微信号, 不是微信的名字)
	ExtAttr    struct {
		Attrs []Attribute `json:"attrs,omitempty"`
	} `json:"extattr"` // 非必须; 扩展属性. 扩展属性需要在WEB管理端创建后才生效, 否则忽略未知属性的赋值
}

// 创建成员
func (clt *Client) UserCreate(para *UserCreateParameters) (err error) {
	if para == nil {
		err = errors.New("nil parameters")
		return
	}

	var result corp.Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/user/create?access_token="
	if err = ((*corp.Client)(clt)).PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result
		return
	}
	return
}

type UserUpdateParameters struct {
	UserId     string  `json:"userid,omitempty"`     // 必须;  员工UserID. 对应管理端的帐号, 企业内必须唯一. 长度为1~64个字符
	Name       string  `json:"name,omitempty"`       // 非必须; 成员名称. 长度为0~64个字符
	Department []int64 `json:"department,omitempty"` // 非必须; 成员所属部门id列表. 注意, 每个部门的直属员工上限为1000个
	Position   string  `json:"position,omitempty"`   // 非必须; 职位信息. 长度为0~64个字符
	Mobile     string  `json:"mobile,omitempty"`     // 非必须; 手机号码. 企业内必须唯一, mobile/weixinid/email三者不能同时为空
	Email      string  `json:"email,omitempty"`      // 非必须; 邮箱. 长度为0~64个字符. 企业内必须唯一
	WeixinId   string  `json:"weixinid,omitempty"`   // 非必须; 微信号. 企业内必须唯一. (注意: 是微信号, 不是微信的名字)
	Enable     *int    `json:"enable,omitempty"`     // 非必须; 启用/禁用成员. 1表示启用成员, 0表示禁用成员
	ExtAttr    struct {
		Attrs []Attribute `json:"attrs,omitempty"`
	} `json:"extattr"` // 非必须; 扩展属性. 扩展属性需要在WEB管理端创建后才生效, 否则忽略未知属性的赋值
}

func (para *UserUpdateParameters) SetEnable(b bool) {
	var x int
	if b {
		x = 1
	}
	para.Enable = &x
}

// 更新成员
func (clt *Client) UserUpdate(para *UserUpdateParameters) (err error) {
	if para == nil {
		err = errors.New("nil parameters")
		return
	}

	var result corp.Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/user/update?access_token="
	if err = ((*corp.Client)(clt)).PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 删除成员
func (clt *Client) UserDelete(userId string) (err error) {
	var result corp.Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/user/delete?userid=" +
		url.QueryEscape(userId) + "&access_token="
	if err = ((*corp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 批量删除成员
func (clt *Client) UserBatchDelete(UserIdList []string) (err error) {
	if len(UserIdList) <= 0 {
		return
	}

	var request = struct {
		UserIdList []string `json:"useridlist,omitempty"`
	}{
		UserIdList: UserIdList,
	}

	var result corp.Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/user/batchdelete?access_token="
	if err = ((*corp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result
		return
	}
	return
}

type UserInfo struct {
	Id         string  `json:"userid"`               // 员工UserID. 对应管理端的帐号
	Name       string  `json:"name"`                 // 成员名称
	Department []int64 `json:"department,omitempty"` // 成员所属部门id列表
	Position   string  `json:"position"`             // 职位信息
	Mobile     string  `json:"mobile"`               // 手机号码
	Email      string  `json:"email"`                // 邮箱
	WeixinId   string  `json:"weixinid"`             // 微信号
	Avatar     string  `json:"avatar"`               // 头像url. 注: 如果要获取小图将url最后的"/0"改成"/64"即可
	Status     int     `json:"status"`               // 关注状态: 1=已关注, 2=已冻结, 4=未关注
	ExtAttr    struct {
		Attrs []Attribute `json:"attrs,omitempty"`
	} `json:"extattr"` // 扩展属性
}

func (clt *Client) UserInfo(userId string) (info *UserInfo, err error) {
	var result struct {
		corp.Error
		UserInfo
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/user/get?userid=" +
		url.QueryEscape(userId) + "&access_token="
	if err = ((*corp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.UserInfo
	return
}

type UserBaseInfo struct {
	Id   string `json:"userid"` // 员工UserID
	Name string `json:"name"`   // 成员名称
}

// 获取部门成员(基本)
//  departmentId: 获取的部门id
//  fetchChild:   是否递归获取子部门下面的成员
//  status:       0 获取全部员工, 1 获取已关注成员列表, 2 获取禁用成员列表, 4 获取未关注成员列表.
//                status可叠加(可用逻辑运算符 | 来叠加, 一般都是后面 3 个叠加).
func (clt *Client) UserSimpleList(departmentId int64,
	fetchChild bool, status int) (UserList []UserBaseInfo, err error) {

	var result struct {
		corp.Error
		UserList []UserBaseInfo `json:"userlist"`
	}

	var fetchChildStr string
	if fetchChild {
		fetchChildStr = "1"
	} else {
		fetchChildStr = "0"
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist" +
		"?department_id=" + strconv.FormatInt(departmentId, 10) +
		"&fetch_child=" + fetchChildStr +
		"&status=" + strconv.FormatInt(int64(status), 10) +
		"&access_token="
	if err = ((*corp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	UserList = result.UserList
	return
}

// 获取部门成员(详情)
//  departmentId: 获取的部门id
//  fetchChild:   是否递归获取子部门下面的成员
//  status:       0 获取全部员工, 1 获取已关注成员列表, 2 获取禁用成员列表, 4 获取未关注成员列表.
//                status可叠加(可用逻辑运算符 | 来叠加, 一般都是后面 3 个叠加).
func (clt *Client) UserList(departmentId int64,
	fetchChild bool, status int) (UserList []UserInfo, err error) {

	var result struct {
		corp.Error
		UserList []UserInfo `json:"userlist"`
	}

	var fetchChildStr string
	if fetchChild {
		fetchChildStr = "1"
	} else {
		fetchChildStr = "0"
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/user/list" +
		"?department_id=" + strconv.FormatInt(departmentId, 10) +
		"&fetch_child=" + fetchChildStr +
		"&status=" + strconv.FormatInt(int64(status), 10) +
		"&access_token="
	if err = ((*corp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	UserList = result.UserList
	return
}
