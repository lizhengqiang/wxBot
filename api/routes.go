package api

import (
	"github.com/cocotyty/summer"
	"github.com/go-macaron/session"
	"github.com/lizhengqiang/wxBot/bot"
	"github.com/lizhengqiang/wxBot/domain"
	"github.com/lizhengqiang/wxBot/provider"
	"github.com/mougeli/beauty"
	"gopkg.in/macaron.v1"
)

func init() {
	summer.Put(&Router{})
}

type Router struct {
	M          *provider.HttpServer `sm:"*"`
	MQ         domain.MessageQueue  `sm:"*"`
	BotManager *bot.BotManager      `sm:"*"`
}

func (router *Router) Ready() {
	m := router.M.M
	m.Group("/api", func() {
		m.Get("/start", router.Start)
		m.Get("/qrcode", router.GetQrCodeUrl)
		m.Get("/status", func(ctx *macaron.Context, sess session.Store, b *bot.WeixinBot, r beauty.Render) {
			// 清理掉上次的
			r.OK(map[string]interface{}{
				"isRunning":     b.IsRunning(),
				"isLoopRunning": b.IsLoopRunning,
				"task":          b.Get("task"),
				"status":        b.Get("status"),
				"me":            b.GetMe(),
			})
		})
	}, router.BotMiddleware)
}
