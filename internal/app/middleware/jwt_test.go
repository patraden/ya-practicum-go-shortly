package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
)

func TestAuthenticateTokenIsMissing(t *testing.T) {
	t.Parallel()

	log := zerolog.New(nil).With().Logger()
	cfg := &config.Config{JWTSecret: "test-secret"}
	authenticate := middleware.Authenticate(&log, cfg)

	handler := authenticate(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		_, ok := middleware.GetUserID(r.Context())
		assert.True(t, ok, "UserID should be present in context")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK")
	cookie := resp.Cookies()[0]
	assert.Equal(t, "auth_token", cookie.Name, "Expected auth cookie")
	assert.NotEmpty(t, cookie.Value, "Auth cookie should have a value")
}

func TestAuthenticateTokenIsPresent(t *testing.T) {
	t.Parallel()

	log := zerolog.New(nil).With().Logger()
	cfg := &config.Config{JWTSecret: "test-secret"}
	authMiddleware := middleware.NewJWTMiddleware(
		func(*jwt.Token) (interface{}, error) { return []byte(cfg.JWTSecret), nil },
		&log,
	)
	authenticate := authMiddleware.AuthenticateHandler

	userID, err := domain.ParseUserID("f3a99c97-7f28-4a16-b020-9b82cfb9883b")
	require.NoError(t, err)

	token, err := authMiddleware.GenerateToken(userID)
	require.NoError(t, err, "Token generation should succeed")

	handler := authenticate(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		ctxUserID, ok := middleware.GetUserID(r.Context())
		assert.True(t, ok, "UserID should be present in context")
		assert.Equal(t, userID.String(), ctxUserID.String(), "UserID in context should match the token")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  middleware.AuthCookieName,
		Value: token,
	})

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK")

	var cookie *http.Cookie

	for _, c := range resp.Cookies() {
		if c.Name == middleware.AuthCookieName {
			cookie = c

			break
		}
	}

	assert.NotNil(t, cookie, "Expected auth cookie to be present")
	assert.Equal(t, middleware.AuthCookieName, cookie.Name, "Expected auth cookie")
	assert.Equal(t, token, cookie.Value, "Auth cookie value should match the original token")
}

func TestAuthorizeTokenIsMissing(t *testing.T) {
	t.Parallel()

	log := zerolog.New(nil).With().Logger()
	cfg := &config.Config{JWTSecret: "test-secret"}
	authorize := middleware.Authorize(&log, cfg)

	handler := authorize(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("handler should not be reached, unauthorized request")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Expected status Unauthorized")
}

func TestAuthorizeTokenIsInvalid(t *testing.T) {
	t.Parallel()

	log := zerolog.New(nil).With().Logger()
	cfg := &config.Config{JWTSecret: "test-secret"}
	authorize := middleware.Authorize(&log, cfg)

	invalidToken := "invalid-token"
	handler := authorize(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("handler should not be reached, unauthorized request")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  middleware.AuthCookieName,
		Value: invalidToken,
	})

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Expected status Unauthorized")
}

func TestAuthorizeTokenIsValid(t *testing.T) {
	t.Parallel()

	log := zerolog.New(nil).With().Logger()
	cfg := config.DefaultConfig()
	authorize := middleware.Authorize(&log, cfg)
	authMiddleware := middleware.NewJWTMiddleware(
		func(*jwt.Token) (interface{}, error) { return []byte(cfg.JWTSecret), nil },
		&log,
	)
	userID, err := domain.ParseUserID("f3a99c97-7f28-4a16-b020-9b82cfb9883b")
	require.NoError(t, err)

	token, err := authMiddleware.GenerateToken(userID)
	require.NoError(t, err, "Token generation should succeed")

	handler := authorize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxUserID, ok := middleware.GetUserID(r.Context())
		assert.True(t, ok, "UserID should be present in context")
		assert.Equal(t, userID.String(), ctxUserID.String(), "UserID in context should match the token")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  middleware.AuthCookieName,
		Value: token,
	})

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK")
}

func TestAuthorizeTokenHasExpired(t *testing.T) {
	t.Parallel()

	log := zerolog.New(nil).With().Logger()
	cfg := &config.Config{JWTSecret: "test-secret"}
	authorize := middleware.Authorize(&log, cfg)
	userID, err := domain.ParseUserID("f3a99c97-7f28-4a16-b020-9b82cfb9883b")
	require.NoError(t, err)

	expiredClaims := &middleware.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredToken, err := token.SignedString([]byte(cfg.JWTSecret))
	require.NoError(t, err, "Token signing should succeed")

	handler := authorize(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("handler should not be reached, unauthorized request")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  middleware.AuthCookieName,
		Value: expiredToken,
	})

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Expected status Unauthorized")
}
