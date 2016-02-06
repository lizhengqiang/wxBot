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
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type WeixinBot struct {
	HttpClient  *http.Client
	UUID        string
	Tip         string
	RedirectUri string
	BaseUri     string
	SKey        string
	WxSid       string
	WxUin       int64
	PassTicket  string
	DeviceId    string

	BaseRequest *BaseRequest

	My *User

	MemberList  []Contact
	ContactList []Contact
	GroupList   []Contact

	SyncKey       *SyncKey
	SyncKeyString string

	hooks map[string]string

	Logs []string
}

func (bot *WeixinBot) timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)

}

// 记录一条日志
func (bot *WeixinBot) log(format string, v ...interface{}) {
	if len(v) > 0 {
		fmt.Printf(format+"\n", v)
		bot.Logs = append(bot.Logs, fmt.Sprintf(format, v))
	} else {
		fmt.Printf(format + "\n")
		bot.Logs = append(bot.Logs, fmt.Sprintf(format))
	}

}

func (bot *WeixinBot) fmt() {
	fmt.Println(bot.UUID)
}

func (bot *WeixinBot) RegisterHookUrl(hookMethod string, hookUrl string) {
	bot.hooks[hookMethod] = hookUrl
}

type HookMessageRequest struct {
	Method   string
	UserName string
	Content  string
}

type HookMessageResponse struct {
	Method   string
	UserName string
	Content  string
}

func (bot *WeixinBot) hookMessage(UserName string, Content string) *HookMessageResponse {
	request := HookMessageRequest{
		Method:   "message",
		UserName: UserName,
		Content:  Content,
	}
	response := HookMessageResponse{}
	bot.PostJson(bot.hooks["message"], request, &response)
	return &response
}

func (bot *WeixinBot) Start() {
	bot.hooks = make(map[string]string)
	bot.DeviceId = "e" + string([]byte(fmt.Sprint(rand.Float64()))[2:17])
	gCurCookieJar, _ := cookiejar.New(nil)

	bot.HttpClient = &http.Client{
		Jar: gCurCookieJar,
	}

	resp, err := bot.HttpClient.PostForm("https://login.weixin.qq.com/jslogin", url.Values{"appid": {"wx782c26e4c19acffb"}, "fun": {"new"}, "lang": {"zh_CN"}, "_": {bot.timestamp()}})
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
		if string(code) == "200" {
			bot.UUID = string(uuid)
		} else {
			panic("请求jslogin返回!200")
		}
	}

	bot.log("执行Start成功, uuid is %s.", bot.UUID)
}

func (bot *WeixinBot) GetQrcodeUrl() string {
	return "https://login.weixin.qq.com/qrcode/" + bot.UUID
}

func (bot *WeixinBot) WaitForLogin() {
	for {
		// 获取登陆返回值
		all, body := func() ([][]byte, []byte) {
			resp, err := bot.HttpClient.Get(fmt.Sprintf("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?tip=%s&uuid=%s&_=%s", bot.Tip, bot.UUID, bot.timestamp()))
			if err != nil {
				bot.log(err.Error())
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			re, _ := regexp.Compile(`window.code=(\d+);`)
			all := re.FindSubmatch(body)
			return all, body
		}()

		if len(all) >= 2 {
			code := string(all[1])
			if code == "201" {
				bot.Tip = "0"
				bot.log("成功扫描, 请在手机上点击确认以登录. ")
				continue
			} else if code == "200" {
				reRedirectUri, _ := regexp.Compile(`window.redirect_uri="(\S+?)";`)
				allRedirectUri := reRedirectUri.FindSubmatch(body)
				if len(allRedirectUri) >= 2 {
					redirectUri := string(allRedirectUri[1])
					bot.RedirectUri = redirectUri + "&fun=new"
					bot.BaseUri = string([]byte(bot.RedirectUri)[0:strings.LastIndex(bot.RedirectUri, "/")])

				}
				bot.log("登陆成功. ")
				break
			} else if code == "408" {
				bot.log("请求login超时. ")
			} else {
				bot.log("请求login返回%s. ", code)
			}
		}
		time.Sleep(time.Second * 3)
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
	resp, err := bot.HttpClient.Get(bot.RedirectUri)
	if err != nil {
		panic("请求login出现错误")
	}
	defer resp.Body.Close()
	doc, htmlErr := html.Parse(resp.Body)
	if htmlErr != nil {
		//fmt.Println(htmlErr.Error())
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
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

type BaseRequest struct {
	Uin      int64
	Sid      string
	Skey     string
	DeviceID string
}

func (bot *WeixinBot) InitBaseRequest() {
	baseRequest := &BaseRequest{
		Uin:      bot.WxUin,
		Sid:      bot.WxSid,
		Skey:     bot.SKey,
		DeviceID: bot.DeviceId,
	}
	bot.BaseRequest = baseRequest
}

func (bot *WeixinBot) SimplePostJson(uri string, params interface{}) (b []byte, err error) {

	paramsBytes, paramsErr := json.Marshal(params)
	if paramsErr != nil {
		return nil, paramsErr
	}
	resp, err := bot.HttpClient.Post(bot.BaseUri+uri, "application/json", bytes.NewReader(paramsBytes))
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

func (bot *WeixinBot) PostJson(uri string, request interface{}, response interface{}) {
	paramsBytes, paramsErr := json.Marshal(request)
	if paramsErr != nil {
		return
	}
	var targetUrl = ""
	if strings.Contains(uri, `http://`) || strings.Contains(uri, `https://`) {
		targetUrl = uri
	} else {
		targetUrl = bot.BaseUri + uri
	}
	//fmt.Println(targetUrl)
	//fmt.Println(string(paramsBytes))
	resp, err := bot.HttpClient.Post(targetUrl, "application/json", bytes.NewReader(paramsBytes))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	//fmt.Println(string(body))
	json.Unmarshal(body, response)
	return
}

func (bot *WeixinBot) GetJson(uri string, request interface{}, response interface{}) {

	var targetUrl = ""
	if strings.Contains(uri, `http://`) || strings.Contains(uri, `https://`) {
		targetUrl = uri
	} else {
		targetUrl = bot.BaseUri + uri
	}
	var params url.Values
	if request != nil {
		var paramsErr error
		params, paramsErr = query.Values(request)
		if paramsErr != nil {
			return
		}
		targetUrl = targetUrl + "?" + params.Encode()
	}
	resp, err := bot.HttpClient.Get(targetUrl)
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

func (bot *WeixinBot) saveSyncKey(syncKey *SyncKey) {
	bot.SyncKey = syncKey
	syncKeyList := make([]string, bot.SyncKey.Count)
	for i, v := range bot.SyncKey.List {
		syncKeyList[i] = strconv.FormatInt(v.Key, 10) + "_" + strconv.FormatInt(v.Val, 10)
	}
	bot.SyncKeyString = strings.Join(syncKeyList, "|")
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
	bot.saveSyncKey(respJson.SyncKey)

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
	respJson := WebWeixinStatusNotifyResponseBody{}

	errJson := json.Unmarshal(respBody, &respJson)

	if errJson != nil {
		fmt.Println(errJson.Error())
	}

	return respJson.BaseResponse.Ret == int64(0)
}

type EmptyRequest struct {
}

type Contact struct {
	VerifyFlag int64
	UserName   string
	RemarkName string
	NickName   string
}
type GetContactResponse struct {
	MemberList []Contact
}

func (bot *WeixinBot) GetContact() bool {
	SpecialUsers := []string{"newsapp", "fmessage", "filehelper", "weibo", "qqmail", "fmessage", "tmessage", "qmessage", "qqsync", "floatbottle", "lbsapp", "shakeapp", "medianote", "qqfriend", "readerapp", "blogapp", "facebookapp", "masssendapp", "meishiapp", "feedsapp", "voip", "blogappweixin", "weixin", "brandsessionholder", "weixinreminder", "wxid_novlwrv3lqwv11", "gh_22b87fa7cb3c", "officialaccounts", "notification_messages", "wxid_novlwrv3lqwv11", "gh_22b87fa7cb3c", "wxitil", "userexperience_alarm", "notification_messages"}
	response := GetContactResponse{}
	bot.PostJson(fmt.Sprintf("/webwxgetcontact?pass_ticket=%s&skey=%s&r=%s", bot.PassTicket, bot.SKey, bot.timestamp()), &EmptyRequest{}, &response)
	bot.MemberList = response.MemberList
	for _, contact := range response.MemberList {
		if contact.VerifyFlag&8 != 0 {

		} else if sort.SearchStrings(SpecialUsers, contact.UserName) >= 0 {

		} else if strings.Contains(contact.UserName, "@@") {
			bot.GroupList = append(bot.GroupList, contact)
		} else {
			bot.ContactList = append(bot.ContactList, contact)
		}
	}
	return true
}

func (bot *WeixinBot) GetUserRemarkName(id string) (name string) {
	for _, contact := range bot.MemberList {
		if contact.UserName == id {
			if contact.RemarkName == "" {
				return contact.NickName
			} else {
				return contact.RemarkName
			}
		}
	}
	return "未知" + id
}

type SyncCheckResponseBody struct {
	retcode  int64
	selector int64
}

func (bot *WeixinBot) SyncCheck() (retcode, selector int64) {
	deviceId := bot.DeviceId
	resp, err := bot.HttpClient.Get(fmt.Sprintf("https://webpush.weixin.qq.com/cgi-bin/mmwebwx-bin/synccheck?synckey=%s&skey=%s&uin=%s&r=%s&deviceid=%s&sid=%s&_=%s", url.QueryEscape(bot.SyncKeyString), url.QueryEscape(bot.SKey), strconv.FormatInt(bot.WxUin, 10), bot.timestamp(), deviceId, url.QueryEscape(bot.WxSid), bot.timestamp()))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		fmt.Println(bodyErr.Error())
	}

	re, _ := regexp.Compile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)
	sub := re.FindSubmatch(body)
	if len(sub) >= 2 {
		retcode, _ := strconv.ParseInt(string(sub[1]), 10, 64)
		selector, _ := strconv.ParseInt(string(sub[2]), 10, 64)
		return retcode, selector
	} else {
		return -1, -1
	}
}

type WebWeixinSyncRequest struct {
	BaseRequest *BaseRequest
	SyncKey     *SyncKey
	rr          int64
}

type AddMsg struct {
	MsgType      int64
	FromUserName string
	ToUserName   string
	Content      string
}

type WebWeixinSyncResponse struct {
	BaseResponse *BaseResponse
	SyncKey      *SyncKey
	AddMsgList   []AddMsg
}

func (bot *WeixinBot) WebWeixinSync() WebWeixinSyncResponse {
	request := WebWeixinSyncRequest{
		BaseRequest: bot.BaseRequest,
		SyncKey:     bot.SyncKey,
		rr:          time.Now().Unix(),
	}

	response := WebWeixinSyncResponse{}

	bot.PostJson(fmt.Sprintf("/webwxsync?sid=%s&skey=%s&pass_ticket=%s", bot.WxSid, bot.SKey, bot.PassTicket), request, &response)

	if response.BaseResponse != nil && response.BaseResponse.Ret == 0 {
		bot.saveSyncKey(response.SyncKey)
	}

	return response
}

type Msg struct {
	Type         int64
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      string
	ClientMsgId  string
}

type SendMsgRequest struct {
	BaseRequest *BaseRequest
	Msg         *Msg
}

type SendMsgResponse struct {
	BaseResponse *BaseResponse
}

func (bot *WeixinBot) SendMsg(Content, UserName string) {
	clientMsgId := strconv.FormatInt(time.Now().Unix()*1000+time.Now().Unix(), 10)
	request := SendMsgRequest{
		BaseRequest: bot.BaseRequest,
		Msg: &Msg{
			Type:         1,
			Content:      Content,
			FromUserName: bot.My.UserName,
			ToUserName:   UserName,
			LocalID:      clientMsgId,
			ClientMsgId:  clientMsgId,
		},
	}
	response := SendMsgResponse{}
	bot.PostJson(fmt.Sprintf("/webwxsendmsg?pass_ticket=%s", bot.PassTicket), request, &response)
	//fmt.Println(response.BaseResponse.Ret)
}

func (bot *WeixinBot) handleMsg(msgList []AddMsg) {
	for _, msg := range msgList {
		fmt.Println("# New Message Received.")
		msgType := msg.MsgType

		if msgType == 51 {
			fmt.Println("# Received Init Message.")
		} else if msgType == 1 {
			if msg.ToUserName == "filehelper" {

			} else if msg.FromUserName == bot.My.UserName {
				continue
			} else if strings.Contains(msg.FromUserName, "@@") {
				msgContensArray := strings.Split(msg.Content, `:<br/>`)
				userName := msgContensArray[0]
				content := strings.Replace(msgContensArray[1], "<br/>", "\n", -1)
				groupUserName := msg.FromUserName
				name := bot.GetUserRemarkName(userName)
				hookMessageResponse := bot.hookMessage(userName, content)
				bot.SendMsg(hookMessageResponse.Content, groupUserName)
				fmt.Printf("# %s(%s) in %s -> %s\n", name, userName, groupUserName, content)
			} else {
				name := bot.GetUserRemarkName(msg.FromUserName)
				hookMessageResponse := bot.hookMessage(msg.FromUserName, msg.Content)
				bot.SendMsg(hookMessageResponse.Content, hookMessageResponse.UserName)
				fmt.Printf("# %s(%s) -> %s\n", name, msg.FromUserName, msg.Content)
			}
		}
	}
}

func (bot *WeixinBot) ListenMsgMode() {
	for {
		retcode, selector := bot.SyncCheck()
		if retcode == 1100 {
			fmt.Println("# u logout, bye.")
			return
		} else if retcode == 0 {
			if selector == 2 {
				msgList := bot.WebWeixinSync()
				if msgList.AddMsgList != nil && len(msgList.AddMsgList) > 0 {
					bot.handleMsg(msgList.AddMsgList)
				}

			} else if selector == 7 {
				fmt.Println("# I found u playing phone.")
			} else if selector == 0 {
				time.Sleep(3 * time.Second)
			}
		}
	}
}

//def listenMsgMode(self):
//info("进入消息监听模式")
//playWeChat = 0
//while True:
//[retcode, selector] = self.synccheck()
//if retcode == '1100':
//print '[*] 你在手机上登出了微信，债见'
//break
//elif retcode == '0':
//if selector == '2':
//r = self.webwxsync()
//if r is not None: self.handleMsg(r)
//elif selector == '7':
//playWeChat += 1
//print '[*] 你在手机上玩微信被我发现了 %d 次' % playWeChat
//r = self.webwxsync()
//elif selector == '0':
//time.sleep(1)
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
