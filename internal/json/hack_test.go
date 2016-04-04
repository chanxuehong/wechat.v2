// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"testing"
)

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
