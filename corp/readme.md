# 微信企业号 golang SDK

## 整体架构图

![架构图](https://github.com/chanxuehong/wechat/blob/master/corp/corp.png)

## 示例

### 主動調用微信 api，corp 子包裏面的 Client 基本都是這樣的調用方式
```Go
package main

import (
	"fmt"

	"github.com/chanxuehong/wechat/corp"
	"github.com/chanxuehong/wechat/corp/menu"
)

var AccessTokenServer = corp.NewDefaultAccessTokenServer("corpId", "corpSecret", nil) // 一個應用只能有一個實例
var corpClient = corp.NewClient(AccessTokenServer, nil)

func main() {
	var subButtons = make([]menu.Button, 2)
	subButtons[0].SetAsViewButton("搜索", "http://www.soso.com/")
	subButtons[1].SetAsClickButton("赞一下我们", "V1001_GOOD")

	var mn menu.Menu
	mn.Buttons = make([]menu.Button, 3)
	mn.Buttons[0].SetAsClickButton("今日歌曲", "V1001_TODAY_MUSIC")
	mn.Buttons[1].SetAsViewButton("视频", "http://v.qq.com/")
	mn.Buttons[2].SetAsSubMenuButton("子菜单", subButtons)

	menuClient := (*menu.Client)(corpClient)
	if err := menuClient.CreateMenu(0 /* agentId */, mn); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ok")
}
```

### 被動接收消息（事件）推送，一个 URL 监听一个企业号应用的消息
```Go
package main

import (
	"log"
	"net/http"

	"github.com/chanxuehong/wechat/corp"
	"github.com/chanxuehong/wechat/corp/message/request"
	"github.com/chanxuehong/wechat/corp/message/response"
	"github.com/chanxuehong/wechat/util"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
}

// 文本消息的 Handler
func TextMessageHandler(w http.ResponseWriter, r *corp.Request) {
	// 简单起见，把用户发送过来的文本原样回复过去
	text := request.GetText(r.MixedMsg) // 可以省略, 直接从 r.MixedMsg 取值
	resp := response.NewText(text.FromUserName, text.ToUserName, text.CreateTime, text.Content)
	corp.WriteResponse(w, r, resp)
}

func main() {
	aesKey, err := util.AESKeyDecode("encodedAESKey") // 这里 encodedAESKey 改成你自己的参数
	if err != nil {
		panic(err)
	}

	messageServeMux := corp.NewMessageServeMux()
	messageServeMux.MessageHandleFunc(request.MsgTypeText, TextMessageHandler) // 注册文本处理 Handler

	// 下面函数的几个参数设置成你自己的参数: corpId, agentId, token
	agentServer := corp.NewDefaultAgentServer("corpId", 0 /* agentId */, "token", aesKey, messageServeMux)

	agentServerFrontend := corp.NewAgentServerFrontend(agentServer, corp.ErrorHandlerFunc(ErrorHandler), nil)

	// 如果你在微信后台设置的回调地址是
	//   http://xxx.yyy.zzz/agent
	// 那么可以这么注册 http.Handler
	http.Handle("/agent", agentServerFrontend)
	http.ListenAndServe(":80", nil)
}
```

### 被動接收消息（事件）推送，一个 URL 监听多个企业号应用的消息
```Go
package main

import (
	"log"
	"net/http"

	"github.com/chanxuehong/wechat/corp"
	"github.com/chanxuehong/wechat/corp/message/request"
	"github.com/chanxuehong/wechat/corp/message/response"
	"github.com/chanxuehong/wechat/util"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
}

// 文本消息的 Handler
func TextMessageHandler(w http.ResponseWriter, r *corp.Request) {
	// 简单起见，把用户发送过来的文本原样回复过去
	text := request.GetText(r.MixedMsg) // 可以省略, 直接从 r.MixedMsg 取值
	resp := response.NewText(text.FromUserName, text.ToUserName, text.CreateTime, text.Content)
	corp.WriteResponse(w, r, resp)
}

func main() {
	// agentServer1
	aesKey1, err := util.AESKeyDecode("encodedAESKey1") // 这里 encodedAESKey1 改成你自己的参数
	if err != nil {
		panic(err)
	}

	messageServeMux1 := corp.NewMessageServeMux()
	messageServeMux1.MessageHandleFunc(request.MsgTypeText, TextMessageHandler) // 注册文本处理 Handler

	// 下面函数的几个参数设置成你自己的参数: corpId1, agentId1, token1
	agentServer1 := corp.NewDefaultAgentServer("corpId1", 1 /* agentId1 */, "token1", aesKey1, messageServeMux1)

	// agentServer2
	aesKey2, err := util.AESKeyDecode("encodedAESKey2") // 这里 encodedAESKey2 改成你自己的参数
	if err != nil {
		panic(err)
	}

	messageServeMux2 := corp.NewMessageServeMux()
	messageServeMux2.MessageHandleFunc(request.MsgTypeText, TextMessageHandler) // 注册文本处理 Handler

	// 下面函数的几个参数设置成你自己的参数: corpId2, agentId2, token2
	agentServer2 := corp.NewDefaultAgentServer("corpId2", 2 /* agentId2 */, "token2", aesKey2, messageServeMux2)

	// multiAgentServerFrontend, 第一个参数可以不是 agent_server, 可以自己指定
	multiAgentServerFrontend := corp.NewMultiAgentServerFrontend("agent_server", corp.ErrorHandlerFunc(ErrorHandler), nil)
	multiAgentServerFrontend.SetAgentServer("agent1", agentServer1) // 回調url上面要加上 agent_server=agent1
	multiAgentServerFrontend.SetAgentServer("agent2", agentServer2) // 回調url上面要加上 agent_server=agent2

	// 如果你在微信后台设置的回调地址是
	//   http://xxx.yyy.zzz/agent
	// 那么可以这么注册 http.Handler
	http.Handle("/agent", multiAgentServerFrontend)
	http.ListenAndServe(":80", nil)
}
```