package provider

import (
	"github.com/cocotyty/summer"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

import (
	_ "github.com/go-macaron/cache/redis"
	_ "github.com/go-macaron/session/redis"
	"github.com/mougeli/beauty"
	"os"
)

type HttpServer struct {
	M                *macaron.Macaron
	SessionRedisConf string
	CacheRedisConf   string
}

func (this *HttpServer) Init() {
	this.SessionRedisConf = os.Getenv("redis.session")
	this.CacheRedisConf = os.Getenv("redis.cache")
}
func (this *HttpServer) Ready() {
	opt := session.Options{
		Provider:       "redis",
		ProviderConfig: this.SessionRedisConf,
	}

	cacheOpt := cache.Options{
		Adapter:       "redis",
		AdapterConfig: this.CacheRedisConf,
	}

	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner(opt))
	m.Use(cache.Cacher(cacheOpt))

	// ! 另一个开源项目,API使用的,可以美化一下输出, github.com/mougeli/beauty
	m.Use(beauty.Renderer())

	this.M = m
}

func init() {
	summer.Add("HttpServer", &HttpServer{})
}
