package bot

import (
	"github.com/cocotyty/cookiejar"
	"net/http"
	"net/url"
	"github.com/cocotyty/json"
)

type CacherCookieJar struct {
	Cacher Cache
}

func NewCookieJar(cacher Cache) http.CookieJar {
	return &CacherCookieJar{
		Cacher: cacher,
	}

}

func (this *CacherCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	var jar *cookiejar.Jar

	bytes := []byte(this.Cacher.Get("cookieJar"))
	err := json.Unmarshal(bytes, jar)
	if err != nil {
		jar, _ = cookiejar.New(nil)
	}
	jar.SetCookies(u, cookies)
	bytes, err = json.Marshal(jar)
	if err != nil {
		return
	}
	this.Cacher.Set("cookieJar", string(bytes))
}
func (this *CacherCookieJar) Cookies(u *url.URL) (cookies []*http.Cookie) {

	bytes := []byte(this.Cacher.Get("cookieJar"))
	jar, err := cookiejar.LoadFromJson(nil, bytes)
	if err != nil {
		return
	}
	return jar.Cookies(u)
}
