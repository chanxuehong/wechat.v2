// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"reflect"
	"regexp"
	"testing"
)

func TestUnmarshalStringToNumberOrBoolean(t *testing.T) {
	var x = struct {
		A []string `json:"a"`
		B string   `json:"b"`
		C string   `json:"c"`
		D []string `json:"d"`
	}{
		A: []string{"1", "2"},
		B: "true",
		C: "100",
		D: []string{"true", "false"},
	}
	b, err := Marshal(&x)
	if err != nil {
		t.Error(err)
		return
	}

	type T struct {
		A []int  `json:"a"`
		B bool   `json:"b"`
		C int    `json:"c"`
		D []bool `json:"d"`
	}
	var y T
	if err := Unmarshal(b, &y); err != nil {
		t.Error(err)
		return
	}

	var y2 = T{
		A: []int{1, 2},
		B: true,
		C: 100,
		D: []bool{true, false},
	}
	if !reflect.DeepEqual(y, y2) {
		t.Error("TestUnmarshal2 failed")
		return
	}
}

func TestNumberBytesIsValid(t *testing.T) {
	// From: http://stackoverflow.com/a/13340826
	var jsonNumberRegexp = regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?$`)

	validTests := []string{
		"0",
		"-0",
		"1",
		"-1",
		"0.1",
		"-0.1",
		"1234",
		"-1234",
		"12.34",
		"-12.34",
		"12E0",
		"12E1",
		"12e34",
		"12E-0",
		"12e+1",
		"12e-34",
		"-12E0",
		"-12E1",
		"-12e34",
		"-12E-0",
		"-12e+1",
		"-12e-34",
		"1.2E0",
		"1.2E1",
		"1.2e34",
		"1.2E-0",
		"1.2e+1",
		"1.2e-34",
		"-1.2E0",
		"-1.2E1",
		"-1.2e34",
		"-1.2E-0",
		"-1.2e+1",
		"-1.2e-34",
		"0E0",
		"0E1",
		"0e34",
		"0E-0",
		"0e+1",
		"0e-34",
		"-0E0",
		"-0E1",
		"-0e34",
		"-0E-0",
		"-0e+1",
		"-0e-34",
	}

	for _, test := range validTests {
		if !isValidNumberBytes([]byte(test)) {
			t.Errorf("%s should be valid", test)
		}

		var f float64
		if err := Unmarshal([]byte(test), &f); err != nil {
			t.Errorf("%s should be valid but Unmarshal failed: %v", test, err)
		}

		if !jsonNumberRegexp.MatchString(test) {
			t.Errorf("%s should be valid but regexp does not match", test)
		}
	}

	invalidTests := []string{
		"",
		"invalid",
		"1.0.1",
		"1..1",
		"-1-2",
		"012a42",
		"01.2",
		"012",
		"12E12.12",
		"1e2e3",
		"1e+-2",
		"1e--23",
		"1e",
		"e1",
		"1e+",
		"1ea",
		"1a",
		"1.a",
		"1.",
		"01",
		"1.e1",
	}

	for _, test := range invalidTests {
		if isValidNumberBytes([]byte(test)) {
			t.Errorf("%s should be invalid", test)
		}

		var f float64
		if err := Unmarshal([]byte(test), &f); err == nil {
			t.Errorf("%s should be invalid but unmarshal wrote %v", test, f)
		}

		if jsonNumberRegexp.MatchString(test) {
			t.Errorf("%s should be invalid but matches regexp", test)
		}
	}
}

func BenchmarkNumberBytesIsValid(b *testing.B) {
	s := []byte("-61657.61667E+61673")
	for i := 0; i < b.N; i++ {
		isValidNumberBytes(s)
	}
}

func TestUnmarshalControlCharacter(t *testing.T) {
	b := make([]byte, 0x20)
	for i := 0; i < len(b); i++ {
		b[i] = byte(i)
	}

	var data []byte
	data = append(data, '"')
	data = append(data, b...)
	data = append(data, '"')

	var s string
	if err := Unmarshal(data, &s); err != nil {
		t.Error(err)
		return
	}
	if string(b) != s {
		t.Error("TestUnmarshalControlCharacter failed")
		return
	}
}
