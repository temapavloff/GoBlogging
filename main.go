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
	r := builder.New(c)

	r.Read(builder.Reader)
	r.Write()
}
