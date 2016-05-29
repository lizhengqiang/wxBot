package main

import (
	"github.com/go-macaron/cache"
	_ "github.com/go-macaron/cache/memcache"
	"github.com/go-macaron/session"
	_ "github.com/go-macaron/session/memcache"
	"github.com/mougeli/beauty"
	"gopkg.in/macaron.v1"
	"os"
	"wxBot/api"
)

func main() {

	memcachedAddr := os.Getenv("memcachedAddr")
	opt := session.Options{
		Provider:       "memcache",
		ProviderConfig: memcachedAddr,
	}

	cacheOpt := cache.Options{
		Adapter:       "memcache",
		AdapterConfig: memcachedAddr,
	}
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner(opt))
	m.Use(cache.Cacher(cacheOpt))

	// ! 另一个开源项目,API使用的,可以美化一下输出, github.com/mougeli/beauty
	m.Use(beauty.Renderer())

	api.RegisterRoutes(m)

	m.Run()
}
