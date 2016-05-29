package bot

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"qiniupkg.com/x/errors.v7"
	"regexp"
	"strings"
	"time"
)

var (
	ErrLoginReq error = errors.New("登录请求出错")

	ErrLoginFailed error = errors.New("登录失败")

	ErrLoginCancels error = errors.New("登陆被取消")

	ErrDoLogin error = errors.New("机器人登录时出错")
)

func (bot *WeixinBot) checkLoginStatus(tip, uuid string) (all [][]byte, body []byte, err error) {
	u := fmt.Sprintf("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?tip=%s&uuid=%s&_=%s", tip, uuid, bot.timestamp())
	resp, err := bot.httpClient.Get(u)
	if err != nil {
		bot.log(err.Error())
		return
	}
	if resp.Body == nil {
		return nil, nil, errors.New("失败")
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	re, _ := regexp.Compile(`window.code=(\d+);`)
	all = re.FindSubmatch(body)
	return
}

// 获取二维码地址
func (bot *WeixinBot) GetQrcodeUrl() string {
	return "https://login.weixin.qq.com/qrcode/" + bot.getProperty(UUID)
}

// 等待登陆
func (this *WeixinBot) waitForScanQrcode() (err error) {
	for this.IsRunning() {
		this.Println("等待QRCode扫描")
		// 获取登陆返回值
		all, body, err := this.checkLoginStatus(this.getProperty(tip), this.getProperty(UUID))

		if err != nil {
			this.log("登录时出错")
			return ErrLoginReq
		}

		if len(all) >= 2 {
			code := string(all[1])

			if code == "201" {
				this.setProperty(tip, "0")
				this.Println("* 成功扫描,请在手机上点击确认以登录.")
				continue
			}

			if code == "200" {
				reRedirectUri, _ := regexp.Compile(`window.redirect_uri="(\S+?)";`)
				allRedirectUri := reRedirectUri.FindSubmatch(body)
				if len(allRedirectUri) >= 2 {
					redirectUri := string(allRedirectUri[1]) + "&fun=new"

					this.setProperty("redirectUri", redirectUri)
					this.setProperty("baseUri", string([]byte(redirectUri)[0:strings.LastIndex(redirectUri, "/")]))
				}
				this.log("* 登陆成功.")
				return nil
			}
			if code == "408" {
				this.log("! 登陆超时.")
				continue
			}

			this.log("! 登录失败 %s", code)
			return ErrLoginFailed

		}
		time.Sleep(time.Second * 3)
	}
	this.Println(this.IsRunning())
	return ErrLoginCancels

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
			bot.set(wxuni, data)
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
