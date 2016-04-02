package util

import (
	"fmt"
	"strconv"
	"strings"
)

// 获取微信客户端的版本.
//  userAgent: 微信内置浏览器的 user-agent;
//  x, y, z:   如果微信版本为 5.3.1 则有 x==5, y==3, z==1
//  err:       错误信息
func WXVersion(userAgent string) (x, y, z int, err error) {
	// Mozilla/5.0(iphone;CPU iphone OS 5_1_1 like Mac OS X) AppleWebKit/534.46(KHTML,like Geocko)Mobile/9B206 MicroMessenger/5.0
	lastSlashIndex := strings.LastIndex(userAgent, "/")
	versionIndex := lastSlashIndex + 1

	if lastSlashIndex == -1 || versionIndex >= len(userAgent) {
		err = fmt.Errorf("不是有效的微信浏览器 user-agent, %s", userAgent)
		return
	}

	strArr := strings.Split(userAgent[versionIndex:], ".")
	verArr := make([]int, len(strArr))

	for i, str := range strArr {
		verArr[i], err = strconv.Atoi(str)
		if err != nil {
			err = fmt.Errorf("不是有效的微信浏览器 user-agent, %s", userAgent)
			return
		}
	}

	// len(verArr) == len(strArr) >= 1, why?
	switch len(verArr) {
	case 3:
		x = verArr[0]
		y = verArr[1]
		z = verArr[2]
	case 2:
		x = verArr[0]
		y = verArr[1]
	case 1:
		x = verArr[0]
	default:
		x = verArr[0]
		y = verArr[1]
		z = verArr[2]
	}
	return
}
