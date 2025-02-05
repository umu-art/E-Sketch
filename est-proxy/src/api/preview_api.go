package api

import (
	"context"
	"encoding/json"
	"est-proxy/src/config"
	"est-proxy/src/errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type PreviewApi struct {
	httpClient *http.Client
}

func NewPreviewApi(httpClient *http.Client) *PreviewApi {
	return &PreviewApi{httpClient}
}

func (p PreviewApi) GetTokens(boardIds []string, ctx context.Context) (map[string]string, *errors.StatusError) {
	endpoint := fmt.Sprintf("%s/internal/get-tokens", config.EST_PREVIEW_URL)

	params := url.Values{}
	params.Add("boardIds", strings.Join(boardIds, ","))

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
		return nil, errors.NewStatusError(resp.StatusCode, "Failed to get token from preview service")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Failed to read response body")
	}

	var tokens []tokenData
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Failed to unmarshal response body")
	}

	mappedTokens := map[string]string{}
	for _, tokenData := range tokens {
		mappedTokens[tokenData.BoardID] = tokenData.Token
	}

	return mappedTokens, nil
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

type tokenData struct {
	BoardID string `json:"boardId"`
	Token   string `json:"token"`
}
