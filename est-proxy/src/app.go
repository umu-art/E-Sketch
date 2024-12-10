package main

import (
	"est-proxy/src/config"
	"est-proxy/src/http"
	"est-proxy/src/listener"

	repoimpl "est-proxy/src/repository/impl"
	serviceimpl "est-proxy/src/service/impl"
	wsimpl "est-proxy/src/ws/impl"

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

	// RabbitMQ
	rabbitService := repoimpl.NewRabbitRepositoryImpl()
	defer rabbitService.Close()
	figureTopic := rabbitService.GetTopic(config.RABBITMQ_TOPIC_EXCHANGE)

	// PostgreSQL
	postgresService := repoimpl.NewPostgresServiceImpl()
	defer postgresService.Release()

	// UserService
	userRepository := repoimpl.NewUserRepositoryImpl(postgresService)
	userService := serviceimpl.NewUserServiceImpl(userRepository)

	//BoardService
	boardService := serviceimpl.NewBoardServiceImpl(backApi.BoardAPI, userRepository)

	//FigureService
	figureService := serviceimpl.NewWsFigureServiceImpl(
		wsimpl.NewChannelImpl(),
		backApi.FigureAPI,
		backApi.BoardAPI,
		figureTopic,
	)

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
