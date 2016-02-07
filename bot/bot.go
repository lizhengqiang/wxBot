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
	MpList      []Contact

	SyncKey       *SyncKey
	SyncKeyString string

	Hooks map[string]string

	Logs []string

	running bool
}

func (bot *WeixinBot) timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// 记录一条日志
func (bot *WeixinBot) log(format string, v ...interface{}) {

	if len(v) > 0 {
		fmt.Printf(format+"\n", v...)
		bot.Logs = append(bot.Logs, fmt.Sprintf(format, v...))
	} else {
		fmt.Printf(format + "\n")
		bot.Logs = append(bot.Logs, fmt.Sprintf(format))
	}
}

func (bot *WeixinBot) RegisterHookUrl(hookMethod string, hookUrl string) {
	bot.Hooks[hookMethod] = hookUrl
}

type HookMessageRequest struct {
	Method        string
	UserName      string
	GroupUserName string
	Content       string
}

type HookMessageResponse struct {
	Method        string
	UserName      string
	GroupUserName string
	Content       string
}

type HookCommonMessageRequest struct {
	HookMessageRequest
	Name      string
	GroupName string
}

func (bot *WeixinBot) hookMessage(Method, UserName, GroupUserName, Content string) {
	request := HookCommonMessageRequest{
		HookMessageRequest: HookMessageRequest{
			Method:        Method,
			UserName:      UserName,
			GroupUserName: GroupUserName,
			Content:       Content,
		},

		Name:      bot.GetRemarkName(UserName),
		GroupName: bot.GetRemarkName(GroupUserName),
	}
	response := HookMessageResponse{}
	bot.PostJson(bot.Hooks[Method], request, &response)
	if response.GroupUserName == "" {
		bot.SendMsg(response.Content, response.UserName)
	} else {
		bot.SendMsg(response.Content, response.GroupUserName)
	}
}

type HookMoneyRequest struct {
	HookMessageRequest
	Name      string
	GroupName string
}

func (bot *WeixinBot) hookMoney(UserName, GroupUserName, Content string) {
	request := HookMoneyRequest{
		HookMessageRequest: HookMessageRequest{
			Method:        "money",
			UserName:      UserName,
			GroupUserName: GroupUserName,
			Content:       Content,
		},
		Name:      bot.GetRemarkName(UserName),
		GroupName: bot.GetRemarkName(GroupUserName),
	}
	response := HookMessageResponse{}
	bot.PostJson(bot.Hooks["money"], request, &response)
	if response.GroupUserName == "" {
		bot.SendMsg(response.Content, response.UserName)
	} else {
		bot.SendMsg(response.Content, response.GroupUserName)
	}
}

// 停止这个会话
func (bot *WeixinBot) Stop() {
	bot.running = false
}

// 开启这个会话
func (bot *WeixinBot) Start() {
	// 等待登陆
	if !bot.WaitForLogin() {
		bot.log("扫描验证码失败. ")
		return
	}
	// 登陆
	if !bot.Login() {
		bot.log("登陆失败. ")
		return
	}
	// 初始化信息
	bot.InitBaseRequest()
	bot.InitWebWeixin()
	// 获取联系人列表
	bot.WebWeixinStatusNotify()
	bot.GetContact()
	// 开始监听消息
	bot.ListenMsgMode()
}

// 初始化这个会话
func (bot *WeixinBot) Init() {
	bot.running = true
	bot.Hooks = make(map[string]string)
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
			bot.log("! 初始化失败. %s", string(code))
		}
	}

	bot.log("* 初始化成功.")
}

// 获取二维码地址
func (bot *WeixinBot) GetQrcodeUrl() string {
	return "https://login.weixin.qq.com/qrcode/" + bot.UUID
}

// 等待登陆
func (bot *WeixinBot) WaitForLogin() bool {
	for bot.running {
		// 获取登陆返回值
		all, body := func() ([][]byte, []byte) {
			resp, err := bot.HttpClient.Get(fmt.Sprintf("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?tip=%s&uuid=%s&_=%s", bot.Tip, bot.UUID, bot.timestamp()))
			if err != nil {
				bot.log(err.Error())
			}
			if resp.Body == nil {
				return nil, nil
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
				bot.log("* 成功扫描,请在手机上点击确认以登录.")
				continue
			} else if code == "200" {
				reRedirectUri, _ := regexp.Compile(`window.redirect_uri="(\S+?)";`)
				allRedirectUri := reRedirectUri.FindSubmatch(body)
				if len(allRedirectUri) >= 2 {
					redirectUri := string(allRedirectUri[1])
					bot.RedirectUri = redirectUri + "&fun=new"
					bot.BaseUri = string([]byte(bot.RedirectUri)[0:strings.LastIndex(bot.RedirectUri, "/")])

				}
				bot.log("* 登陆成功.")
				return true
			} else if code == "408" {
				bot.log("! 登陆超时.")
			} else {
				bot.log("! 登录失败 %s", code)
				return false
			}
		}
		time.Sleep(time.Second * 3)
	}
	return false

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
func (bot *WeixinBot) Login() bool {
	resp, err := bot.HttpClient.Get(bot.RedirectUri)
	if err != nil {
		bot.log(err.Error())
		return false
	}
	if resp.Body == nil {
		return false
	}
	defer resp.Body.Close()
	doc, htmlErr := html.Parse(resp.Body)
	if htmlErr != nil {
		bot.log(htmlErr.Error())
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
	return true
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
	resp, err := bot.HttpClient.Post(targetUrl, "application/json", bytes.NewReader(paramsBytes))
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
		bot.log(errJson.Error())
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
		bot.log(errJson.Error())
	}

	return respJson.BaseResponse.Ret == int64(0)
}

type EmptyRequest struct {
}

type Contact struct {
	VerifyFlag  int64
	UserName    string
	RemarkName  string
	NickName    string
	DisplayName string
}
type GetContactResponse struct {
	MemberList []Contact
}

// 获取联系人列表
func (bot *WeixinBot) GetContact() bool {
	// SpecialUsers := []string{"newsapp", "fmessage", "filehelper", "weibo", "qqmail", "fmessage", "tmessage", "qmessage", "qqsync", "floatbottle", "lbsapp", "shakeapp", "medianote", "qqfriend", "readerapp", "blogapp", "facebookapp", "masssendapp", "meishiapp", "feedsapp", "voip", "blogappweixin", "weixin", "brandsessionholder", "weixinreminder", "wxid_novlwrv3lqwv11", "gh_22b87fa7cb3c", "officialaccounts", "notification_messages", "wxid_novlwrv3lqwv11", "gh_22b87fa7cb3c", "wxitil", "userexperience_alarm", "notification_messages"}
	response := GetContactResponse{}
	bot.PostJson(fmt.Sprintf("/webwxgetcontact?pass_ticket=%s&skey=%s&r=%s", bot.PassTicket, bot.SKey, bot.timestamp()), &EmptyRequest{}, &response)
	bot.MemberList = response.MemberList
	for _, contact := range response.MemberList {
		if contact.VerifyFlag != 0 {
			bot.MpList = append(bot.MpList, contact)
		} else if strings.Contains(contact.UserName, "@@") {
			bot.GroupList = append(bot.GroupList, contact)
		} else if strings.Contains(contact.UserName, "@") {
			bot.ContactList = append(bot.ContactList, contact)
		}
	}
	return true
}

func (bot *WeixinBot) GetRemarkName(id string) (name string) {
	for _, contact := range bot.MemberList {
		if contact.UserName == id {
			if contact.RemarkName == "" {
				return contact.NickName
			} else {
				return contact.RemarkName
			}
		}
	}
	return "未知"
}

type SyncCheckResponseBody struct {
	retcode  int64
	selector int64
}

func (bot *WeixinBot) SyncCheck() (retcode, selector int64) {
	deviceId := bot.DeviceId
	resp, err := bot.HttpClient.Get(fmt.Sprintf("https://webpush.weixin.qq.com/cgi-bin/mmwebwx-bin/synccheck?synckey=%s&skey=%s&uin=%s&r=%s&deviceid=%s&sid=%s&_=%s", url.QueryEscape(bot.SyncKeyString), url.QueryEscape(bot.SKey), strconv.FormatInt(bot.WxUin, 10), bot.timestamp(), deviceId, url.QueryEscape(bot.WxSid), bot.timestamp()))
	if err != nil {
		bot.log(err.Error())
	}
	defer resp.Body.Close()
	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		bot.log(bodyErr.Error())
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
	name := bot.GetRemarkName(UserName)
	bot.log("# 我->%s:%s", name, Content)
}

func (bot *WeixinBot) handleMsg(msgList []AddMsg) {
	for _, msg := range msgList {

		msgType := msg.MsgType
		userName := ""
		groupUserName := ""
		name := ""
		groupName := ""
		content := ""
		fromGroup := false

		if msg.ToUserName == "filehelper" {
			// 文件助手
		} else if msg.FromUserName == bot.My.UserName {
			// 自己的
			continue
		} else if strings.Contains(msg.FromUserName, "@@") {
			contents := strings.Split(msg.Content, `:<br/>`)
			userName = contents[0]
			content = contents[1]
			groupUserName = msg.FromUserName
			name = bot.GetRemarkName(userName)
			groupName = bot.GetRemarkName(groupUserName)
			fromGroup = true

		} else {
			userName = msg.FromUserName
			content = msg.Content
			name = bot.GetRemarkName(msg.FromUserName)
			fromGroup = false
		}

		if msgType == 1 {
			method := "message"
			if fromGroup {
				method = "groupMessage"
			} else {
				method = "contactMessage"
			}
			bot.hookMessage(method, userName, groupUserName, content)
			bot.log("# %s(%s):%s", name, groupName, content)
		} else if msgType == 10000 {
			bot.hookMoney(userName, groupUserName, content)
			bot.log("# %s(%s):%s", name, groupName, content)
		} else {
			bot.log("# %s(%s):%s(%s) ", name, groupName, content, strconv.FormatInt(msgType, 10))
		}
	}
}

func (bot *WeixinBot) ListenMsgMode() {
	for bot.running {
		retcode, selector := bot.SyncCheck()
		if retcode == 1100 {
			bot.log("# 你退出了,债见~_~")
			return
		} else if retcode == 0 {
			if selector == 2 {
				msgList := bot.WebWeixinSync()
				if msgList.AddMsgList != nil && len(msgList.AddMsgList) > 0 {
					bot.handleMsg(msgList.AddMsgList)
				}

			} else if selector == 7 {
				bot.log("# 发现你玩手机了!")
			} else if selector == 0 {
				time.Sleep(3 * time.Second)
			}
		}
	}
}
