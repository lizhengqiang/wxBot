package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"qiniupkg.com/x/errors.v7"
	"regexp"
	"time"
)

// 初始化这个会话
func (this *WeixinBot) Init() {
	this.set(isRunning, TRUE)

	this.setProperty(deviceId, "e"+string([]byte(fmt.Sprint(rand.Float64()))[2:17]))

	resp, err := this.httpClient.PostForm("https://login.weixin.qq.com/jslogin", url.Values{"appid": {"wx782c26e4c19acffb"}, "fun": {"new"}, "lang": {"zh_CN"}, "_": {this.timestamp()}})
	if err != nil {
		this.log(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	re, _ := regexp.Compile(`window.QRLogin.code = (\d+); window.QRLogin.uuid = "(\S+?)"`)
	all := re.FindSubmatch(body)
	if len(all) >= 3 {
		code := all[1]
		uuid := all[2]
		this.Println(code, uuid)
		if string(code) == "200" {
			this.setProperty(UUID, string(uuid))
		} else {
			this.log("! 初始化失败. %s", string(code))
		}
	}

	this.log("* 初始化成功.")
}

func (this *WeixinBot) scanForLogin() (err error) {

	// 等待登陆
	err = this.waitForScanQrcode()
	if err != nil {
		this.Println(err)
		return
	}

	// 登陆
	err = this.Login()
	if err != nil {
		this.Println("登陆失败. ")
		return
	}

	// 初始化信息
	this.InitWebWeixin()

	this.Println("初始化信息完毕")

	this.WebWeixinStatusNotify()

	this.Println("初始化信息完毕")

	return nil
}

type InitWebWeixinRequestBody struct {
	BaseRequest *BaseRequest
}

type InitWebWeixinResponseBody struct {
	BaseResponse *BaseResponse
	SyncKey      *SyncKey
	User         *User
}

func (this *WeixinBot) InitWebWeixin() int64 {
	requestBody := InitWebWeixinRequestBody{
		BaseRequest: this.getBaseRequest(),
	}
	u := fmt.Sprintf("/webwxinit?pass_ticket=%s&skey=%s&r=%s", this.getProperty(passTicket), this.getProperty(skey), this.timestamp())

	respJson := &InitWebWeixinResponseBody{}

	this.PostJson(u, requestBody, respJson)

	this.marshal(me, respJson.User)

	this.saveSyncKey(respJson.SyncKey)

	return respJson.BaseResponse.Ret

}

type WebWeixinStatusNotifyRequest struct {
	BaseRequest  *BaseRequest
	Code         int64
	FromUserName string
	ToUserName   string
	ClientMsgId  int64
}

type WebWeixinStatusNotifyResponseBody struct {
	BaseResponse *BaseResponse
}

var ErrStatusNotify error = errors.New("检查状态出错")

func (bot *WeixinBot) WebWeixinStatusNotify() (err error) {
	my := &User{}
	bot.unmarshal(me, my)
	requestBody := WebWeixinStatusNotifyRequest{
		BaseRequest:  bot.getBaseRequest(),
		Code:         int64(3),
		FromUserName: my.UserName,
		ToUserName:   my.UserName,
		ClientMsgId:  time.Now().Unix(),
	}

	respBody, _ := bot.SimplePostJson(fmt.Sprintf("/webwxstatusnotify?lang=zh_CN&pass_ticket=%s", bot.getProperty(passTicket)), requestBody)
	respJson := WebWeixinStatusNotifyResponseBody{}

	err = json.Unmarshal(respBody, &respJson)

	if err != nil {
		bot.log(err.Error())
		return
	}

	if respJson.BaseResponse.Ret == int64(0) {
		return nil
	}
	return ErrStatusNotify
}
