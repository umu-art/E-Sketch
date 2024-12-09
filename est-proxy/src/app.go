package main

import (
	"est-proxy/src/config"
	"est-proxy/src/http"
	"est-proxy/src/listener"
	"est-proxy/src/repository/postgres"
	"est-proxy/src/repository/user_repository"
	"est-proxy/src/service"
	"est-proxy/src/ws/ws_channel"
	estbackapi "est_back_go"
	"log"
	nethttp "net/http"
)

func main() {
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
	postgresService := postgres.NewPostgresService()
	defer postgresService.Release()

	// UserService
	userRepository := user_repository.NewUserRepository(postgresService)
	userService := service.NewUserService(userRepository)

	//BoardService
	boardService := service.NewBoardService(backApi.BoardAPI, userRepository)

	//FigureService
	figureService := service.NewWsFigureService(ws_channel.NewChannel())

	// Хандлеры
	boardListener := listener.NewBoardListener(boardService)
	userListener := listener.NewUserListener(userService)
	figureListener := listener.NewWsFigureListener(figureService)

	// Проксирование запросов

	// HTTP сервер
	echoListener := http.NewListener(boardListener, userListener, figureListener)

	go echoListener.Serve()

	log.Println("est-proxy started")
	select {}
}
