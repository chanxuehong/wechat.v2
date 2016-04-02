// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package addresslist

import (
	"errors"
	"strconv"

	"github.com/chanxuehong/wechat/corp"
)

// 创建部门参数
type DepartmentCreateParameters struct {
	DepartmentName string `json:"name,omitempty"`  // 必须, 部门名称. 长度限制为1~64个字符
	ParentId       int64  `json:"parentid"`        // 必须, 父亲部门id. 根部门id为1
	Order          *int   `json:"order,omitempty"` // 可选, 在父部门中的次序. 从1开始, 数字越大排序越靠后
	DepartmentId   *int64 `json:"id,omitempty"`    // 可选, 部门ID. 用指定部门ID新建部门, 不指定此参数时, 则自动生成
}

// 创建部门
func (clt *Client) DepartmentCreate(para *DepartmentCreateParameters) (id int64, err error) {
	if para == nil {
		err = errors.New("nil parameters")
		return
	}

	var result struct {
		corp.Error
		Id int64 `json:"id"`
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/department/create?access_token="
	if err = ((*corp.Client)(clt)).PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	id = result.Id
	return
}

// 更新部门参数
type DepartmentUpdateParameters struct {
	DepartmentId   int64  `json:"id"`                 // 必须, 部门id
	DepartmentName string `json:"name,omitempty"`     // 可选, 更新的部门名称. 长度限制为1~64个字符. 修改部门名称时指定该参数
	ParentId       *int64 `json:"parentid,omitempty"` // 可选, 父亲部门id. 根部门id为1
	Order          *int   `json:"order,omitempty"`    // 可选, 在父部门中的次序. 从1开始, 数字越大排序越靠后, 当数字大于该层部门数时表示移动到最末尾.
}

// 更新部门
func (clt *Client) DepartmentUpdate(para *DepartmentUpdateParameters) (err error) {
	if para == nil {
		err = errors.New("nil parameters")
		return
	}

	var result corp.Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/department/update?access_token="
	if err = ((*corp.Client)(clt)).PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 删除部门
func (clt *Client) DepartmentDelete(id int64) (err error) {
	var result corp.Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/department/delete?id=" +
		strconv.FormatInt(id, 10) + "&access_token="
	if err = ((*corp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result
		return
	}
	return
}

type Department struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	ParentId int64  `json:"parentid"`
}

// 获取 rootId 部门的子部门
func (clt *Client) DepartmentList(rootId int64) (departments []Department, err error) {
	var result struct {
		corp.Error
		Departments []Department `json:"department"`
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/department/list?id=" +
		strconv.FormatInt(rootId, 10) + "&access_token="
	if err = ((*corp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	departments = result.Departments
	return
}
