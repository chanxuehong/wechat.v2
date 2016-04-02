package util

// ToLower 返回字符串对应的小写版本.
//  如果确定 s 是[a-zA-Z0-9_-] 的组合, 那么可以用这个函数, 否则请用 strings.ToLower!
func ToLower(s string) string {
	var b []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c > 'Z' || c < 'A' {
			if b != nil {
				b[i] = c
			}
		} else {
			c += 'a' - 'A'
			if b == nil {
				b = make([]byte, len(s))
				copy(b, s[:i])
			}
			b[i] = c
		}
	}
	if b != nil {
		return string(b)
	}
	return s
}
