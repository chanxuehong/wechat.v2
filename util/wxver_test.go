package util

import (
	"testing"
)

func TestWXVersion(t *testing.T) {
	userAgent := `Mozilla/5.0(iphone;CPU iphone OS 5_1_1 like Mac OS X) AppleWebKit/534.46(KHTML,like Geocko)Mobile/9B206 MicroMessenger/5`
	x, y, z, err := WXVersion(userAgent)
	if err != nil {
		t.Error(err)
		return
	}
	if x != 5 || y != 0 || z != 0 {
		t.Error("获取了错误的版本号")
		return
	}

	userAgent = `Mozilla/5.0(iphone;CPU iphone OS 5_1_1 like Mac OS X) AppleWebKit/534.46(KHTML,like Geocko)Mobile/9B206 MicroMessenger/5.3`
	x, y, z, err = WXVersion(userAgent)
	if err != nil {
		t.Error(err)
		return
	}
	if x != 5 || y != 3 || z != 0 {
		t.Error("获取了错误的版本号")
		return
	}

	userAgent = `Mozilla/5.0(iphone;CPU iphone OS 5_1_1 like Mac OS X) AppleWebKit/534.46(KHTML,like Geocko)Mobile/9B206 MicroMessenger/5.3.1`
	x, y, z, err = WXVersion(userAgent)
	if err != nil {
		t.Error(err)
		return
	}
	if x != 5 || y != 3 || z != 1 {
		t.Error("获取了错误的版本号")
		return
	}

	userAgent = `Mozilla/5.0(iphone;CPU iphone OS 5_1_1 like Mac OS X) AppleWebKit/534.46(KHTML,like Geocko)Mobile/9B206 MicroMessenger/5.3.1.5`
	x, y, z, err = WXVersion(userAgent)
	if err != nil {
		t.Error(err)
		return
	}
	if x != 5 || y != 3 || z != 1 {
		t.Error("获取了错误的版本号")
		return
	}

	userAgent = `Mozilla5.0(iphone;CPU iphone OS 5_1_1 like Mac OS X) AppleWebKit534.46(KHTML,like Geocko)Mobile9B206 MicroMessenger5.0`
	_, _, _, err = WXVersion(userAgent)
	if err == nil {
		t.Errorf("从 %#s 获取版本号应该出错, 但是目前却没有错误!", userAgent)
		return
	}

	userAgent = `Mozilla/5.0(iphone;CPU iphone OS 5_1_1 like Mac OS X) AppleWebKit/534.46(KHTML,like Geocko)Mobile/9B206 MicroMessenger/`
	_, _, _, err = WXVersion(userAgent)
	if err == nil {
		t.Errorf("从 %#s 获取版本号应该出错, 但是目前却没有错误!", userAgent)
		return
	}

	userAgent = `Mozilla/5.0(iphone;CPU iphone OS 5_1_1 like Mac OS X) AppleWebKit/534.46(KHTML,like Geocko)Mobile/9B206 MicroMessenger/5x`
	_, _, _, err = WXVersion(userAgent)
	if err == nil {
		t.Errorf("从 %#s 获取版本号应该出错, 但是目前却没有错误!", userAgent)
		return
	}

	userAgent = `Mozilla/5.0(iphone;CPU iphone OS 5_1_1 like Mac OS X) AppleWebKit/534.46(KHTML,like Geocko)Mobile/9B206 MicroMessenger/5.3x`
	_, _, _, err = WXVersion(userAgent)
	if err == nil {
		t.Errorf("从 %#s 获取版本号应该出错, 但是目前却没有错误!", userAgent)
		return
	}

	userAgent = `Mozilla/5.0(iphone;CPU iphone OS 5_1_1 like Mac OS X) AppleWebKit/534.46(KHTML,like Geocko)Mobile/9B206 MicroMessenger/5.3.1x`
	_, _, _, err = WXVersion(userAgent)
	if err == nil {
		t.Errorf("从 %#s 获取版本号应该出错, 但是目前却没有错误!", userAgent)
		return
	}
}
