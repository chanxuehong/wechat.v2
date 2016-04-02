// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mch

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"sort"
)

// 微信支付签名.
//  parameters: 待签名的参数集合
//  apiKey:     API密钥
//  fn:         func() hash.Hash, 如果 fn == nil 则默认用 md5.New
func Sign(parameters map[string]string, apiKey string, fn func() hash.Hash) string {
	ks := make([]string, 0, len(parameters))
	for k := range parameters {
		if k == "sign" {
			continue
		}
		ks = append(ks, k)
	}
	sort.Strings(ks)

	if fn == nil {
		fn = md5.New
	}
	h := fn()

	buf := make([]byte, 256)
	for _, k := range ks {
		v := parameters[k]
		if v == "" {
			continue
		}

		buf = buf[:0]
		buf = append(buf, k...)
		buf = append(buf, '=')
		buf = append(buf, v...)
		buf = append(buf, '&')
		h.Write(buf)
	}
	buf = buf[:0]
	buf = append(buf, "key="...)
	buf = append(buf, apiKey...)
	h.Write(buf)

	signature := make([]byte, h.Size()*2)
	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
}

// 收货地址共享接口签名
func EditAddressSign(appId, url, timestamp, nonceStr, accessToken string) string {
	h := sha1.New()
	buf := make([]byte, 256)

	// accesstoken
	// appid
	// noncestr
	// timestamp
	// url
	buf = buf[:0]
	buf = append(buf, "accesstoken="...)
	buf = append(buf, accessToken...)
	h.Write(buf)

	buf = buf[:0]
	buf = append(buf, "&appid="...)
	buf = append(buf, appId...)
	h.Write(buf)

	buf = buf[:0]
	buf = append(buf, "&noncestr="...)
	buf = append(buf, nonceStr...)
	h.Write(buf)

	buf = buf[:0]
	buf = append(buf, "&timestamp="...)
	buf = append(buf, timestamp...)
	h.Write(buf)

	buf = buf[:0]
	buf = append(buf, "&url="...)
	buf = append(buf, url...)
	h.Write(buf)

	return hex.EncodeToString(h.Sum(nil))
}
