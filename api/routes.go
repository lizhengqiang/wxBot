package api

import (
	"github.com/go-macaron/session"
	"github.com/mougeli/beauty"
	"gopkg.in/macaron.v1"
	"wxBot/bot"
)

func RegisterRoutes(m *macaron.Macaron) {
	m.Group("/api", func() {
		m.Get("/start", Start)
		m.Get("/qrcode", GetQrCodeUrl)
		m.Get("/status", func(ctx *macaron.Context, sess session.Store, b *bot.WeixinBot, r beauty.Render) {
			// 清理掉上次的
			r.OK(map[string]interface{}{
				"isRunning": b.IsRunning(),
				"task":      b.Get("task"),
				"status":    b.Get("status"),
				"me":        b.GetMe(),
			})
		})
	}, BotMiddleware)
}
