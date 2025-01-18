package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"est-proxy/src/config"
	"est-proxy/src/errors"
	"io"
	"log"
	"net/http"
)

type GptApi struct {
	httpClient *http.Client
}

func NewGptApi(httpClient *http.Client) *GptApi {
	return &GptApi{httpClient}
}

type Payload struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string          `json:"role"`
	Content []ContentObject `json:"content,omitempty"`
}

type ContentObject struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageUrl Image  `json:"image_url,omitempty"`
}

type Image struct {
	Url    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

type Response struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message RespMessage `json:"message"`
}

type RespMessage struct {
	Content string `json:"content"`
}

func (g GptApi) Request(
	prompt string,
	image []byte,
	ctx context.Context,
) (string, *errors.StatusError) {

	payload := toPayload(prompt, image)

	body, err := json.Marshal(payload)
	if err != nil {
		return "", errors.NewStatusError(http.StatusInternalServerError, "Failed to marshal payload")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", config.GPT_API_PATH, bytes.NewReader(body))
	if err != nil {
		return "", errors.NewStatusError(http.StatusInternalServerError, "Failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.GPT_API_TOKEN)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", errors.NewStatusError(http.StatusInternalServerError, "Failed to send request to GPT API")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("GPT API response: %s\n", string(body))
		return "", errors.NewStatusError(http.StatusInternalServerError, "Failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("GPT API response: %s\n", string(body))
		return "", errors.NewStatusError(resp.StatusCode, "Failed to request response from GPT API")
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("GPT API response: %s\n", string(body))
		return "", errors.NewStatusError(http.StatusInternalServerError, "Failed to unmarshal response body")
	}

	return response.Choices[0].Message.Content, nil
}

func toPayload(prompt string, image []byte) Payload {
	preContent := ContentObject{
		Type: "text",
		Text: prompt,
	}

	imageContent := ContentObject{
		Type: "image_url",
		ImageUrl: Image{
			Url:    "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(image),
			Detail: "high",
		},
	}

	return Payload{
		Model: "gpt-4o",
		Messages: []Message{
			{
				Role:    "user",
				Content: []ContentObject{preContent, imageContent},
			},
		},
	}
}
