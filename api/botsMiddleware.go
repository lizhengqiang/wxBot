package api

import (
	"github.com/go-macaron/cache"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

// 写个操作BOT的MiddleWare
func (router *Router) BotMiddleware(ctx *macaron.Context, sess session.Store, cache cache.Cache) {


	lastBot :=router.BotManager.Get(sess.ID())
	ctx.Map(lastBot)
	return

}
