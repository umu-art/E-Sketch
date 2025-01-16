package listener

import (
	"est-proxy/src/service"
	"est-proxy/src/service/impl"
	"est_proxy_go/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GptListener struct {
	gptService service.GptService
}

func NewGptListener(gptService service.GptService) *GptListener {
	return &GptListener{gptService: gptService}
}

func (g GptListener) Request(ctx echo.Context) error {
	sessionUserId := impl.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	var requestDto models.GptRequestDto
	err := ctx.Bind(&requestDto)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорректный запрос: %v", err))
	}

	responseDto, statusError := g.gptService.Request(requestDto, ctx.Request())
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return ctx.JSON(http.StatusOK, responseDto)
}
