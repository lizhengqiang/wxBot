package main

import (
	"github.com/cocotyty/summer"
	"wxBot/api"
	"wxBot/provider"
)

func main() {
	summer.Start()

	m := summer.GetStoneWithName("HttpServer").(*provider.HttpServer).M

	api.RegisterRoutes(m)

	m.Run()
}
