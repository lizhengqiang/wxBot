package bot

import (
	"github.com/lizhengqiang/wxBot/domain"
	"github.com/robertkrimen/otto"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// TODO 得先把这个改成分布式
type WeixinBot struct {
	ID            string
	httpClient    *http.Client
	Cacher        Cache
	MQ            domain.MessageQueue
	newMsgLock    sync.Mutex
	rw            sync.Mutex
	IsLoopRunning bool
	otto          *otto.Otto
	hears         map[*regexp.Regexp]otto.Value
}

func (bot *WeixinBot) timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

type BaseRequest struct {
	Uin      int64
	Sid      string
	Skey     string
	DeviceID string
}

type BaseResponse struct {
	Ret int64
}

type EmptyRequest struct {
}
