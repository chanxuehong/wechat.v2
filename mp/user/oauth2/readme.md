# 网页授权获取用户基本信息

```Go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/session"
	"github.com/chanxuehong/sid"
	"github.com/chanxuehong/wechat/mp/user/oauth2"
)

var (
	sessionStorage = session.New(20*60, 60*60)

	oauth2Config = oauth2.NewOAuth2Config(
		"appid",                     // 填上自己的参数
		"appsecret",                 // 填上自己的参数
		"http://192.168.1.80/page2", // 授权后跳转地址
		"snsapi_userinfo",           // 需要用户授权, snsapi_base 不需要
	)
)

func init() {
	http.HandleFunc("/page1", Page1Handler)
	http.HandleFunc("/page2", Page2Handler)
}

// 建立必要的 session, 然后跳转到授权页面
func Page1Handler(w http.ResponseWriter, r *http.Request) {
	state := string(rand.NewHex())
	sid := sid.New()

	if err := sessionStorage.Add(sid, state); err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	cookie := http.Cookie{
		Name:     "sid",
		Value:    sid,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	AuthCodeURL := oauth2Config.AuthCodeURL(state, nil)
	log.Println("AuthCodeURL:", AuthCodeURL)

	http.Redirect(w, r, AuthCodeURL, http.StatusFound)
}

// 授权后回调页面
func Page2Handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)

	cookie, err := r.Cookie("sid")
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	session, err := sessionStorage.Get(cookie.Value)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	savedState := session.(string)

	urlValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	code := urlValues.Get("code")
	urlState := urlValues.Get("state")

	if savedState != urlState {
		io.WriteString(w, fmt.Sprintf("state 不匹配, session 中的为 %q, url 传递过来的是 %q", savedState, urlState))
		return
	}

	if code == "" {
		log.Println("客户禁止授权")
		return
	}

	oauth2Client := oauth2.Client{
		Config: oauth2Config,
	}

	token, err := oauth2Client.Exchange(code)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}
	log.Printf("%+v\n", token)

	userinfo, err := oauth2Client.UserInfo(oauth2.Language_zh_CN)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(userinfo)
	log.Printf("%+v\n", userinfo)
	return
}

func main() {
	http.ListenAndServe(":80", nil)
}
```
