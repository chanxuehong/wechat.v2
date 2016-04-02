// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package card

const (
	// 卡券类型
	CardTypeGeneralCoupon = "GENERAL_COUPON" // 优惠券
	CardTypeGroupon       = "GROUPON"        // 团购券
	CardTypeCash          = "CASH"           // 代金券
	CardTypeDiscount      = "DISCOUNT"       // 折扣券
	CardTypeGift          = "GIFT"           // 礼品券
	CardTypeMemberCard    = "MEMBER_CARD"    // 会员卡
	CardTypeMeetingTicket = "MEETING_TICKET" // 会议门票
	CardTypeScenicTicket  = "SCENIC_TICKET"  // 景区门票
	CardTypeMovieTicket   = "MOVIE_TICKET"   // 电影票
	CardTypeBoardingPass  = "BOARDING_PASS"  // 飞机票
)

// 卡券数据结构
type Card struct {
	CardType string `json:"card_type,omitempty"`

	GeneralCoupon *GeneralCoupon `json:"general_coupon,omitempty"`
	Groupon       *Groupon       `json:"groupon,omitempty"`
	Cash          *Cash          `json:"cash,omitempty"`
	Discount      *Discount      `json:"discount,omitempty"`
	Gift          *Gift          `json:"gift,omitempty"`
	MemberCard    *MemberCard    `json:"member_card,omitempty"`
	MeetingTicket *MeetingTicket `json:"meeting_ticket,omitempty"`
	ScenicTicket  *ScenicTicket  `json:"scenic_ticket,omitempty"`
	MovieTicket   *MovieTicket   `json:"movie_ticket,omitempty"`
	BoardingPass  *BoardingPass  `json:"boarding_pass,omitempty"`
}

// 优惠券
type GeneralCoupon struct {
	BaseInfo      *CardBaseInfo `json:"base_info,omitempty"`
	DefaultDetail string        `json:"default_detail,omitempty"` // 优惠券专用, 填写优惠详情
}

// 团购券
type Groupon struct {
	BaseInfo   *CardBaseInfo `json:"base_info,omitempty"`
	DealDetail string        `json:"deal_detail,omitempty"` // 团购券专用，团购详情
}

// 代金券
type Cash struct {
	BaseInfo   *CardBaseInfo `json:"base_info,omitempty"`
	LeastCost  *int          `json:"least_cost,omitempty"`  // 代金券专用, 表示起用金额(单位为分)
	ReduceCost *int          `json:"reduce_cost,omitempty"` // 代金券专用, 表示减免金额(单位为分)
}

// 折扣券
type Discount struct {
	BaseInfo *CardBaseInfo `json:"base_info,omitempty"`
	Discount *int          `json:"discount,omitempty"` // 折扣券专用, 表示打折额度(百分比). 填30 就是七折.
}

// 礼品券
type Gift struct {
	BaseInfo *CardBaseInfo `json:"base_info,omitempty"`
	Gift     string        `json:"gift,omitempty"` // 礼品券专用, 表示礼品名字
}

// 会员卡
type MemberCard struct {
	BaseInfo *CardBaseInfo `json:"base_info,omitempty"`

	Prerogative       string                 `json:"prerogative,omitempty"`       // 会员卡特权说明
	SupplyBonus       *bool                  `json:"supply_bonus,omitempty"`      // 显示积分，填写true或false，如填写true，积分相关字段均为必填
	BonusURL          string                 `json:"bonus_url,omitempty"`         // 设置跳转外链查看积分详情。仅适用于积分无法通过激活接口同步的情况下使用该字段。
	SupplyBalance     *bool                  `json:"supply_balance,omitempty"`    // 是否支持储值，填写true或false。如填写true，储值相关字段均为必填。
	BalanceURL        string                 `json:"balance_url,omitempty"`       // 设置跳转外链查看余额详情。仅适用于余额无法通过激活接口同步的情况下使用该字段。
	BonusClearedRules string                 `json:"bonus_cleared,omitempty"`     // 积分清零规则。
	BonusRules        string                 `json:"bonus_rules,omitempty"`       // 积分规则。
	BalanceRules      string                 `json:"balance_rules,omitempty"`     // 储值说明。
	ActivateURL       string                 `json:"activate_url,omitempty"`      // 激活会员卡的url。
	NeedPushOnView    *bool                  `json:"need_push_on_view,omitempty"` // 填写true为用户点击进入会员卡时推送事件，默认为false。
	CustomField1      *MemberCardCustomField `json:"custom_field1,omitempty"`     // 自定义会员信息类目，会员卡激活后显示。
	CustomField2      *MemberCardCustomField `json:"custom_field2,omitempty"`     // 自定义会员信息类目，会员卡激活后显示。
	CustomField3      *MemberCardCustomField `json:"custom_field3,omitempty"`     // 自定义会员信息类目，会员卡激活后显示。
	CustomCell1       *MemberCardCustomCell  `json:"custom_cell1,omitempty"`      // 自定义会员信息类目，会员卡激活后显示。
}

type MemberCardCustomField struct {
	// 会员信息类目名称:
	//
	// FIELD_NAME_TYPE_LEVEL        等级
	// FIELD_NAME_TYPE_COUPON       优惠券
	// FIELD_NAME_TYPE_STAMP        印花
	// FIELD_NAME_TYPE_DISCOUNT     折扣
	// FIELD_NAME_TYPE_ACHIEVEMEN   成就
	// FIELD_NAME_TYPE_MILEAGE      里程
	NameType string `json:"name_type,omitempty"`
	URL      string `json:"url,omitempty"` // 点击类目跳转外链url
}

type MemberCardCustomCell struct {
	Name string `json:"name,omitempty"` // 入口名称。
	Tips string `json:"tips,omitempty"` // 入口右侧提示语，6个汉字内。
	URL  string `json:"url,omitempty"`  // 入口跳转链接。
}

// 会议门票
type MeetingTicket struct {
	BaseInfo      *CardBaseInfo `json:"base_info,omitempty"`
	MeetingDetail string        `json:"meeting_detail,omitempty"` // 会议详情
	MapURL        string        `json:"map_url,omitempty"`        // 会议导览图
}

// 景区门票
type ScenicTicket struct {
	BaseInfo    *CardBaseInfo `json:"base_info,omitempty"`
	TicketClass string        `json:"ticket_class,omitempty"` // 票类型, 例如平日全票, 套票等
	GuideURL    string        `json:"guide_url,omitempty"`    // 导览图url
}

// 电影票
type MovieTicket struct {
	BaseInfo *CardBaseInfo `json:"base_info,omitempty"`
	Detail   string        `json:"detail,omitempty"` // 电影票详情
}

// 飞机票
type BoardingPass struct {
	BaseInfo *CardBaseInfo `json:"base_info,omitempty"`

	From          string `json:"from,omitempty"`           // 起点, 上限为18 个汉字
	To            string `json:"to,omitempty"`             // 终点, 上限为18 个汉字
	Flight        string `json:"flight,omitempty"`         // 航班
	Gate          string `json:"gate,omitempty"`           // 登机口. 如发生登机口变更, 建议商家实时调用该接口变更
	CheckinURL    string `json:"check_in_url,omitempty"`   // 在线值机的链接
	AirModel      string `json:"air_model,omitempty"`      // 机型, 上限为8 个汉字
	DepartureTime int64  `json:"departure_time,omitempty"` // 起飞时间. Unix 时间戳格式
	LandingTime   int64  `json:"landing_time,omitempty"`   // 降落时间. Unix 时间戳格式
	BoardingTime  string `json:"boarding_time,omitempty"`  // 登机时间, 只显示"时分"不显示日期, 按时间戳格式填写. 如发生登机时间变更, 建议商家实时调用该接口变更.
}
