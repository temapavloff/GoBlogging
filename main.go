package main

import (
	"flag"

	"GoBlogging/builder"
	"GoBlogging/config"
	"GoBlogging/layout"
)

func main() {
	configPath := flag.String("c", "config.json", "Path to config file")
	flag.Parse()

	c := config.New(*configPath)
	l := layout.New(c)
	b := builder.New(c)
	w := builder.NewWriter(c, l)

	b.Read(builder.Reader)
	b.Write(w)
}
