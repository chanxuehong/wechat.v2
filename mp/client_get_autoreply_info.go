// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mp

// 获取自动回复规则
func (clt *Client) GetAutoReplyInfo() (info *AutoReplyInfo, err error) {
	var result struct {
		Error
		AutoReplyInfo
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/get_current_autoreply_info?access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.AutoReplyInfo
	return
}

type AutoReplyInfo struct {
	IsAddFriendReplyOpen int `json:"is_add_friend_reply_open"` // 关注后自动回复是否开启，0代表未开启，1代表开启
	IsAutoReplyOpen      int `json:"is_autoreply_open"`        // 消息自动回复是否开启，0代表未开启，1代表开启

	AddFriendAutoReplyInfo      *AddFriendAutoReplyInfo      `json:"add_friend_autoreply_info,omitempty"`      // 关注后自动回复的信息
	MessageDefaultAutoReplyInfo *MessageDefaultAutoReplyInfo `json:"message_default_autoreply_info,omitempty"` // 消息自动回复的信息
	KeywordAutoReplyInfo        *KeywordAutoReplyInfo        `json:"keyword_autoreply_info,omitempty"`         // 关键词自动回复的信息
}

// 关注后自动回复的信息
type AddFriendAutoReplyInfo struct {
	Type    string `json:"type"`    // 自动回复的类型。关注后自动回复和消息自动回复的类型仅支持文本（text）、图片（img）、语音（voice）、视频（video），关键词自动回复则还多了图文消息
	Content string `json:"content"` // 对于文本类型，content是文本内容，对于图片、语音、视频类型，content是mediaID
}

// 消息自动回复的信息
type MessageDefaultAutoReplyInfo struct {
	Type    string `json:"type"`    // 自动回复的类型。关注后自动回复和消息自动回复的类型仅支持文本（text）、图片（img）、语音（voice）、视频（video），关键词自动回复则还多了图文消息
	Content string `json:"content"` // 对于文本类型，content是文本内容，对于图片、语音、视频类型，content是mediaID
}

// 关键词自动回复的信息
type KeywordAutoReplyInfo struct {
	RuleList []KeywordAutoReplyRule `json:"list,omitempty"`
}

// 关键词自动回复的规则
type KeywordAutoReplyRule struct {
	RuleName        string        `json:"rule_name"`                   // 规则名称
	CreateTime      int64         `json:"create_time"`                 // 创建时间
	ReplyMode       string        `json:"reply_mode"`                  // 回复模式，reply_all代表全部回复，random_one代表随机回复其中一条
	KeywordInfoList []KeywordInfo `json:"keyword_list_info,omitempty"` // 匹配的关键词列表
	ReplyInfoList   []ReplyInfo   `json:"reply_list_info,omitempty"`   // 回复信息列表
}

// 关键词匹配规则
type KeywordInfo struct {
	Type      string `json:"type"` // 一般都是文本吧???
	Content   string `json:"content"`
	MatchMode string `json:"match_mode"` // 匹配模式，contain代表消息中含有该关键词即可，equal表示消息内容必须和关键词严格相同
}

// 关键词回复信息
type ReplyInfo struct {
	Type string `json:"type"` // 自动回复的类型。关注后自动回复和消息自动回复的类型仅支持文本（text）、图片（img）、语音（voice）、视频（video），关键词自动回复则还多了图文消息

	// 下面两个字段不会同时有效, 根据 ReplyInfo.Type 来做选择
	Content  string    `json:"content,omitempty"`   // 对于文本类型，content是文本内容，对于图片、语音、视频类型，content是mediaID
	NewsInfo *NewsInfo `json:"news_info,omitempty"` // 图文消息的信息
}

type NewsInfo struct {
	ArticleList []Article `json:"list,omitempty"`
}

type Article struct {
	Title      string `json:"title"`       // 图文消息的标题
	Author     string `json:"author"`      // 作者
	Digest     string `json:"digest"`      // 摘要
	ShowCover  int    `json:"show_cover"`  // 是否显示封面，0为不显示，1为显示
	CoverURL   string `json:"cover_url"`   // 封面图片的URL
	ContentURL string `json:"content_url"` // 正文的URL
	SourceURL  string `json:"source_url"`  // 原文的URL，若置空则无查看原文入口
}
