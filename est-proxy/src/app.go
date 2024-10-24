package main

import (
	"est-proxy/src/http"
	"est-proxy/src/listener"
	estbackapi "est_back_go"
	"log"
	nethttp "net/http"
	"os"
)

func main() {
	checkEnv()

	// Клиент для est-back
	backApiConfig := estbackapi.NewConfiguration()
	backApiConfig.Servers = estbackapi.ServerConfigurations{
		{
			URL: os.Getenv("EST_BACK_URL"),
		},
	}
	backApiConfig.HTTPClient = &nethttp.Client{
		Transport: http.NewTransportWithTraceparentHeaders(),
	}
	backApi := estbackapi.NewAPIClient(backApiConfig)

	// Хандлеры
	boardListener := listener.NewBoardListener(backApi.BoardAPI)
	userListener := listener.NewUserListener()

	// HTTP сервер
	echoListener := http.NewListener(boardListener, userListener)

	go echoListener.Serve()

	log.Println("est-proxy started")
	select {}
}

func checkEnv() {
	if os.Getenv("EST_BACK_URL") == "" {
		log.Fatal("EST_BACK_URL is not set")
	}
}
