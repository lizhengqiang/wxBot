package bot

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/qiniu/log"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
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
		log.Println(re)
		bot.hears[RE] = cb
		return otto.TrueValue()
	})

	vm.Set("setTimeout", func(call otto.FunctionCall) otto.Value {
		cb := call.Argument(0)
		after, _ := call.Argument(1).ToInteger()

		go func() {
			time.Sleep(time.Duration(after) * time.Millisecond)
			cb.Call(call.This)
		}()

		return otto.NullValue()
	})

	vm.Set("set", func(call otto.FunctionCall) otto.Value {
		key, _ := call.Argument(0).ToString()
		value, _ := call.Argument(1).ToString()
		bot.Set("/otto/"+key, value)
		return otto.NullValue()
	})

	vm.Set("get", func(call otto.FunctionCall) otto.Value {
		key, _ := call.Argument(0).ToString()
		value := bot.Get("/otto/" + key)
		result, _ := otto.ToValue(value)
		return result
	})

	vm.Set("query", func(call otto.FunctionCall) otto.Value {
		dbAddr, _ := call.Argument(0).ToString()
		SQL, _ := call.Argument(1).ToString()
		db, err := sqlx.Connect("mysql", dbAddr)

		if err != nil {
			return otto.NullValue()
		}

		rows, err := db.Query(SQL)
		if err != nil {
			return otto.NullValue()
		}

		defer rows.Close()

		result := ""

		colums, err := rows.Columns()
		if err != nil {
			return otto.NullValue()
		}
		row := make([]interface{}, len(colums))
		for i, _ := range colums {
			t := ""
			row[i] = &t
		}

		result += (strings.Join(colums, "\t"))
		result += "\n"

		for rows.Next() {

			err := rows.Scan(row...)
			if err != nil {
				log.Println(err)
			}

			for _, item := range row {
				result += (*(item.(*string)) + "\t")
			}

			result += "\n"

		}

		r, _ := otto.ToValue(result)
		return r

	})
	vm.Run(string(code))
	bot.otto = vm
}
