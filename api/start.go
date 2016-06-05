package api

import (
	"github.com/cocotyty/summer"
	"github.com/go-macaron/session"
	"github.com/mougeli/beauty"
	"gopkg.in/macaron.v1"
	"wxBot/bot"
	"wxBot/mns"
)

func Start(ctx *macaron.Context, sess session.Store, b *bot.WeixinBot, r beauty.Render) {
	// 清理掉上次的
	summer.GetStoneWithName("Trigger").(*mns.Trigger).Send(sess.ID(), "start", nil)
	r.OK(nil)
}
