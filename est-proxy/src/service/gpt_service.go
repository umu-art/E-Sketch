package service

import (
	"est-proxy/src/errors"
	"est_proxy_go/models"
	"net/http"
)

type GptService interface {
	Request(requestDto models.GptRequestDto, request *http.Request) (models.GptResponseDto, *errors.StatusError)
}
