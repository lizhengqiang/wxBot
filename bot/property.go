package bot

import (
	"encoding/json"
	"qiniupkg.com/x/errors.v7"
	"strconv"
)

const (
	TRUE  = "true"
	FALSE = "false"
)
const (
	IsRunning     = "isRunning"
	IsSigned      = "isSigned"
	IsSigning     = "isLoading"
	wxuni         = "wxuni"
	wxsid         = "wxsid"
	skey          = "skey"
	deviceId      = "deviceId"
	passTicket    = "pass_ticket"
	baseUri       = "baseUri"
	syncKeyString = "syncKeyString"
	UUID          = "UUID"
	tip           = "tip"
	syncKeyHost   = "syncKeyHost"
)

const (
	me      = "me"
	syncKey = "syncKey"
)

const (
	memberList  = "memberList"
	mpList      = "mpList"
	groupList   = "groupList"
	contactList = "contactList"
)

var (
	ErrCacheMiss = errors.New("找不到数据")
)

func (bot *WeixinBot) Get(key string) (value string) {
	value = bot.Cacher.Get("data/" + key)
	return
}
func (bot *WeixinBot) unmarshal(key string, value interface{}) {
	str := bot.Cacher.Get("data/" + key)
	json.Unmarshal([]byte(str), value)
}

func (bot *WeixinBot) Set(key string, value string) {
	bot.Cacher.Set("data/"+key, value)
}

func (bot *WeixinBot) marshal(key string, value interface{}) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return
	}
	bot.Cacher.Set("data/"+key, string(bytes))
	return
}

func (bot *WeixinBot) getProperty(key string) (value string) {
	value = bot.Get("property/" + key)
	return
}

func (bot *WeixinBot) setProperty(key string, value string) {
	bot.Set("property/"+key, value)
}

type User struct {
	Uin        int64
	UserName   string
	NickName   string
	HeadImgUrl string
}

func (bot *WeixinBot) GetMe() (user *User) {
	user = &User{}
	bot.unmarshal(me, user)
	return user
}

func (bot *WeixinBot) getBaseRequest() (req *BaseRequest) {
	wxUni, _ := strconv.ParseInt(bot.getProperty(wxuni), 10, 64)
	return &BaseRequest{
		Uin:      wxUni,
		Sid:      bot.getProperty(wxsid),
		Skey:     bot.getProperty(skey),
		DeviceID: bot.getProperty(deviceId),
	}
}

func (bot *WeixinBot) IsRunning() bool {
	return bot.Get(IsRunning) == TRUE
}

func (bot *WeixinBot) IsSigned() bool {
	return bot.Get(IsSigning) == TRUE
}
