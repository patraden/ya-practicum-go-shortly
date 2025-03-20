package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
)

func TestLoggerMiddleware(t *testing.T) {
	t.Parallel()

	// Prepare a buffer to capture logs
	var logBuffer bytes.Buffer
	logger := zerolog.New(&logBuffer)

	// Simple handler that returns 200 OK with "Hello, World!"
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			t.Errorf("failed to wite payload %v", err)
		}
	})

	loggedHandler := middleware.Logger(&logger)(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	loggedHandler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "Hello, World!", rec.Body.String())

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, `"uri":"/test"`)
	assert.Contains(t, logOutput, `"method":"GET"`)
	assert.Contains(t, logOutput, `"status":200`)
	assert.Contains(t, logOutput, `"size":13`)  // "Hello, World!" is 13 bytes
	assert.Contains(t, logOutput, `"duration"`) // Ensure duration is logged
}
