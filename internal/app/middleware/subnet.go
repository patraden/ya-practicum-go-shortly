package middleware

import (
	"net"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
)

// SubnetMiddleware is a middleware handler that verifies that request has been received from a trusted subnet.
func SubnetMiddleware(log *zerolog.Logger, config *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			if config.TrustedSubnet == "" {
				next.ServeHTTP(response, request)
				return
			}

			ipStr := request.Header.Get("X-Real-IP")

			if ipStr == "" {
				log.Info().
					Str("subnet", config.TrustedSubnet).
					Msg("X-Real-IP header is nil")

				response.WriteHeader(http.StatusForbidden)

				return
			}

			_, trustedNet, err := net.ParseCIDR(config.TrustedSubnet)
			if err != nil {
				log.Error().Err(err).
					Str("subnet", config.TrustedSubnet).
					Msg("invalid subnet")

				response.WriteHeader(http.StatusInternalServerError)

				return
			}

			ip := net.ParseIP(ipStr)
			if !trustedNet.Contains(ip) {
				log.Info().
					Str("subnet", config.TrustedSubnet).
					Str("ip", ip.String()).
					Msg("ip does not belong to subnet")

				response.WriteHeader(http.StatusForbidden)

				return
			}

			next.ServeHTTP(response, request)
		})
	}
}
