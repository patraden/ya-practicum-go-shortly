package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
)

func TestSubnetMiddleware(t *testing.T) {
	t.Parallel()

	logger := logger.NewLogger(zerolog.DebugLevel).GetLogger()

	// Create a handler that just returns OK
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		ipHeader       string
		expectedStatus int
		config         *config.Config
	}{
		{"Allowed IP", "192.168.1.100", http.StatusOK, &config.Config{TrustedSubnet: "192.168.1.0/24"}},
		{"Forbidden IP", "10.0.0.1", http.StatusForbidden, &config.Config{TrustedSubnet: "192.168.1.0/24"}},
		{"Missing IP Header", "", http.StatusForbidden, &config.Config{TrustedSubnet: "192.168.1.0/24"}},
		{"Broken subnet", "192.168.1.100", http.StatusInternalServerError, &config.Config{TrustedSubnet: "111"}},
		{"No subnet", "192.168.1.100", http.StatusOK, &config.Config{TrustedSubnet: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			middleware := middleware.SubnetMiddleware(logger, tt.config)

			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			if tt.ipHeader != "" {
				req.Header.Set("X-Real-IP", tt.ipHeader)
			}

			recorder := httptest.NewRecorder()
			middleware(handler).ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
		})
	}
}
