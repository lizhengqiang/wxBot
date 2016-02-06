package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"
	"wxBot/bot"
)

// Encode via Gob to file
func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Decode Gob file
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func main() {
	bot := &bot.WeixinBot{}
	if len(os.Args) <= 1 || os.Args[1] != "reload" {
		Load("./bot.obj", bot)
	}

	fmt.Println(bot)
	if bot.UUID == "" {
		bot.Start()
		fmt.Println(bot.UUID)
	}

	if bot.RedirectUri == "" {
		qrcodeUrl := bot.GetQrcodeUrl()
		fmt.Println(qrcodeUrl)
		for code := bot.WaitForLogin(); code != "200"; code = bot.WaitForLogin() {
			fmt.Println(code)
			time.Sleep(time.Second * 3)
		}
	}

	if bot.PassTicket == "" || bot.WxUin == 0 || bot.WxSid == "" || bot.SKey == "" {
		bot.Login()
		bot.InitBaseRequest()
	}

	if bot.My == nil || bot.SyncKey == nil {
		bot.InitWebWeixin()
	}

	fmt.Println(bot.WebWeixinStatusNotify())

	//fmt.Println(bot)
	bot.ListenMsgMode()
	bot.SyncCheck()

	Save("./bot.obj", bot)
}
