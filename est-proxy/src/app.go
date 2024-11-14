package main

import (
	"est-proxy/src/config"
	"est-proxy/src/http"
	"est-proxy/src/listener"
	"est-proxy/src/service"
	estbackapi "est_back_go"
	"log"
	nethttp "net/http"
)

func main() {
	checkEnv()

	// Клиент для est-back
	backApiConfig := estbackapi.NewConfiguration()
	backApiConfig.Servers = estbackapi.ServerConfigurations{
		{
			URL: config.EST_BACK_URL,
		},
	}
	backApiConfig.HTTPClient = &nethttp.Client{
		Transport: http.NewTransportWithTraceparentHeaders(),
	}
	backApi := estbackapi.NewAPIClient(backApiConfig)

	// PostgreSQL
	userRepository := service.NewUserRepository()
	defer userRepository.Release()

	// Хандлеры
	boardListener := listener.NewBoardListener(backApi.BoardAPI, userRepository)
	userListener := listener.NewUserListener(userRepository)

	// HTTP сервер
	echoListener := http.NewListener(boardListener, userListener)

	go echoListener.Serve()

	log.Println("est-proxy started")
	select {}
}

func checkEnv() {
	if config.EST_BACK_URL == "" {
		log.Fatal("EST_BACK_URL is not set")
	}
}
