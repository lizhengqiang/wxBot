package main

import (
	"github.com/cocotyty/summer"

	_ "github.com/lizhengqiang/wxBot/api"
	_ "github.com/lizhengqiang/wxBot/bot_router"
	"github.com/lizhengqiang/wxBot/provider"
)

func main() {
	summer.Start()

	m := summer.GetStoneWithName("HttpServer").(*provider.HttpServer).M

	m.Run()
}
