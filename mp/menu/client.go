// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package menu

import (
	"net/http"

	"github.com/chanxuehong/wechat/mp"
)

type Client mp.Client

func NewClient(srv mp.AccessTokenServer, clt *http.Client) *Client {
	return (*Client)(mp.NewClient(srv, clt))
}

// 创建自定义菜单.
func (clt *Client) CreateMenu(menu Menu) (err error) {
	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/menu/create?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &menu, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 删除自定义菜单
func (clt *Client) DeleteMenu() (err error) {
	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/menu/delete?access_token="
	if err = ((*mp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 获取自定义菜单
func (clt *Client) GetMenu() (menu *Menu, conditionalMenus []Menu, err error) {
	var result struct {
		mp.Error
		Menu             Menu   `json:"menu"`
		ConditionalMenus []Menu `json:"conditionalmenu"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/menu/get?access_token="
	if err = ((*mp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	menu = &result.Menu
	conditionalMenus = result.ConditionalMenus
	return
}

// 获取自定义菜单配置接口
func (clt *Client) GetMenuInfo() (info MenuInfo, isMenuOpen bool, err error) {
	var result struct {
		mp.Error
		IsMenuOpen int      `json:"is_menu_open"`
		MenuInfo   MenuInfo `json:"selfmenu_info"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token="
	if err = ((*mp.Client)(clt)).GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = result.MenuInfo
	if result.IsMenuOpen == 1 {
		isMenuOpen = true
	}
	return
}

// 创建个性化菜单.
func (clt *Client) CreateConditionalMenu(menu *Menu) (menuId int64, err error) {
	var result struct {
		mp.Error
		MenuId int64 `json:"menuid"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/menu/addconditional?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, menu, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	menuId = result.MenuId
	return
}

// 删除个性化菜单.
func (clt *Client) DeleteConditionalMenu(menuId int64) (err error) {
	var request = struct {
		MenuId int64 `json:"menuid"`
	}{
		MenuId: menuId,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/menu/delconditional?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 测试个性化菜单匹配结果.
// userId 可以是粉丝的 OpenID, 也可以是粉丝的微信号
func (clt *Client) TryMatch(userId string) (menu *Menu, err error) {
	var request = struct {
		UserId string `json:"user_id"`
	}{
		UserId: userId,
	}
	var result struct {
		mp.Error
		Menu Menu `json:"menu"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/menu/trymatch?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	menu = &result.Menu
	return
}
