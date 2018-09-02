package main

import (
	"flag"
	"log"
	"net/http"

	"GoBlogging/builder"
	"GoBlogging/config"
	"GoBlogging/layout"
)

func main() {
	configPath := flag.String("c", "config.json", "Path to config file")
	serve := flag.Bool("s", false, "Start server after building complete")
	port := flag.String("p", "8080", "Port to serve on")
	flag.Parse()

	c := config.New(*configPath)
	l := layout.New(c)
	b := builder.New(c)
	w := builder.NewWriter(c, l)

	b.Read(builder.Reader)
	b.Write(w)

	if *serve {
		dir := c.GetAbsPath(c.Output)

		http.Handle("/", http.FileServer(http.Dir(dir)))
		log.Printf("Serving %s on HTTP port: %s\n", dir, *port)
		log.Fatal(http.ListenAndServe(":"+*port, nil))
	}
}
