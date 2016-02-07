package main

import (
	"gopkg.in/macaron.v1"

	"github.com/go-macaron/session"
	"strconv"
	"strings"
	"time"
	"wxBot/bot"
)

type QRCodeResponse struct {
	Code      int64
	QRCodeUrl string
}

type RegisterResponse struct {
	Code       int64
	HookMethod string
	HookUrl    string
}

type LogsResponse struct {
	Code int64
	Logs []string
}

type StatusResponse struct {
	Code int64
}

var bots map[string]*bot.WeixinBot

func main() {
	bots = make(map[string]*bot.WeixinBot)
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())

	m.Get("/:sessionId/start", func(ctx *macaron.Context) {
		lastBot, ok := bots[ctx.Params(":sessionId")]
		if ok {
			lastBot.Stop()
		}
		bot := &bot.WeixinBot{}
		bots[ctx.Params(":sessionId")] = bot
		bot.Init()
		go bot.Start()
		qrcodeUrl := bot.GetQrcodeUrl()

		ctx.JSON(200, &QRCodeResponse{
			QRCodeUrl: qrcodeUrl,
		})
	})
	m.Get("/:sessionId/register", func(ctx *macaron.Context) {
		bot, ok := bots[ctx.Params(":sessionId")]
		response := RegisterResponse{}
		if ok {
			response.Code = 0
			response.HookMethod = ctx.Query("HookMethod")
			response.HookUrl = ctx.Query("HookUrl")
			bot.RegisterHookUrl(response.HookMethod, response.HookUrl)
		} else {
			response.Code = 400
		}
		ctx.JSON(200, &response)
	})
	m.Get("/:sessionId/logs", func(ctx *macaron.Context) {
		bot, ok := bots[ctx.Params(":sessionId")]
		response := LogsResponse{}
		if ok {
			response.Code = 0
			response.Logs = bot.Logs
		} else {
			response.Code = 400
		}
		ctx.JSON(200, &response)
	})
	m.Get("/:sessionId/status", func(ctx *macaron.Context) {
		_, ok := bots[ctx.Params(":sessionId")]
		response := StatusResponse{}
		if ok {
			response.Code = 0
		} else {
			response.Code = 400
		}
		ctx.JSON(200, &response)
	})
	m.Get("/", func(ctx *macaron.Context, sess session.Store) {
		sessionId := ""
		sessionIdInterface := sess.Get("sessionId")
		if sessionIdInterface != nil {
			sessionId = sessionIdInterface.(string)
		}
		if sessionId == "" {
			sessionId = strconv.FormatInt(time.Now().Unix(), 10)
			sess.Set("sessionId", sessionId)
		}
		u := "/" + sessionId + "/qrcode"

		ctx.Redirect(u, 302)
	})
	m.Get("/:sessionId/qrcode", func(ctx *macaron.Context) {
		bot, ok := bots[ctx.Params(":sessionId")]
		response := QRCodeResponse{}
		QRCodeUrl := ""
		if ok {
			response.Code = 0
			QRCodeUrl = bot.GetQrcodeUrl()
		} else {
			response.Code = 400
		}

		ctx.Render.HTML(200, "pages/qrcode", map[string]interface{}{"QRCodeUrl": QRCodeUrl})
	})
	m.Get("/:sessionId/log", log)
	m.Get("/:sessionId/hook", hook)
	m.Get("/:sessionId/setHook", setHook)
	m.Get("/:sessionId/tool", tool)
	m.Get("/:sessionId/useTool", useTool)
	m.Get("/:sessionId/groups", groups)
	m.Get("/:sessionId/contacts", contacts)
	m.Get("/:sessionId/members", members)
	m.Get("/:sessionId/sendMsg", sendMsg)
	m.Run()
}

type Hook struct {
	Method string
	Url    string
	Title  string
}

//func hook(ctx *macaron.Context) {
//	hook := Hook{hook := Hook{
//	Mehotd: "contactMessage",
//	Url:    "http://wechat.lizhengqiang.alpha.mouge.cc/BlackHole/messageHook",
//	Title:  "图灵机器人",
//	}}
//	ctx.Render.HTML(200, "pages/log", []Hook{	})
//}

func hook(ctx *macaron.Context) {
	interfaces := map[string]interface{}{"Hooks": []Hook{
		Hook{
			Method: "contactMessage",
			Url:    "http://wechat.lizhengqiang.alpha.mouge.cc/BlackHole/messageHook",
			Title:  "图灵机器人",
		},
	}}
	ctx.Render.HTML(200, "pages/hook", interfaces)
}

func tool(ctx *macaron.Context) {
	interfaces := map[string]interface{}{"Tools": []Hook{
		Hook{
			Method: "singleHappyNewYear",
			Url:    "http://wechat.lizhengqiang.alpha.mouge.cc/BlackHole/messageHook",
			Title:  "随机单发新年快乐",
		},
	}}
	ctx.Render.HTML(200, "pages/tool", interfaces)
}

type UserToolResponse struct {
	Code int64
}

func useTool(ctx *macaron.Context) {
	bot, ok := bots[ctx.Params(":sessionId")]
	response := &SetHookResponse{}
	method := ctx.Query("Method")
	if ok {
		response.Code = 0

		if method == "singleHappyNewYear" {
			go func() {
				for _, contact := range bot.MemberList {
					if strings.Contains(contact.UserName, "@@") {

					} else if contact.VerifyFlag == 0 {
						bot.SendMsg("猴年猴赛雷,新年快乐!", contact.UserName)
						time.Sleep(1 * time.Second)
					}

				}
			}()

		}
	} else {
		response.Code = 400
	}

	ctx.JSON(200, response)
}

type SetHookResponse struct {
	Code  int64
	Hooks map[string]string
}

func setHook(ctx *macaron.Context) {
	bot, ok := bots[ctx.Params(":sessionId")]
	response := &SetHookResponse{}
	if ok {
		response.Code = 0
		bot.RegisterHookUrl(ctx.Query("Method"), ctx.Query("Url"))
		response.Hooks = bot.Hooks
	} else {
		response.Code = 400
	}

	ctx.JSON(200, response)

}

func log(ctx *macaron.Context) {
	bot, ok := bots[ctx.Params(":sessionId")]
	var Logs []string
	if ok {
		Logs = bot.Logs
	} else {
		Logs = make([]string, 0)
	}
	ctx.Render.HTML(200, "pages/log", map[string]interface{}{"Logs": Logs})
}

type SendMsgResponse struct {
	Code int64
}

func sendMsg(ctx *macaron.Context) {
	bot, ok := bots[ctx.Params(":sessionId")]
	response := SendMsgResponse{}
	content := ctx.Query("Content")
	userName := ctx.Query("UserName")
	if ok {
		response.Code = 0
		bot.SendMsg(content, userName)
	} else {
		response.Code = 400
	}
	ctx.JSON(200, &response)
}

type ContactResponse struct {
	Code int64
	List []bot.Contact
}

func members(ctx *macaron.Context) {
	bot, ok := bots[ctx.Params(":sessionId")]
	response := ContactResponse{}
	if ok {
		response.Code = 0
		response.List = bot.MemberList
	} else {
		response.Code = 400
	}
	ctx.JSON(200, &response)
}

func contacts(ctx *macaron.Context) {
	bot, ok := bots[ctx.Params(":sessionId")]
	response := ContactResponse{}
	if ok {
		response.Code = 0
		response.List = bot.ContactList
	} else {
		response.Code = 400
	}
	ctx.JSON(200, &response)
}
func groups(ctx *macaron.Context) {
	bot, ok := bots[ctx.Params(":sessionId")]
	response := ContactResponse{}
	if ok {
		response.Code = 0
		response.List = bot.GroupList
	} else {
		response.Code = 400
	}
	ctx.JSON(200, &response)
}
