package main

import (
	"est-proxy/src/api"
	"est-proxy/src/config"
	"est-proxy/src/http"
	"est-proxy/src/listener"
	"est-proxy/src/utils"

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
		Transport: http.NewTransportWithTraceparentHeaders("est-back"),
	}
	backApi := estbackapi.NewAPIClient(backApiConfig)

	// Клиент для est-preview
	previewApiClient := &nethttp.Client{
		Transport: http.NewTransportWithTraceparentHeaders("est-preview"),
	}
	previewApi := api.NewPreviewApi(previewApiClient)

	// Клиент для GPT
	gptApiClient := &nethttp.Client{
		Transport: http.NewTransportWithTraceparentHeaders("gpt"),
	}
	gptApi := api.NewGptApi(gptApiClient)

	// Mail api
	mailApi := api.NewMailApi()

	// RabbitMQ
	rabbitService := repoimpl.NewRabbitRepositoryImpl()
	defer rabbitService.Close()
	go rabbitService.Refresh()
	figureTopic := rabbitService.GetTopic(config.RABBITMQ_FIGURE_TOPIC_EXCHANGE)
	markerTopic := rabbitService.GetTopic(config.RABBITMQ_MARKER_TOPIC_EXCHANGE)

	// Redis
	redisClient := repoimpl.NewRedisClientImpl()
	defer redisClient.Close()
	go redisClient.Refresh()

	// PostgreSQL
	postgresService := repoimpl.NewPostgresServiceImpl()
	defer postgresService.Release()

	// UserService
	userRepository := repoimpl.NewUserRepositoryImpl(postgresService)
	userService := serviceimpl.NewUserServiceImpl(userRepository, mailApi, redisClient)

	//FigureBuffer
	figureBuffer := utils.NewFigureBuffer()

	//BoardService
	boardService := serviceimpl.NewBoardServiceImpl(backApi.BoardAPI, previewApi, userRepository)

	//FigureService
	figureService := serviceimpl.NewWsFigureServiceImpl(
		wsimpl.NewChannelImpl(),
		backApi.FigureAPI,
		backApi.BoardAPI,
		figureTopic,
		figureBuffer,
	)
	go figureBuffer.ServeFlush(figureService.UpdateFigure)

	//MarkerService
	markerService := serviceimpl.NewWsMarkerServiceImpl(
		wsimpl.NewChannelImpl(),
		backApi.BoardAPI,
		userService,
		markerTopic,
	)

	//GptService
	gptService := serviceimpl.NewGptServiceImpl(previewApi, gptApi)

	// Хандлеры
	boardListener := listener.NewBoardListener(boardService)
	userListener := listener.NewUserListener(userService)
	figureListener := listener.NewWsFigureListener(figureService)
	markerListener := listener.NewWsMarkerListener(markerService)
	gptListener := listener.NewGptListener(gptService)

	// HTTP сервер
	echoListener := http.NewListener(
		boardListener,
		userListener,
		figureListener,
		markerListener,
		gptListener,
	)

	go echoListener.Serve()

	log.Println("est-proxy started")
	select {}
}
