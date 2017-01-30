package bot

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"regexp"
)

// 初始化这个会话
func (bot *WeixinBot) Init() {
	bot.Set(IsRunning, TRUE)
	bot.Set(IsSigning, TRUE)

	bot.setProperty(deviceId, "e"+string([]byte(fmt.Sprint(rand.Float64()))[2:17]))

	resp, err := bot.httpClient.PostForm("https://login.weixin.qq.com/jslogin", url.Values{"appid": {"wx782c26e4c19acffb"}, "fun": {"new"}, "lang": {"zh_CN"}, "_": {bot.timestamp()}})
	if err != nil {
		bot.log(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	re, _ := regexp.Compile(`window.QRLogin.code = (\d+); window.QRLogin.uuid = "(\S+?)"`)
	all := re.FindSubmatch(body)
	if len(all) >= 3 {
		code := all[1]
		uuid := all[2]
		bot.Println(code, uuid)
		if string(code) == "200" {
			bot.setProperty(UUID, string(uuid))
		} else {
			bot.log("! 初始化失败. %s", string(code))
		}
	}

	bot.log("* 初始化成功.")
}
