// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package menu

const (
	MenuButtonCountLimit    = 3 // 一级菜单最多包含 3 个按钮
	SubMenuButtonCountLimit = 5 // 二级菜单最多包含 5 个按钮
)

const (
	MenuButtonNameLenLimit    = 16 // 菜单标题不超过16个字节
	SubMenuButtonNameLenLimit = 40 // 子菜单标题不超过40个字节
)

const (
	ButtonKeyLenLimit = 128 // 菜单KEY值不能超过128字节
	ButtonURLLenLimit = 256 // 网页链接不能超过256字节
)

const (
	// 下面六個類型(包括view)是在公众平台官网通过网站功能发布菜单的按鈕類型
	ButtonTypeText  = "text"
	ButtonTypeImage = "img"
	ButtonTypePhoto = "photo"
	ButtonTypeVideo = "video"
	ButtonTypeVoice = "voice"

	// NOTE: 上面的按鈕類型不能通過API設置

	ButtonTypeView  = "view"  // 跳转URL
	ButtonTypeClick = "click" // 点击推事件

	// 仅支持微信iPhone5.4.1以上版本, 和Android5.4以上版本的微信用户,
	// 旧版本微信用户点击后将没有回应, 开发者也不能正常接收到事件推送.
	ButtonTypeScanCodePush    = "scancode_push"      // 扫码推事件
	ButtonTypeScanCodeWaitMsg = "scancode_waitmsg"   // 扫码带提示
	ButtonTypePicSysPhoto     = "pic_sysphoto"       // 系统拍照发图
	ButtonTypePicPhotoOrAlbum = "pic_photo_or_album" // 拍照或者相册发图
	ButtonTypePicWeixin       = "pic_weixin"         // 微信相册发图
	ButtonTypeLocationSelect  = "location_select"    // 发送位置

	// 专门给第三方平台旗下未微信认证(具体而言, 是资质认证未通过)的订阅号准备的事件类型,
	// 它们是没有事件推送的, 能力相对受限, 其他类型的公众号不必使用.
	ButtonTypeMediaId     = "media_id"     // 下发消息
	ButtonTypeViewLimited = "view_limited" // 跳转图文消息URL
)

type Menu struct {
	Buttons   []Button   `json:"button,omitempty"` // 一级菜单数组, 个数应为1~3个
	MatchRule *MatchRule `json:"matchrule,omitempty"`
	MenuId    int64      `json:"menuid,omitempty"` // 有个性化菜单时查询接口返回值包含这个字段
}

// 菜单的按钮
type Button struct {
	Type       string   `json:"type,omitempty"`       // 非必须; 菜单的响应动作类型
	Name       string   `json:"name,omitempty"`       // 必须;  菜单标题, 不超过16个字节, 子菜单不超过40个字节
	Key        string   `json:"key,omitempty"`        // 非必须; 菜单KEY值, 用于消息接口推送, 不超过128字节
	URL        string   `json:"url,omitempty"`        // 非必须; 网页链接, 用户点击菜单可打开链接, 不超过256字节
	MediaId    string   `json:"media_id,omitempty"`   // 非必须; 调用新增永久素材接口返回的合法media_id
	SubButtons []Button `json:"sub_button,omitempty"` // 非必须; 二级菜单数组, 个数应为1~5个
}

type MatchRule struct {
	GroupId            *int64 `json:"group_id,omitempty"`
	Sex                *int   `json:"sex,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	ClientPlatformType *int   `json:"client_platform_type,omitempty"`
	Language           string `json:"language,omitempty"`
}

// 设置 btn 指向的 Button 为 子菜单 类型按钮
func (btn *Button) SetAsSubMenuButton(name string, subButtons []Button) {
	btn.Name = name
	btn.SubButtons = subButtons

	btn.Type = ""
	btn.Key = ""
	btn.URL = ""
	btn.MediaId = ""
}

// 设置 btn 指向的 Button 为 click 类型按钮
func (btn *Button) SetAsClickButton(name, key string) {
	btn.Type = ButtonTypeClick
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 view 类型按钮
func (btn *Button) SetAsViewButton(name, url string) {
	btn.Type = ButtonTypeView
	btn.Name = name
	btn.URL = url

	btn.Key = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 扫码推事件 类型按钮
func (btn *Button) SetAsScanCodePushButton(name, key string) {
	btn.Type = ButtonTypeScanCodePush
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 扫码推事件且弹出"消息接收中"提示框 类型按钮
func (btn *Button) SetAsScanCodeWaitMsgButton(name, key string) {
	btn.Type = ButtonTypeScanCodeWaitMsg
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 弹出系统拍照发图 类型按钮
func (btn *Button) SetAsPicSysPhotoButton(name, key string) {
	btn.Type = ButtonTypePicSysPhoto
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 弹出拍照或者相册发图 类型按钮
func (btn *Button) SetAsPicPhotoOrAlbumButton(name, key string) {
	btn.Type = ButtonTypePicPhotoOrAlbum
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 弹出微信相册发图器 类型按钮
func (btn *Button) SetAsPicWeixinButton(name, key string) {
	btn.Type = ButtonTypePicWeixin
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 弹出地理位置选择器 类型按钮
func (btn *Button) SetAsLocationSelectButton(name, key string) {
	btn.Type = ButtonTypeLocationSelect
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 下发消息(除文本消息) 类型按钮
func (btn *Button) SetAsMediaIdButton(name, mediaId string) {
	btn.Type = ButtonTypeMediaId
	btn.Name = name
	btn.MediaId = mediaId

	btn.Key = ""
	btn.URL = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 跳转图文消息URL 类型按钮
func (btn *Button) SetAsViewLimitedButton(name, mediaId string) {
	btn.Type = ButtonTypeViewLimited
	btn.Name = name
	btn.MediaId = mediaId

	btn.Key = ""
	btn.URL = ""
	btn.SubButtons = nil
}
