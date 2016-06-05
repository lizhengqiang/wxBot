package bot

import "net/http"

func NewBot(id string, cache Cache) (bot *WeixinBot) {

	bot = &WeixinBot{
		ID:     id,
		Cacher: cache,
		httpClient: &http.Client{
			Jar: NewCookieJar(cache),
		},
	}

	return
}
