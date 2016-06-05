package main

import (
	"github.com/cocotyty/summer"
	"os"
	"wxBot/api"
	"wxBot/provider"
)

func main() {

	summer.TomlFile(os.Getenv("config"))
	summer.Start()

	m := summer.GetStoneWithName("HttpServer").(*provider.HttpServer).M

	api.RegisterRoutes(m)

	m.Run()
}
