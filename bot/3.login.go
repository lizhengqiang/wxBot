package bot

import (
	"golang.org/x/net/html"
	"strings"
)

func init() {

}

type LoginHtml struct {
	Html struct {
		Head struct {
		} `xml:"head"`
		Body struct {
			Error struct {
				Ret         string `xml:"ret"`
				Message     string `xml:"message"`
				Skey        string `xml:"skey"`
				Wxsid       string `xml:"wxsid"`
				Wxuin       string `xml:"wxuin"`
				PassTicket  string `xml:"pass_ticket"`
				IsGrayscale string `xml:"isgrayscale"`
			} `xml:"error"`
		} `xml:"body"`
	} `xml:"html"`
}

// 登陆
func (bot *WeixinBot) Login() error {
	resp, err := bot.httpClient.Get(bot.getProperty("redirectUri"))

	if err != nil {
		bot.Println(err.Error())
		return ErrDoLogin
	}
	if resp.Body == nil {
		return ErrDoLogin
	}
	defer resp.Body.Close()

	// 卧槽 这地下是要解析HTML呀
	doc, err := html.Parse(resp.Body)
	if err != nil {
		bot.Println(err.Error())
		return ErrDoLogin
	}

	var f func(*html.Node)

	f = func(n *html.Node) {
		name := strings.TrimSpace(n.Data)
		data := ""
		if n.FirstChild != nil {
			data = strings.TrimSpace(n.FirstChild.Data)
		}

		if name == "skey" {
			bot.setProperty(skey, data)
			//bot.SKey = data
		} else if name == "wxsid" {
			bot.setProperty(wxsid, data)
			//bot.WxSid = data
		} else if name == "wxuin" {
			//wxUin, _ := strconv.ParseInt(data, 10, 64)
			bot.Set(wxuni, data)
			//bot.WxUin = wxUin
		} else if name == "pass_ticket" {
			bot.setProperty(passTicket, data)
			//bot.PassTicket = data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return nil
}
