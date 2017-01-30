package bot

import (
	"encoding/json"
	"fmt"
	"qiniupkg.com/x/log.v7"
)

// 写入日志
func (bot *WeixinBot) appendLog(log string) {
	var logs []string
	str := bot.Cacher.Get("logs")
	err := json.Unmarshal([]byte(str), logs)
	if err != nil {
		logs = []string{}
	}
	logs = append(logs, log)

	bytes, err := json.Marshal(logs)
	if err == nil {
		bot.Cacher.Set("logs", string(bytes))
	}

}

func (bot *WeixinBot) Println(v ...interface{}) {
	log.Println(v...)
	bot.appendLog(fmt.Sprintln(v...))
}

// 记录一条日志
func (bot *WeixinBot) log(format string, v ...interface{}) {
	if len(v) > 0 {
		log.Printf(format+"\n", v...)
		bot.appendLog(fmt.Sprintf(format, v...))
	} else {
		log.Printf(format + "\n")
		bot.appendLog(fmt.Sprintf(format))
	}
}
