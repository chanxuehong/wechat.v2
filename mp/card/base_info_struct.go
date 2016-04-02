// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package card

const (
	// 卡券Code展示类型
	CodeTypeText        = "CODE_TYPE_TEXT"         // 文本
	CodeTypeBarCode     = "CODE_TYPE_BARCODE"      // 一维码
	CodeTypeQRCode      = "CODE_TYPE_QRCODE"       // 二维码
	CodeTypeOnlyBarCode = "CODE_TYPE_ONLY_BARCODE" // 一维码无code显示
	CodeTypeOnlyQRCode  = "CODE_TYPE_ONLY_QRCODE"  // 二维码无code显示
)

const (
	// 卡券的状态
	CardStatusNotVerify    = "CARD_STATUS_NOT_VERIFY"    // 待审核
	CardStatusVerifyFail   = "CARD_STATUS_VERIFY_FALL"   // 审核失败
	CardStatusVerifyOk     = "CARD_STATUS_VERIFY_OK"     // 通过审核
	CardStatusUserDelete   = "CARD_STATUS_USER_DELETE"   // 卡券被用户删除
	CardStatusUserDispatch = "CARD_STATUS_USER_DISPATCH" // 在公众平台投放过的卡券
)

const (
	// DateInfo类型
	DateInfoTypeFixTimeRange = "DATE_TYPE_FIX_TIME_RANGE" // 表示固定日期区间
	DateInfoTypeFixTerm      = "DATE_TYPE_FIX_TERM"       // 表示固定时长（自领取后按天算）
	DateInfoTypePermanent    = "DATE_TYPE_PERMANENT"      // 表示永久有效（会员卡类型专用）
)

// 基本的卡券数据, 所有卡券通用
type CardBaseInfo struct {
	CardId string `json:"id,omitempty"`     // 查询的时候有返回
	Status string `json:"status,omitempty"` // 查询的时候有返回
	AppId  string `json:"appid,omitempty"`  // 查询的时候有返回

	LogoURL     string    `json:"logo_url,omitempty"`    // 卡券的商户logo，建议像素为300*300。
	CodeType    string    `json:"code_type,omitempty"`   // Code展示类型
	BrandName   string    `json:"brand_name,omitempty"`  // 商户名字,字数上限为12个汉字。
	Title       string    `json:"title,omitempty"`       // 卡券名，字数上限为9个汉字。(建议涵盖卡券属性、服务及金额)。
	SubTitle    string    `json:"sub_title,omitempty"`   // 券名的副标题, 字数上限为18个汉字。
	Color       string    `json:"color,omitempty"`       // 券颜色。按色彩规范标注填写Color010-Color100。
	Notice      string    `json:"notice,omitempty"`      // 卡券使用提醒，字数上限为16个汉字。
	Description string    `json:"description,omitempty"` // 卡券使用说明，字数上限为1024个汉字。
	SKU         *SKU      `json:"sku,omitempty"`         // 商品信息。
	DateInfo    *DateInfo `json:"date_info,omitempty"`   // 使用日期，有效期的信息。

	UseCustomCode        *bool   `json:"use_custom_code,omitempty"`         // 是否自定义Code码。填写true或false，默认为false。通常自有优惠码系统的开发者选择自定义Code码，并在卡券投放时带入Code码
	BindOpenId           *bool   `json:"bind_openid,omitempty"`             // 是否指定用户领取，填写true或false。默认为false。通常指定特殊用户群体投放卡券或防止刷券时选择指定用户领取。
	ServicePhone         string  `json:"service_phone,omitempty"`           // 客服电话。
	LocationIdList       []int64 `json:"location_id_list,omitempty"`        // 门店位置poiid。
	Source               string  `json:"source,omitempty"`                  // 第三方来源名，例如同程旅游、大众点评。
	CustomURLName        string  `json:"custom_url_name,omitempty"`         // 自定义跳转外链的入口名字。
	CustomURLSubTitle    string  `json:"custom_url_sub_title,omitempty"`    // 显示在入口右侧的提示语。
	CustomURL            string  `json:"custom_url,omitempty"`              // 自定义跳转的URL。
	PromotionURLName     string  `json:"promotion_url_name,omitempty"`      // 营销场景的自定义入口名称。
	PromotionURLSubTitle string  `json:"promotion_url_sub_title,omitempty"` // 显示在营销入口右侧的提示语。
	PromotionURL         string  `json:"promotion_url,omitempty"`           // 入口跳转外链的地址链接。
	GetLimit             *int    `json:"get_limit,omitempty"`               // 每人可领券的数量限制,不填写默认为50。
	UseLimit             *int    `json:"use_limit,omitempty"`               // 每人使用次数限制.
	CanShare             *bool   `json:"can_share,omitempty"`               // 卡券领取页面是否可分享。
	CanGiveFriend        *bool   `json:"can_give_friend,omitempty"`         // 卡券是否可转赠。
}

type SKU struct {
	Quantity int `json:"quantity"` // 卡券库存的数量，不支持填写0，上限为100000000。
}

type DateInfo struct {
	Type           string `json:"type"`                       // 使用时间的类型, DATE_TYPE_FIX_TIME_RANGE 表示固定日期区间，DATE_TYPE_FIX_TERM表示固定时长（自领取后按天算）
	BeginTimestamp int64  `json:"begin_timestamp,omitempty"`  // type为DATE_TYPE_FIX_TIME_RANGE时专用，表示起用时间。从1970年1月1日00:00:00至起用时间的秒数，最终需转换为字符串形态传入。（东八区时间，单位为秒）
	EndTimestamp   int64  `json:"end_timestamp,omitempty"`    // type为DATE_TYPE_FIX_TIME_RANGE时专用，表示结束时间，建议设置为截止日期的23:59:59过期。（东八区时间，单位为秒）
	FixedTerm      *int   `json:"fixed_term,omitempty"`       // type为DATE_TYPE_FIX_TERM时专用，表示自领取后多少天内有效，领取后当天有效填写0。（单位为天）
	FixedBeginTerm *int   `json:"fixed_begin_term,omitempty"` // type为DATE_TYPE_FIX_TERM时专用，表示自领取后多少天开始生效。（单位为天）
}
