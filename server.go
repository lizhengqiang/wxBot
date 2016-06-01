package main

import (
	"github.com/cocotyty/summer"
	"github.com/go-macaron/cache"
	_ "github.com/go-macaron/cache/redis"
	"github.com/go-macaron/session"
	_ "github.com/go-macaron/session/redis"
	"github.com/mougeli/beauty"
	"gopkg.in/macaron.v1"
	"os"
	"wxBot/api"
	"wxBot/mns"
"github.com/gogap/ali_mns"
)

func main() {

	summer.TomlFile(os.Getenv("configPath"))
	summer.Start()

	mq := summer.GetStoneWithName("MQ").(*mns.MQ)
	respChan := make(chan ali_mns.MessageReceiveResponse)
	errChan := make(chan error)

	for {
		select {
		case resp := <-respChan:
			{

			}
		case err := <-errChan:
			{
				if ali_mns.ERR_MNS_QUEUE_NOT_EXIST.IsEqual(err) {
				} else if ali_mns.ERR_MNS_MESSAGE_NOT_EXIST.IsEqual(err) {
				} else {
				}
			}
		}
	}

	redisAddr := os.Getenv("redisAddr")
	opt := session.Options{
		Provider:       "redis",
		ProviderConfig: redisAddr,
	}

	cacheOpt := cache.Options{
		Adapter:       "redis",
		AdapterConfig: redisAddr,
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
