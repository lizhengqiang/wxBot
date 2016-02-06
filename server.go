package main

import (
	"gopkg.in/macaron.v1"

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

type LogsResponse struct {
	Code int64
	Logs []string
}

type StatusResponse struct {
	Code int64
}

func main() {
	bots := make(map[string]*bot.WeixinBot)
	m := macaron.Classic()
	m.Use(macaron.Renderer())
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
	m.Run()
}

func startLogin(bot *bot.WeixinBot) {
	// 等待登陆
	bot.WaitForLogin()
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
