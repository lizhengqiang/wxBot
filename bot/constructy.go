package bot

import "net/http"

func NewBot(cache Cache) (bot *WeixinBot) {

	bot = &WeixinBot{
		Cacher: cache,
		httpClient: &http.Client{
			Jar: NewCookieJar(cache),
		},
	}

	return
}
