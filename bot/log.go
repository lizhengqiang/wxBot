package bot

import (
	"encoding/json"
	"fmt"
	"qiniupkg.com/x/log.v7"
)

// 写入日志
func (this *WeixinBot) appendLog(log string) {
	var logs []string
	str := this.Cacher.Get("logs")
	err := json.Unmarshal([]byte(str), logs)
	if err != nil {
		logs = []string{}
	}
	logs = append(logs, log)

	bytes, err := json.Marshal(logs)
	if err == nil {
		this.Cacher.Set("logs", string(bytes))
	}

}

func (this *WeixinBot) Println(v ...interface{}) {
	log.Println(v...)
	this.appendLog(fmt.Sprintln(v...))
}

// 记录一条日志
func (this *WeixinBot) log(format string, v ...interface{}) {
	if len(v) > 0 {
		log.Printf(format+"\n", v...)
		this.appendLog(fmt.Sprintf(format, v...))
	} else {
		log.Printf(format + "\n")
		this.appendLog(fmt.Sprintf(format))
	}
}
