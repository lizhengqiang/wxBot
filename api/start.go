package api

import (
	"github.com/go-macaron/session"
	"github.com/lizhengqiang/wxBot/bot"
	"github.com/lizhengqiang/wxBot/domain"
	"github.com/mougeli/beauty"
	"gopkg.in/macaron.v1"
)

func (router *Router) Start(ctx *macaron.Context, sess session.Store, b *bot.WeixinBot, r beauty.Render) {
	// 清理掉上次的
	if b.IsRunning() {
		if !b.IsLoopRunning {
			router.MQ.Send(&domain.Message{BotID: sess.ID(), Type: "handleMsg", Body: nil})
		}

	} else {
		router.MQ.Send(&domain.Message{BotID: sess.ID(), Type: "start", Body: nil})
	}

	r.OK(nil)
}
