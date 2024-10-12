package main

import (
	"log"
	"nb-proxy/src/http"
)

func main() {
	listener := http.NewListener()

	go listener.Serve()

	log.Println("nb-proxy started")
	select {}
}
