package main

import (
	"log"
	"est-proxy/src/http"
)

func main() {
	listener := http.NewListener()

	go listener.Serve()

	log.Println("est-proxy started")
	select {}
}
