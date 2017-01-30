package bot

import (
	"github.com/qiniu/log"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"regexp"
)

func (bot *WeixinBot) OnMessageJS(group, user, content string) {
	vm := bot.otto
	vm.Call("onMessage", nil, group, user, content)
}
func (bot *WeixinBot) Hear(group, user, content string) {
	for RE, cb := range bot.hears {
		if RE.Match([]byte(content)) {
			cb.Call(otto.NullValue(), group, user, content)
		}
	}
}
func (bot *WeixinBot) ReloadJS() {
	code, err := ioutil.ReadFile("./scripts/main.js")
	if err != nil {
		log.Println(err)
		return
	}
	bot.hears = map[*regexp.Regexp]otto.Value{}
	vm := otto.New()
	vm.Set("me", bot.GetMe())
	vm.Set("sendMsg", func(call otto.FunctionCall) otto.Value {
		userName, _ := call.Argument(0).ToString()
		content, _ := call.Argument(1).ToString()
		bot.SendMsg(content, userName)
		return otto.NullValue()
	})

	vm.Set("reloadJS", func(call otto.FunctionCall) otto.Value {
		bot.ReloadJS()
		return otto.NullValue()
	})

	vm.Set("hear", func(call otto.FunctionCall) otto.Value {
		re, _ := call.Argument(0).ToString()
		cb := call.Argument(1)
		RE, err := regexp.Compile(re)
		if err != nil {
			return otto.FalseValue()
		}
		bot.hears[RE] = cb
		return otto.TrueValue()
	})
	vm.Run(string(code))
	bot.otto = vm
}
