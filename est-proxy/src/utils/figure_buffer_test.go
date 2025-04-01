package utils

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAddAndRemove(t *testing.T) {
	fb := NewFigureBuffer()
	figureID := "test-1"
	testData := []byte{1, 2, 3}

	fb.Add(figureID, testData)
	data, exists := fb.data.Get(figureID)
	require.True(t, exists)
	require.Equal(t, testData, data.(figureUpdateData).data)

	fb.Add(figureID, []byte{4, 5})
	data, _ = fb.data.Get(figureID)
	require.Equal(t, []byte{1, 2, 3, 4, 5}, data.(figureUpdateData).data)

	require.True(t, fb.Remove(figureID))
	_, exists = fb.data.Get(figureID)
	require.False(t, exists)

	require.False(t, fb.Remove("non-existent"))
}

func TestServeFlush(t *testing.T) {
	originalExpiration := bufferedFigureExpirationTime
	bufferedFigureExpirationTime = time.Millisecond
	defer func() { bufferedFigureExpirationTime = originalExpiration }()

	fb := NewFigureBuffer()
	figureID := "test-1"

	var flushedData []byte
	var flushedID string
	callback := func(id string, data []byte) {
		flushedID = id
		flushedData = data
	}

	go func() {
		fb.ServeFlush(callback)
	}()

	fb.Add(figureID, []byte{1, 2, 3})

	require.Eventually(t, func() bool {
		return flushedID == figureID && bytes.Equal(flushedData, []byte{1, 2, 3})
	}, 5*time.Millisecond, 1*time.Millisecond)

	_, exists := fb.data.Get(figureID)
	require.False(t, exists)
}
