package impl_test

import (
	"context"
	"net/http"
	"testing"

	"est-proxy/src/errors"
	"est-proxy/src/service/impl"
	"est_proxy_go/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPreviewApi struct {
	mock.Mock
}

func (m *MockPreviewApi) GetPreview(
	boardId string,
	width, height,
	xLeft, yTop,
	xRight, yBottom float32,
	ctx context.Context,
) ([]byte, *errors.StatusError) {
	args := m.Called(boardId, width, height, xLeft, yTop, xRight, yBottom, ctx)
	return args.Get(0).([]byte), args.Get(1).(*errors.StatusError)
}

func (m *MockPreviewApi) GetTokens(boardIds []string, ctx context.Context) (map[string]string, *errors.StatusError) {
	return nil, nil
}

// MockGptApi реализация
type MockGptApi struct {
	mock.Mock
}

func (m *MockGptApi) Request(prompt string, image []byte, ctx context.Context) (string, *errors.StatusError) {
	args := m.Called(prompt, image, ctx)
	return args.String(0), args.Get(1).(*errors.StatusError)
}

func TestGptServiceImpl_Request(t *testing.T) {
	validRequestDto := models.GptRequestDto{
		BoardId:   "test-board",
		LeftUp:    models.Point{X: 10, Y: 20},
		RightDown: models.Point{X: 30, Y: 40},
		Prompt:    "test-prompt",
	}

	t.Run("successful request", func(t *testing.T) {
		previewMock := new(MockPreviewApi)
		gptMock := new(MockGptApi)

		// Настройка моков
		previewMock.On("GetPreview",
			"test-board",
			float32(1200), float32(800),
			float32(10), float32(20),
			float32(30), float32(40),
			mock.Anything,
		).Return([]byte("image-data"), (*errors.StatusError)(nil))

		gptMock.On("Request",
			"test-prompt",
			[]byte("image-data"),
			mock.Anything,
		).Return("test-response", (*errors.StatusError)(nil))

		service := impl.NewGptServiceImpl(previewMock, gptMock)
		req, _ := http.NewRequest("GET", "/", nil)

		response, err := service.Request(validRequestDto, req)

		assert.Nil(t, err)
		assert.Equal(t, "test-response", response.Text)
		previewMock.AssertExpectations(t)
		gptMock.AssertExpectations(t)
	})

	t.Run("preview api returns error", func(t *testing.T) {
		previewMock := new(MockPreviewApi)
		gptMock := new(MockGptApi)

		expectedErr := errors.NewStatusError(http.StatusInternalServerError, "preview error")

		previewMock.On("GetPreview", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]byte{}, expectedErr)

		service := impl.NewGptServiceImpl(previewMock, gptMock)
		req, _ := http.NewRequest("GET", "/", nil)

		response, err := service.Request(validRequestDto, req)

		assert.NotNil(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, response.Text)
		gptMock.AssertNotCalled(t, "Request")
	})

	t.Run("gpt api returns error", func(t *testing.T) {
		previewMock := new(MockPreviewApi)
		gptMock := new(MockGptApi)

		expectedErr := errors.NewStatusError(http.StatusBadGateway, "gpt error")

		previewMock.On("GetPreview", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]byte("image-data"), (*errors.StatusError)(nil))

		gptMock.On("Request", mock.Anything, mock.Anything, mock.Anything).
			Return("", expectedErr)

		service := impl.NewGptServiceImpl(previewMock, gptMock)
		req, _ := http.NewRequest("GET", "/", nil)

		response, err := service.Request(validRequestDto, req)

		assert.NotNil(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, response.Text)
	})

	t.Run("context propagation", func(t *testing.T) {
		previewMock := new(MockPreviewApi)
		gptMock := new(MockGptApi)
		ctx := context.WithValue(context.Background(), "test-key", "test-value")

		previewMock.On("GetPreview", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.MatchedBy(func(c context.Context) bool {
				return c.Value("test-key") == "test-value"
			})).
			Return([]byte("image-data"), (*errors.StatusError)(nil))

		gptMock.On("Request", mock.Anything, mock.Anything, mock.MatchedBy(func(c context.Context) bool {
			return c.Value("test-key") == "test-value"
		})).
			Return("response", (*errors.StatusError)(nil))

		service := impl.NewGptServiceImpl(previewMock, gptMock)
		req, _ := http.NewRequest("GET", "/", nil)
		req = req.WithContext(ctx)

		_, err := service.Request(validRequestDto, req)

		assert.Nil(t, err)
		previewMock.AssertExpectations(t)
		gptMock.AssertExpectations(t)
	})
}
