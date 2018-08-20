package main

import (
	"flag"

	"GoBlogging/builder"
	"GoBlogging/config"
)

func main() {
	configPath := flag.String("c", "config.json", "Path to config file")
	flag.Parse()

	c := config.New(*configPath)
	b := builder.New(c)

	b.Read(builder.Reader)
	b.Write()
}
