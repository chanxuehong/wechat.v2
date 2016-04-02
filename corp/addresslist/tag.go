// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package addresslist

import (
	"strconv"
	"strings"

	"github.com/chanxuehong/wechat/corp"
)

// 创建标签
func (clt *Client) TagCreate(tagName string) (id int64, err error) {
	var request = struct {
		TagName string `json:"tagname"`
	}{
		TagName: tagName,
	}

	var result struct {
		corp.Error
		TagId int64 `json:"tagid"`
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/tag/create?access_token="
	if err = ((*corp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	id = result.TagId
	return
}

// 更新标签名字
func (clt *Client) TagUpdate(id int64, name string) (err error) {
	var request = struct {
		TagId   int64  `json:"tagid"`
		TagName string `json:"tagname"`
	}{
		TagId:   id,
		TagName: name,
	}

	var result corp.Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/tag/update?access_token="
	if err = ((*corp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 删除标签
func (clt *Client) TagDelete(id int64) (err error) {
	var result corp.Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/tag/delete?tagid=" +
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

// 获取标签成员
func (clt *Client) TagInfo(id int64) (userList []UserBaseInfo, departmentList []int64, err error) {
	var result struct {
		corp.Error
		UserList       []UserBaseInfo `json:"userlist"`
		DepartmentList []int64        `json:"partylist"`
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/tag/get?tagid=" +
		strconv.FormatInt(id, 10) + "&access_token="
	if err = ((*corp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	userList = result.UserList
	departmentList = result.DepartmentList
	return
}

// 增加标签成员
func (clt *Client) TagAddUser(id int64, userList []string,
	departmentList []int64) (invalidUserList []string, invalidDepartmentList []int64, err error) {

	if len(userList) <= 0 && len(departmentList) <= 0 {
		return
	}

	var request = struct {
		TagId          int64    `json:"tagid"`
		Userlist       []string `json:"userlist,omitempty"`
		DepartmentList []int64  `json:"partylist,omitempty"`
	}{
		TagId:          id,
		Userlist:       userList,
		DepartmentList: departmentList,
	}

	var result struct {
		corp.Error
		InvalidUserList       string  `json:"invalidlist"`
		InvalidDepartmentList []int64 `json:"invalidparty"`
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/tag/addtagusers?access_token="
	if err = ((*corp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	switch result.ErrCode {
	case corp.ErrCodeOK:
		if result.InvalidUserList != "" {
			invalidUserList = strings.Split(result.InvalidUserList, "|")
		}
		invalidDepartmentList = result.InvalidDepartmentList
		return
	case 40070:
		invalidUserList = userList
		invalidDepartmentList = departmentList
		return
	default:
		err = &result.Error
		return
	}
}

// 删除标签成员
func (clt *Client) TagDeleteUser(id int64, userList []string,
	departmentList []int64) (invalidUserList []string, invalidDepartmentList []int64, err error) {

	if len(userList) <= 0 && len(departmentList) <= 0 {
		return
	}

	var request = struct {
		TagId          int64    `json:"tagid"`
		Userlist       []string `json:"userlist,omitempty"`
		DepartmentList []int64  `json:"partylist,omitempty"`
	}{
		TagId:          id,
		Userlist:       userList,
		DepartmentList: departmentList,
	}

	var result struct {
		corp.Error
		InvalidUserList       string  `json:"invalidlist"`
		InvalidDepartmentList []int64 `json:"invalidparty"`
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/tag/deltagusers?access_token="
	if err = ((*corp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	switch result.ErrCode {
	case corp.ErrCodeOK:
		if result.InvalidUserList != "" {
			invalidUserList = strings.Split(result.InvalidUserList, "|")
		}
		invalidDepartmentList = result.InvalidDepartmentList
		return
	case 40031:
		invalidUserList = userList
		invalidDepartmentList = departmentList
		return
	default:
		err = &result.Error
		return
	}
}

type Tag struct {
	Id   int64  `json:"tagid"`
	Name string `json:"tagname"`
}

// 获取标签列表
func (clt *Client) TagList() (list []Tag, err error) {
	var result struct {
		corp.Error
		TagList []Tag `json:"taglist"`
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/tag/list?access_token="
	if err = ((*corp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.TagList
	return
}
