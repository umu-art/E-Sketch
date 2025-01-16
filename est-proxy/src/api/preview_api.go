package api

import (
	"context"
	"est-proxy/src/config"
	"est-proxy/src/errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type PreviewApi struct {
	httpClient *http.Client
}

func NewPreviewApi(httpClient *http.Client) *PreviewApi {
	return &PreviewApi{httpClient}
}

func (p PreviewApi) GetToken(boardId string, ctx context.Context) (string, *errors.StatusError) {
	endpoint := fmt.Sprintf("%s/internal/get-token", config.EST_PREVIEW_URL)

	params := url.Values{}
	params.Add("boardId", boardId)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return "", errors.NewStatusError(http.StatusInternalServerError, "Failed to create request")
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", errors.NewStatusError(http.StatusInternalServerError, "Failed to send request to preview service")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.NewStatusError(resp.StatusCode, "Failed to get token from preview service")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.NewStatusError(http.StatusInternalServerError, "Failed to read response body")
	}

	return string(body), nil
}

func (p PreviewApi) GetPreview(
	boardId string,
	width float32, height float32,
	xLeft float32, yTop float32,
	xRight float32, yBottom float32,
	ctx context.Context,
) ([]byte, *errors.StatusError) {
	endpoint := fmt.Sprintf("%s/internal/preview", config.EST_PREVIEW_URL)

	params := url.Values{}
	params.Add("boardId", boardId)
	params.Add("width", fmt.Sprintf("%f", width))
	params.Add("height", fmt.Sprintf("%f", height))
	params.Add("xLeft", fmt.Sprintf("%f", xLeft))
	params.Add("yUp", fmt.Sprintf("%f", yTop))
	params.Add("xRight", fmt.Sprintf("%f", xRight))
	params.Add("yDown", fmt.Sprintf("%f", yBottom))

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Failed to create request")
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Failed to send request to preview service")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewStatusError(resp.StatusCode, "Failed to get preview from preview service")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Failed to read response body")
	}

	return body, nil
}
