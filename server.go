package main

import (
	"github.com/cocotyty/summer"
	"github.com/lizhengqiang/wxBot/api"
	"github.com/lizhengqiang/wxBot/provider"
)

func main() {
	summer.Start()

	m := summer.GetStoneWithName("HttpServer").(*provider.HttpServer).M

	api.RegisterRoutes(m)

	m.Run()
}
