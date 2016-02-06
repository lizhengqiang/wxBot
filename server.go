package main

import (
	"fmt"
	"gopkg.in/macaron.v1"
	"time"
	"wxBot/bot"
)

type QRCodeResponse struct {
	QRCodeUrl string
}

type RegisterResponse struct {
	Code       int64
	HookMethod string
	HookUrl    string
}

func main() {
	bots := make(map[string]*bot.WeixinBot)
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Get("/:sessionId/start", func(ctx *macaron.Context) {
		bot := &bot.WeixinBot{}
		bots[ctx.Params(":sessionId")] = bot
		bot.Start()
		qrcodeUrl := bot.GetQrcodeUrl()
		go startLogin(bot)

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
	m.Run()
}

func startLogin(bot *bot.WeixinBot) {
	// 等待登陆
	for code := bot.WaitForLogin(); code != "200"; code = bot.WaitForLogin() {
		fmt.Println(code)
		time.Sleep(time.Second * 3)
	}
	// 登陆
	bot.Login()
	// 初始化信息
	bot.InitBaseRequest()
	bot.InitWebWeixin()
	// 获取联系人列表
	bot.WebWeixinStatusNotify()
	bot.GetContact()
	// 开始监听消息
	bot.ListenMsgMode()
}
