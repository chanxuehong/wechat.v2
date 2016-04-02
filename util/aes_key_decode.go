package util

import (
	"encoding/base64"
	"errors"
)

// 把长度为 43 的字符串 base64 decode 到 32 字节的 []byte
//  encodedAESKey 由 a-z,A-Z,0-9 组成, 一般在微信管理后台随机生成
func AESKeyDecode(encodedAESKey string) (aesKey []byte, err error) {
	if len(encodedAESKey) != 43 {
		err = errors.New("the length of encodedAESKey must be equal to 43")
		return
	}
	return base64.StdEncoding.DecodeString(encodedAESKey + "=")
}
