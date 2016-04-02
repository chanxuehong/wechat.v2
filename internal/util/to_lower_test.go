package util

import (
	"strings"
	"testing"
)

func TestToLower(t *testing.T) {
	strs := []string{
		"aaaa_bbbb-cccc",
		"aaaA_Bbbb-CCCC",
		"AAAA_BBBB-cccc",
	}
	for _, str := range strs {
		dst1 := ToLower(str)
		dst2 := strings.ToLower(str)
		if dst1 != dst2 {
			t.Errorf("TestToLower failed, have: %s, want %s\n", dst1, dst2)
		}
	}
}

func BenchmarkToLower(b *testing.B) {
	s := "scancode_waitmsg"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToLower(s)
	}
}

func BenchmarkStringsToLower(b *testing.B) {
	s := "scancode_waitmsg"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strings.ToLower(s)
	}
}
