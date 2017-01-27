package api

import (
//"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"github.com/lizhengqiang/wxBot/bot"
	"github.com/cocotyty/summer"
	"github.com/lizhengqiang/wxBot/provider"
)


// 写个操作BOT的MiddleWare
func BotMiddleware(ctx *macaron.Context, sess session.Store, cache cache.Cache) {

	factory := summer.GetStoneWithName("CacherFactory").(*provider.CacherFactory)
	lastBot := bot.NewBot(sess.ID(), factory.NewCacher(sess.ID()))
	ctx.Map(lastBot)
	return

}
