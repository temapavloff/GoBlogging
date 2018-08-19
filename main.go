package main

import (
	"flag"

	"GoBlogging/config"
	"GoBlogging/reader"
)

func main() {
	configPath := flag.String("c", "config.json", "Path to config file")
	flag.Parse()

	c := config.New(*configPath)
	r := reader.New(c, reader.Worker)

	r.Run()
}
