package main

import (
	"est-proxy/src/http"
	"est-proxy/src/listener"
	"est-proxy/src/service"
	estbackapi "est_back_go"
	"log"
	nethttp "net/http"
	"os"
)

const kDBConfig string = "postgres://est-admin:wjq49t3q_i29f@postgres.est-dbs.svc.cluster.local:5432/est-data"

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

	// PostgreSQL
	userRepository := service.NewUserRepository(kDBConfig)
	defer userRepository.Release()

	// Хандлеры
	boardListener := listener.NewBoardListener(backApi.BoardAPI)
	userListener := listener.NewUserListener(userRepository)

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
