package bot

import (
	"github.com/lizhengqiang/wxBot/domain"
	"net/http"
	"regexp"
	"github.com/robertkrimen/otto"
)

func NewBot(id string, cache Cache, mq domain.MessageQueue) (bot *WeixinBot) {

	bot = &WeixinBot{
		ID:     id,
		Cacher: cache,
		httpClient: &http.Client{
			Jar: NewCookieJar(cache),
		},
		MQ: mq,
		hears:map[*regexp.Regexp]otto.Value{},
	}
	bot.ReloadJS()
	return
}
