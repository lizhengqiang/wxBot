package api

import (
	"github.com/go-macaron/session"
	"github.com/mougeli/beauty"
	"gopkg.in/macaron.v1"
	"wxBot/bot"
)

func Start(ctx *macaron.Context, sess session.Store, b *bot.WeixinBot, r beauty.Render) {
	// 清理掉上次的
	if ctx.Query("restart") == "true" {
		b.Stop()
	}
	go b.Start()
	r.OK(nil)
}