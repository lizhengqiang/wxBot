package bot

import (
	"strings"
	"time"
	"regexp"
	"errors"
	"io/ioutil"
	"fmt"
)

func init() {

}

var (
	ErrLoginReq     error = errors.New("登录请求出错")
	ErrLoginFailed  error = errors.New("登录失败")
	ErrLoginCancels error = errors.New("登陆被取消")
	ErrDoLogin      error = errors.New("机器人登录时出错")
	WaitScanQRCode        = errors.New("等待扫描")
)

func (bot *WeixinBot) checkQRCodeScanStatus(tip, uuid string) (all [][]byte, body []byte, err error) {
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

func (bot *WeixinBot) HandleQRCodeScanStatus() (err error) {
	bot.Println("等待QRCode扫描")
	// 获取登陆返回值
	all, body, err := bot.checkQRCodeScanStatus(bot.getProperty(tip), bot.getProperty(UUID))

	if err != nil {
		bot.log("登录时出错")
		err = ErrLoginReq
		return
	}

	if len(all) >= 2 {
		code := string(all[1])

		if code == "201" {
			bot.setProperty(tip, "0")
			bot.Println("* 成功扫描,请在手机上点击确认以登录.")
			err = WaitScanQRCode
			return
		}

		if code == "200" {
			reRedirectUri, _ := regexp.Compile(`window.redirect_uri="(\S+?)";`)
			allRedirectUri := reRedirectUri.FindSubmatch(body)
			if len(allRedirectUri) >= 2 {
				redirectUri := string(allRedirectUri[1]) + "&fun=new"
				bot.setProperty("redirectUri", redirectUri)
				bot.setProperty("baseUri", string([]byte(redirectUri)[0:strings.LastIndex(redirectUri, "/")]))
			}
			bot.log("* 登陆成功.")
			return
		}

		if code == "408" {
			bot.log("! 登陆超时.")
			err = WaitScanQRCode
			return
		}

		bot.log("! 登录失败 %s", code)
		err = ErrLoginFailed
		return

	}

	err = ErrLoginReq
	return
}

// 获取二维码地址
func (bot *WeixinBot) GetQrcodeUrl() string {
	return "https://login.weixin.qq.com/qrcode/" + bot.getProperty(UUID)
}

//等待登陆
func (bot *WeixinBot) waitForScanQrcode() (err error) {
	for bot.IsRunning() {
		bot.HandleQRCodeScanStatus()
		time.Sleep(time.Second * 3)
	}
	bot.Println(bot.IsRunning())
	return ErrLoginCancels

}

func (bot *WeixinBot) scanForLogin() (err error) {

	// 等待登陆
	err = bot.waitForScanQrcode()
	if err != nil {
		bot.Println(err)
		return
	}

	// 登陆
	err = bot.Login()
	if err != nil {
		bot.Println("登陆失败. ")
		return
	}

	// 初始化信息
	bot.InitWebWeixin()

	bot.Println("初始化信息完毕")

	my := &User{}
	bot.unmarshal(me, my)

	bot.WebWeixinStatusNotify(my.UserName, my.UserName, 3)

	bot.Println("初始化信息完毕")

	return nil
}
