package impl

import (
	"est-proxy/src/api"
	"est-proxy/src/errors"
	"est_proxy_go/models"
	"net/http"
)

type GptServiceImpl struct {
	previewApi *api.PreviewApi
	gptApi     *api.GptApi
}

func NewGptServiceImpl(previewApi *api.PreviewApi, gptApi *api.GptApi) *GptServiceImpl {
	return &GptServiceImpl{previewApi, gptApi}
}

func (g GptServiceImpl) Request(
	requestDto models.GptRequestDto,
	request *http.Request,
) (models.GptResponseDto, *errors.StatusError) {
	image, err := g.previewApi.GetPreview(requestDto.BoardId,
		1200, 800,
		requestDto.LeftUp.X, requestDto.LeftUp.Y,
		requestDto.RightDown.X, requestDto.RightDown.Y,
		request.Context())

	if err != nil {
		return models.GptResponseDto{}, err
	}

	gptResponse, statusError := g.gptApi.Request(requestDto.Prompt, image, request.Context())

	if statusError != nil {
		return models.GptResponseDto{}, statusError
	}

	return models.GptResponseDto{
		Text: gptResponse,
	}, nil
}
