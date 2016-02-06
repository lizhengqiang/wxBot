package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"golang.org/x/net/html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type WeixinBot struct {
	UUID          string
	Tip           string
	RedirectUri   string
	BaseUri       string
	SKey          string
	WxSid         string
	WxUin         int64
	PassTicket    string
	BaseRequest   *BaseRequest
	My            *User
	SyncKey       *SyncKey
	SyncKeyString string
}

func (bot *WeixinBot) timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)

}

func (bot *WeixinBot) fmt() {
	fmt.Println(bot.UUID)
}

func (bot *WeixinBot) Start() {

	resp, err := http.PostForm("https://login.weixin.qq.com/jslogin", url.Values{"appid": {"wx782c26e4c19acffb"}, "fun": {"new"}, "lang": {"zh_CN"}, "_": {bot.timestamp()}})
	if err != nil {
		panic("请求jslogin出现错误")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	re, _ := regexp.Compile(`window.QRLogin.code = (\d+); window.QRLogin.uuid = "(\S+?)"`)
	all := re.FindSubmatch(body)
	if len(all) >= 3 {
		code := all[1]
		uuid := all[2]
		fmt.Println(string(code), string(uuid), string(code) == "200")
		if string(code) == "200" {
			bot.UUID = string(uuid)
		} else {
			panic("请求jslogin返回!200")
		}
	}
}

func (bot *WeixinBot) GetQrcodeUrl() string {
	return "https://login.weixin.qq.com/qrcode/" + bot.UUID
}

func (bot *WeixinBot) WaitForLogin() string {
	resp, err := http.Get(fmt.Sprintf("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?tip=%s&uuid=%s&_=%s", bot.Tip, bot.UUID, bot.timestamp()))
	if err != nil {
		panic("请求login出现错误")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	re, _ := regexp.Compile(`window.code=(\d+);`)
	all := re.FindSubmatch(body)
	if len(all) >= 2 {
		code := string(all[1])
		fmt.Println(string(code))
		if code == "201" {
			bot.Tip = "0"
			fmt.Println("成功扫描, 请在手机上点击确认以登录")
		} else if code == "200" {
			reRedirectUri, _ := regexp.Compile(`window.redirect_uri="(\S+?)";`)
			allRedirectUri := reRedirectUri.FindSubmatch(body)
			if len(allRedirectUri) >= 2 {
				redirectUri := string(allRedirectUri[1])
				bot.RedirectUri = redirectUri + "&fun=new"
				bot.BaseUri = string([]byte(bot.RedirectUri)[0:strings.LastIndex(bot.RedirectUri, "/")])

			}
		} else if code == "408" {
			fmt.Println("请求login超时")
		} else {
			panic("请求login返回:" + code)
		}
		return code
	} else {
		panic("请求login出现错误")
		return "-1"
	}
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

func (bot *WeixinBot) Login() {
	resp, err := http.Get(bot.RedirectUri)
	if err != nil {
		panic("请求login出现错误")
	}
	defer resp.Body.Close()
	doc, htmlErr := html.Parse(resp.Body)
	if htmlErr != nil {
		fmt.Println(htmlErr.Error())
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		name := strings.TrimSpace(n.Data)
		data := ""
		if n.FirstChild != nil {
			data = strings.TrimSpace(n.FirstChild.Data)
		}

		if name == "skey" {

			bot.SKey = data
		} else if name == "wxsid" {
			bot.WxSid = data
		} else if name == "wxuin" {
			wxUin, _ := strconv.ParseInt(data, 10, 64)
			bot.WxUin = wxUin
		} else if name == "pass_ticket" {
			bot.PassTicket = data
		}
		fmt.Println(name + ":" + data)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	//html := LoginHtml{}
	//xmlErr := xml.Unmarshal(body, &html)
	//if xmlErr != nil {
	//	fmt.Println(xmlErr.Error())
	//}
	//fmt.Println(html)
}

type BaseRequest struct {
	Uin      int64
	Sid      string
	Skey     string
	DeviceID string
}

func (bot *WeixinBot) InitBaseRequest() {
	deviceId := string([]byte(fmt.Sprint(rand.Float64()))[2:17])
	baseRequest := &BaseRequest{
		Uin:      bot.WxUin,
		Sid:      bot.WxSid,
		Skey:     bot.SKey,
		DeviceID: deviceId,
	}
	bot.BaseRequest = baseRequest
}

func (bot *WeixinBot) SimplePostJson(uri string, params interface{}) (b []byte, err error) {

	paramsBytes, paramsErr := json.Marshal(params)
	if paramsErr != nil {
		return nil, paramsErr
	}
	resp, err := http.Post(bot.BaseUri+uri, "application/json", bytes.NewReader(paramsBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (bot *WeixinBot) PostJson(uri string, request interface{}, response *interface{}) {
	paramsBytes, paramsErr := json.Marshal(request)
	if paramsErr != nil {
		return
	}

	resp, err := http.Post(bot.BaseUri+uri, "application/json", bytes.NewReader(paramsBytes))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	json.Unmarshal(body, response)
	return
}

func (bot *WeixinBot) GetJson(uri string, request *interface{}, response *interface{}) {
	var params url.Values
	if request != nil {
		var paramsErr error
		params, paramsErr = query.Values(request)
		if paramsErr != nil {
			return
		}
	}
	var targetUrl = ""
	if strings.Contains(uri, `http://`) {
		targetUrl = uri
	} else {
		targetUrl = bot.BaseUri + uri
	}
	resp, err := http.Get(targetUrl + "?" + params.Encode())
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	json.Unmarshal(body, response)
	return
}

type InitWebWeixinRequestBody struct {
	BaseRequest *BaseRequest
}

type User struct {
	Uin        int64
	UserName   string
	NickName   string
	HeadImgUrl string
}

type SyncKey struct {
	Count int64
	List  []struct {
		Key int64
		Val int64
	}
}

type BaseResponse struct {
	Ret int64
}
type InitWebWeixinResponseBody struct {
	BaseResponse *BaseResponse
	SyncKey      *SyncKey
	User         *User
}

func (bot *WeixinBot) InitWebWeixin() int64 {
	requestBody := InitWebWeixinRequestBody{
		BaseRequest: bot.BaseRequest,
	}
	respBody, _ := bot.SimplePostJson(fmt.Sprintf("/webwxinit?pass_ticket=%s&skey=%s&r=%s", bot.PassTicket, bot.SKey, bot.timestamp()), requestBody)

	respJson := InitWebWeixinResponseBody{}

	errJson := json.Unmarshal(respBody, &respJson)
	if errJson != nil {
		fmt.Println(errJson.Error())
	}
	bot.My = respJson.User
	bot.SyncKey = respJson.SyncKey
	syncKeyList := make([]string, bot.SyncKey.Count)
	for i, v := range bot.SyncKey.List {
		syncKeyList[i] = strconv.FormatInt(v.Key, 10) + "_" + strconv.FormatInt(v.Val, 10)
	}
	bot.SyncKeyString = strings.Join(syncKeyList, "|")

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

func (bot *WeixinBot) WebWeixinStatusNotify() bool {
	requestBody := WebWeixinStatusNotifyRequest{
		BaseRequest:  bot.BaseRequest,
		Code:         int64(3),
		FromUserName: bot.My.UserName,
		ToUserName:   bot.My.UserName,
		ClientMsgId:  time.Now().Unix(),
	}

	respBody, _ := bot.SimplePostJson(fmt.Sprintf("/webwxstatusnotify?lang=zh_CN&pass_ticket=%s", bot.PassTicket), requestBody)
	fmt.Println(string(respBody))
	respJson := WebWeixinStatusNotifyResponseBody{}

	errJson := json.Unmarshal(respBody, &respJson)

	if errJson != nil {
		fmt.Println(errJson.Error())
	}

	return respJson.BaseResponse.Ret == int64(0)
}

type SyncCheckRequest struct {
	r        int64
	sid      string
	uin      string
	skey     string
	deviceid string
	synckey  string
	_        int64
}

type SyncCheckResponseBody struct {
}

func (bot *WeixinBot) SyncCheck() {
	requestBody := SyncCheckRequest{
		r:       time.Now().Unix(),
		sid:     bot.WxSid,
		uin:     bot.WxUin,
		skey:    bot.SKey,
		synckey: bot.SyncKeyString,
		_:       time.Now().Unix(),
	}
	var responseBody string
	bot.GetJson("https://webpush.weixin.qq.com/cgi-bin/mmwebwx-bin/synccheck", &requestBody, &responseBody)

}

//def synccheck(self):
//params = {
//'r': int(time.time()),
//'sid': self.wxsid,
//'uin': self.wxuin,
//'skey': self.skey,
//'deviceid': self.device_id,
//'synckey': self.synckey,
//'_': int(time.time()),
//}
//url = 'https://webpush.weixin.qq.com/cgi-bin/mmwebwx-bin/synccheck?' + urlencode(params)
//data = get(url)
//
//pm = re.search(r'window.synccheck={retcode:"(\d+)",selector:"(\d+)"}', data)
//retcode = pm.group(1)
//selector = pm.group(2)
//
//return [retcode, selector]
//
//url = self.base_uri + '/webwxstatusnotify?lang=zh_CN&pass_ticket=%s' % (self.pass_ticket)
//params = {
//'BaseRequest': self.base_request,
//"Code": 3,
//"FromUserName": self.my['UserName'],
//"ToUserName": self.my['UserName'],
//"ClientMsgId": int(time.time())
//}
//dic = post(url, params)
//
//return dic['BaseResponse']['Ret'] == 0

//def web_weixin_init(self):
//url = self.base_uri + '/webwxinit?pass_ticket=%s&skey=%s&r=%s' % (self.pass_ticket, self.skey, int(time.time()))
//params = {
//'BaseRequest': self.base_request
//}
//
//dic = post(url, params)
//self.contact_list = dic['ContactList']
//self.my = dic['User']
//self.SyncKey = dic['SyncKey']
//self.synckey = '|'.join([str(keyVal['Key']) + '_' + str(keyVal['Val']) for keyVal in self.SyncKey['List']])
//
//return dic['BaseResponse']['Ret'] == 0
