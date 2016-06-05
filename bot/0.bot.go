package bot

import (
	"net/http"
	"strconv"
	"time"
)

// TODO 得先把这个改成分布式
type WeixinBot struct {
	ID         string
	httpClient *http.Client
	Cacher     Cache
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
