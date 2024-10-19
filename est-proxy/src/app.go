package main

import (
	"est-proxy/src/http"
	"log"
)

func main() {
	listener := http.NewListener()

	go listener.Serve()

	log.Println("est-proxy started")
	select {}
}
