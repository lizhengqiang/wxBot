package bot

import (
	"encoding/json"
	"qiniupkg.com/x/errors.v7"
	"strconv"
)

const (
	TRUE = "true"
	FALSE = "false"
)
const (
	IsRunning = "isRunning"
	isLogining = "isLoading"
	wxuni = "wxuni"
	wxsid = "wxsid"
	skey = "skey"
	deviceId = "deviceId"
	passTicket = "pass_ticket"
	baseUri = "baseUri"
	syncKeyString = "syncKeyString"
	UUID = "UUID"
	tip = "tip"
)

const (
	me = "me"
	syncKey = "syncKey"
)

const (
	memberList = "memberList"
	mpList = "mpList"
	groupList = "groupList"
	contactList = "contactList"
)

var (
	ErrCacheMiss = errors.New("找不到数据")
)

func (this *WeixinBot) Get(key string) (value string) {
	value = this.Cacher.Get("data/" + key)
	return
}
func (this *WeixinBot) unmarshal(key string, value interface{}) {
	str := this.Cacher.Get("data/" + key)
	json.Unmarshal([]byte(str), value)
}

func (this *WeixinBot) Set(key string, value string) {
	this.Cacher.Set("data/" + key, value)
}

func (this *WeixinBot) marshal(key string, value interface{}) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return
	}
	this.Cacher.Set("data/" + key, string(bytes))
	return
}

func (this *WeixinBot) getProperty(key string) (value string) {
	value = this.Get("property/" + key)
	return
}

func (this *WeixinBot) setProperty(key string, value string) {
	this.Set("property/" + key, value)
}



type User struct {
	Uin        int64
	UserName   string
	NickName   string
	HeadImgUrl string
}



func (this *WeixinBot) GetMe() (user *User) {
	user = &User{}
	this.unmarshal(me, user)
	return user
}

func (this *WeixinBot) getBaseRequest() (req *BaseRequest) {
	wxUni, _ := strconv.ParseInt(this.getProperty(wxuni), 10, 64)
	return &BaseRequest{
		Uin:      wxUni,
		Sid:      this.getProperty(wxsid),
		Skey:     this.getProperty(skey),
		DeviceID: this.getProperty(deviceId),
	}
}

func (bot *WeixinBot) IsRunning() bool {
	return bot.Get(IsRunning) == TRUE
}

func (this *WeixinBot) IsLogining() bool {
	return this.Get(isLogining) == TRUE
}
