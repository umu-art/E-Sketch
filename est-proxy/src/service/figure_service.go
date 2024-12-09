package service

import (
	"est-proxy/src/service/impl"
	"est-proxy/src/ws/ws_channel"
	estbackapi "est_back_go"
	"net/http"
)

type WsFigureService interface {
	Listen(writer http.ResponseWriter, request *http.Request) error
}

func NewWsFigureService(channel ws_channel.Channel, figureApi *estbackapi.FigureAPIService) WsFigureService {
	return impl.NewWsServiceListenerImpl(channel, figureApi)
}
