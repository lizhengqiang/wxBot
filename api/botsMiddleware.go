package api

import (
	//"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"wxBot/bot"
)

type Cacher struct {
	ID    string
	cache cache.Cache
}

func NewCacher(ID string, cache cache.Cache) *Cacher {
	return &Cacher{
		ID:    ID,
		cache: cache,
	}
}
func (this *Cacher) Get(key string) (value string) {
	value, _ = this.cache.Get(this.ID + ":" + key).(string)
	return
}

func (this *Cacher) Set(key string, value string) (err error) {
	return this.cache.Put(this.ID+":"+key, value, 0)
}

// 写个操作BOT的MiddleWare
func BotMiddleware(ctx *macaron.Context, sess session.Store, cache cache.Cache) {

	lastBot := bot.NewBot(NewCacher(sess.ID(), cache))

	ctx.Map(lastBot)
	return

}
