package mch

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"hash"
	"sort"
	"testing"
)

var (
	testSignParameters = map[string]string{
		"asjdfsadfsadfasdfasd1":  "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd2":  "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd3":  "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd4":  "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd5":  "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd6":  "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd7":  "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd8":  "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd9":  "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd11": "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd12": "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd13": "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd14": "sdfasdfasdfsadfsadfasdfasdfasd",
		"asjdfsadfsadfasdfasd15": "sdfasdfasdfsadfsadfasdfasdfasd",
	}

	testAPIKey = "afadskfjaskldjflkasdjflkashdljkfhalsdjkfhl"
)

func BenchmarkSign(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sign(testSignParameters, testAPIKey, nil)
	}
}

func BenchmarkSign2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sign2(testSignParameters, testAPIKey, nil)
	}
}

// 传统的签名代码, Sign 是优化后的代码, 要提高 35% 的速度
func Sign2(parameters map[string]string, apiKey string, fn func() hash.Hash) string {
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
	signature := make([]byte, h.Size()*2)

	for _, k := range ks {
		v := parameters[k]
		if v == "" {
			continue
		}
		h.Write([]byte(k))
		h.Write([]byte{'='})
		h.Write([]byte(v))
		h.Write([]byte{'&'})
	}
	h.Write([]byte("key="))
	h.Write([]byte(apiKey))

	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
}
