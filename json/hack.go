// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
)

var jsonNumberType = reflect.TypeOf(json.Number(""))

var trueLiteral = []byte("true")
var falseLiteral = []byte("false")

func unmarshalStringToNumberOrBoolean(d *decodeState, s []byte, v reflect.Value, fromQuoted bool) {
	if fromQuoted {
		d.saveError(&UnmarshalTypeError{"string", v.Type(), int64(d.off)})
		return
	}
	switch {
	default:
		d.saveError(&UnmarshalTypeError{"string", v.Type(), int64(d.off)})
	case isValidNumberBytes(s):
		switch v.Kind() {
		default:
			d.saveError(&UnmarshalTypeError{"string", v.Type(), int64(d.off)})
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			s := string(s)
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil || v.OverflowInt(n) {
				d.saveError(&UnmarshalTypeError{"string " + s, v.Type(), int64(d.off)})
				break
			}
			v.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			s := string(s)
			n, err := strconv.ParseUint(s, 10, 64)
			if err != nil || v.OverflowUint(n) {
				d.saveError(&UnmarshalTypeError{"string " + s, v.Type(), int64(d.off)})
				break
			}
			v.SetUint(n)
		case reflect.Float32, reflect.Float64:
			s := string(s)
			n, err := strconv.ParseFloat(s, v.Type().Bits())
			if err != nil || v.OverflowFloat(n) {
				d.saveError(&UnmarshalTypeError{"string " + s, v.Type(), int64(d.off)})
				break
			}
			v.SetFloat(n)
		}
	case bytes.Equal(s, trueLiteral):
		switch v.Kind() {
		default:
			d.saveError(&UnmarshalTypeError{"string", v.Type(), int64(d.off)})
		case reflect.Bool:
			v.SetBool(true)
		}
	case bytes.Equal(s, falseLiteral):
		switch v.Kind() {
		default:
			d.saveError(&UnmarshalTypeError{"string", v.Type(), int64(d.off)})
		case reflect.Bool:
			v.SetBool(false)
		}
	}
}

// isValidNumberBytes reports whether s is a valid JSON number literal.
func isValidNumberBytes(s []byte) bool {
	// This function implements the JSON numbers grammar.
	// See https://tools.ietf.org/html/rfc7159#section-6
	// and http://json.org/number.gif

	if len(s) == 0 {
		return false
	}

	// Optional -
	if s[0] == '-' {
		s = s[1:]
		if len(s) == 0 {
			return false
		}
	}

	// Digits
	switch {
	default:
		return false

	case s[0] == '0':
		s = s[1:]

	case '1' <= s[0] && s[0] <= '9':
		s = s[1:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// . followed by 1 or more digits.
	if len(s) >= 2 && s[0] == '.' && '0' <= s[1] && s[1] <= '9' {
		s = s[2:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// e or E followed by an optional - or + and
	// 1 or more digits.
	if len(s) >= 2 && (s[0] == 'e' || s[0] == 'E') {
		s = s[1:]
		if s[0] == '+' || s[0] == '-' {
			s = s[1:]
			if len(s) == 0 {
				return false
			}
		}
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// Make sure we are at the end.
	return len(s) == 0
}
