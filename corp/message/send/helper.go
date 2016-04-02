// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package send

import (
	"strconv"
	"strings"
)

// 用 '|' 连接 a 的各个元素
func JoinString(a []string) string {
	return strings.Join(a, "|")
}

// 用 '|' 连接 a 的各个元素的十进制字符串
func JoinInt64(a []int64) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(a[0], 10)
	default:
		strs := make([]string, len(a))
		for i, n := range a {
			strs[i] = strconv.FormatInt(n, 10)
		}
		return strings.Join(strs, "|")
	}
}

// 用 '|' 分离 str
func SplitString(str string) []string {
	return strings.Split(str, "|")
}

// 用 '|' 分离 str, 然后将分离后的字符串都转换为整数
//  NOTE: 要求 str 都是整数合并的, 否则会出错
func SplitInt64(str string) (dst []int64, err error) {
	strs := strings.Split(str, "|")

	dst = make([]int64, len(strs))
	for i, str := range strs {
		dst[i], err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			return
		}
	}
	return
}
