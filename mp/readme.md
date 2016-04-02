# 微信公众平台（订阅号、服务号） golang SDK

## 整体架构图

![架构图](https://github.com/chanxuehong/wechat/blob/master/mp/mp.png)

## 示例

### 主動調用微信 api，mp 子包裏面的 Client 基本都是這樣的調用方式
```Go
package main

import (
	"fmt"

	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/menu"
)

var AccessTokenServer = mp.NewDefaultAccessTokenServer("appId", "appSecret", nil) // 一個應用只能有一個實例
var mpClient = mp.NewClient(AccessTokenServer, nil)

func main() {
	var subButtons = make([]menu.Button, 2)
	subButtons[0].SetAsViewButton("搜索", "http://www.soso.com/")
	subButtons[1].SetAsClickButton("赞一下我们", "V1001_GOOD")

	var mn menu.Menu
	mn.Buttons = make([]menu.Button, 3)
	mn.Buttons[0].SetAsClickButton("今日歌曲", "V1001_TODAY_MUSIC")
	mn.Buttons[1].SetAsViewButton("视频", "http://v.qq.com/")
	mn.Buttons[2].SetAsSubMenuButton("子菜单", subButtons)

	menuClient := (*menu.Client)(mpClient)
	if err := menuClient.CreateMenu(mn); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ok")
}
```

### 被動接收消息（事件）推送，一个 URL 监听一个公众号的消息
```Go
package main

import (
	"log"
	"net/http"

	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/message/request"
	"github.com/chanxuehong/wechat/mp/message/response"
	"github.com/chanxuehong/wechat/util"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
}

// 文本消息的 Handler
func TextMessageHandler(w http.ResponseWriter, r *mp.Request) {
	// 简单起见，把用户发送过来的文本原样回复过去
	text := request.GetText(r.MixedMsg) // 可以省略, 直接从 r.MixedMsg 取值
	resp := response.NewText(text.FromUserName, text.ToUserName, text.CreateTime, text.Content)
	//mp.WriteRawResponse(w, r, resp) // 明文模式
	mp.WriteAESResponse(w, r, resp) // 安全模式
}

func main() {
	aesKey, err := util.AESKeyDecode("encodedAESKey") // 这里 encodedAESKey 改成你自己的参数
	if err != nil {
		panic(err)
	}

	messageServeMux := mp.NewMessageServeMux()
	messageServeMux.MessageHandleFunc(request.MsgTypeText, TextMessageHandler) // 注册文本处理 Handler

	// 下面函数的几个参数设置成你自己的参数: oriId, token, appId
	mpServer := mp.NewDefaultServer("oriId", "token", "appId", aesKey, messageServeMux)

	mpServerFrontend := mp.NewServerFrontend(mpServer, mp.ErrorHandlerFunc(ErrorHandler), nil)

	// 如果你在微信后台设置的回调地址是
	//   http://xxx.yyy.zzz/wechat
	// 那么可以这么注册 http.Handler
	http.Handle("/wechat", mpServerFrontend)
	http.ListenAndServe(":80", nil)
}
```

### 被動接收消息（事件）推送，一个 URL 监听多个公众号的消息
```Go
package main

import (
	"log"
	"net/http"

	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/message/request"
	"github.com/chanxuehong/wechat/mp/message/response"
	"github.com/chanxuehong/wechat/util"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
}

// 文本消息的 Handler
func TextMessageHandler(w http.ResponseWriter, r *mp.Request) {
	// 简单起见，把用户发送过来的文本原样回复过去
	text := request.GetText(r.MixedMsg) // 可以省略, 直接从 r.MixedMsg 取值
	resp := response.NewText(text.FromUserName, text.ToUserName, text.CreateTime, text.Content)
	//mp.WriteRawResponse(w, r, resp) // 明文模式
	mp.WriteAESResponse(w, r, resp) // 安全模式
}

func main() {
	// mpServer1
	aesKey1, err := util.AESKeyDecode("encodedAESKey1") // 这里 encodedAESKey1 改成你自己的参数
	if err != nil {
		panic(err)
	}

	messageServeMux1 := mp.NewMessageServeMux()
	messageServeMux1.MessageHandleFunc(request.MsgTypeText, TextMessageHandler) // 注册文本处理 Handler

	// 下面函数的几个参数设置成你自己的参数: oriId1, token1, appId1
	mpServer1 := mp.NewDefaultServer("oriId1", "token1", "appId1", aesKey1, messageServeMux1)

	// mpServer2
	aesKey2, err := util.AESKeyDecode("encodedAESKey2") // 这里 encodedAESKey2 改成你自己的参数
	if err != nil {
		panic(err)
	}

	messageServeMux2 := mp.NewMessageServeMux()
	messageServeMux2.MessageHandleFunc(request.MsgTypeText, TextMessageHandler) // 注册文本处理 Handler

	// 下面函数的几个参数设置成你自己的参数: oriId2, token2, appId2
	mpServer2 := mp.NewDefaultServer("oriId2", "token2", "appId2", aesKey2, messageServeMux2)

	// multiServerFrontend, 第一个参数可以不是 wechat_server, 可以自己指定
	multiServerFrontend := mp.NewMultiServerFrontend("wechat_server", mp.ErrorHandlerFunc(ErrorHandler), nil)
	multiServerFrontend.SetServer("wechat1", mpServer1) // 回調url上面要加上 wechat_server=wechat1
	multiServerFrontend.SetServer("wechat2", mpServer2) // 回調url上面要加上 wechat_server=wechat2

	// 如果你在微信后台设置的回调地址是
	//   http://xxx.yyy.zzz/wechat
	// 那么可以这么注册 http.Handler
	http.Handle("/wechat", multiServerFrontend)
	http.ListenAndServe(":80", nil)
}
```