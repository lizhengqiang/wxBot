package bot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

func (bot *WeixinBot) Tuling(content string) {
	args := strings.Split(content, " ")
	bot.setProperty("tuling.contact", args[1])
	bot.setProperty("tuling.group", args[2])
	bot.setProperty("tuling.mp", args[3])
	bot.fileHelperResponse("联系人:" + args[1] + ",群:" + args[2] + ",公众号:" + args[3])
}

type turingReq struct {
	Key    string `json:"key"`
	Info   string `json:"info"`
	UserID string `json:"userid"`
}
type turingResp struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func (bot *WeixinBot) callTuling(content, username string) (r string, err error) {
	req := &turingReq{Info: content, Key: "4aa2411ed509a4f7209d95b7ee4dfc9a", UserID: username}
	tResp := &turingResp{}
	buf := bytes.NewBuffer(nil)
	data, _ := json.Marshal(req)
	buf.Write(data)
	resp, err := http.Post("http://www.tuling123.com/openapi/api", "application/json", buf)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}
	d := json.NewDecoder(resp.Body)
	err = d.Decode(tResp)
	if err != nil {
		return
	}
	return tResp.Text, nil
}
