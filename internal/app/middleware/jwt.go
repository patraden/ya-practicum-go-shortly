package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
)

// Aux constants.
const (
	UserIDKey            contextKey = "user_id"
	AuthCookieName                  = "auth_token"
	defaultTokenDuration            = 365 * 24 * time.Hour
)

// Authenticate is a middleware handler that checks if the request has a valid JWT token.
// If the token is not valid, it generates a new one and sets it in the response cookie.
func Authenticate(log *zerolog.Logger, config *config.Config) func(next http.Handler) http.Handler {
	auth := NewJWTMiddleware(
		func(*jwt.Token) (interface{}, error) { return []byte(config.JWTSecret), nil },
		log,
	)

	return auth.AuthenticateHandler
}

// Authorize is a middleware handler that checks if the request is authorized with a valid JWT token.
// It ensures that the user is authenticated before allowing access to the resource.
func Authorize(log *zerolog.Logger, config *config.Config) func(next http.Handler) http.Handler {
	auth := NewJWTMiddleware(
		func(*jwt.Token) (interface{}, error) { return []byte(config.JWTSecret), nil },
		log,
	)

	return auth.AuthorizeHandler
}

// Aux types.
type (
	contextKey     string
	TokenExtractor func(r *http.Request) (string, error)
)

// Claims represents the JWT claims for a user. It includes the user ID and standard JWT claims.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// JWTMiddleware is a struct that provides JWT-based authentication and authorization middleware.
type JWTMiddleware struct {
	keyFunc   jwt.Keyfunc
	log       *zerolog.Logger
	extractor TokenExtractor
}

// CookieTokenExtractor extracts the JWT token from a cookie in the HTTP request.
func CookieTokenExtractor(r *http.Request) (string, error) {
	cookie, err := r.Cookie(AuthCookieName)

	if errors.Is(err, http.ErrNoCookie) || cookie.Value == `` {
		return ``, e.ErrAuthNoCookie
	}

	if err != nil {
		return ``, e.Wrap(`failed to extract cookie`, err, errLabel)
	}

	return cookie.Value, nil
}

// NewJWTMiddleware creates a new JWTMiddleware instance with the provided key function and logger.
func NewJWTMiddleware(keyFunc jwt.Keyfunc, log *zerolog.Logger) *JWTMiddleware {
	// add signature method validation
	kfunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, e.ErrAuthUnexpectedSign
		}

		return keyFunc(t)
	}

	return &JWTMiddleware{
		keyFunc:   kfunc,
		log:       log,
		extractor: CookieTokenExtractor,
	}
}

// ValidateToken validates the JWT token string and returns the claims if valid.
func (auth *JWTMiddleware) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, auth.keyFunc)
	if err != nil {
		msg := "failed to parse JWT token"
		auth.log.Error().Err(err).
			Str("token_string", tokenString).
			Msg(msg)

		return claims, e.Wrap(msg, err, errLabel)
	}

	if !token.Valid {
		auth.log.Error().
			Str("token_string", token.Raw).
			Str("user_id", claims.UserID).
			Msg("invalid token")

		return claims, e.ErrAuthInvalidToken
	}

	if userID, err := domain.ParseUserID(claims.UserID); err != nil || userID.IsNil() {
		auth.log.Error().
			Str("token_string", token.Raw).
			Msg("missing user_id in token")

		return claims, e.ErrAuthInvalidToken
	}

	return claims, nil
}

// GenerateToken generates a new JWT token for the provided user ID.
func (auth *JWTMiddleware) GenerateToken(userID domain.UserID) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(defaultTokenDuration)),
		},
	}

	auth.log.Info().
		Str("userID", userID.String()).
		Msg("generating token")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingKey, err := auth.keyFunc(token)
	if err != nil {
		msg := `failed to retrieve signing key`
		auth.log.Error().Err(err).Msg(msg)

		return ``, e.Wrap(msg, err, errLabel)
	}

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		msg := `failed to sign token`
		auth.log.Error().Err(err).Msg(msg)

		return ``, e.Wrap(msg, err, errLabel)
	}

	return tokenString, nil
}

func (auth *JWTMiddleware) extractAndValidateToken(r *http.Request) (domain.UserID, string, error) {
	token, err := auth.extractor(r)
	if err != nil {
		auth.log.Error().Err(err).Msg(`failed to extract token`)

		return domain.UserID{}, ``, err
	}

	claims, err := auth.ValidateToken(token)
	if err != nil {
		auth.log.Error().Err(err).Msg(`failed to validate token`)

		return domain.UserID{}, ``, err
	}

	userID, err := domain.ParseUserID(claims.UserID)
	if err != nil {
		auth.log.Error().Err(err).Msg(`failed to parse user from token claims`)

		return domain.UserID{}, ``, e.Wrap("failed to parse user from token claims", err, errLabel)
	}

	return userID, token, nil
}

// AuthenticateHandler is the handler for the Authenticate middleware.
// It validates or generates a token and adds the user ID to the request context.
func (auth *JWTMiddleware) AuthenticateHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, token, err := auth.extractAndValidateToken(r)
		if err != nil {
			if userID.IsNil() {
				userID = domain.NewUserID()
			}

			token, err = auth.GenerateToken(userID)
			if err != nil {
				http.Error(w, "failed to generate token", http.StatusInternalServerError)

				return
			}
		}

		r = r.Clone(context.WithValue(r.Context(), UserIDKey, userID))

		http.SetCookie(w, &http.Cookie{
			Name:     AuthCookieName,
			HttpOnly: true,
			Value:    token,
		})

		next.ServeHTTP(w, r)
	})
}

// AuthorizeHandler is the handler for the Authorize middleware.
// It ensures that the request is authorized with a valid token before proceeding.
func (auth *JWTMiddleware) AuthorizeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _, err := auth.extractAndValidateToken(r)
		if err != nil {
			http.Error(w, `unauthorized`, http.StatusUnauthorized)

			return
		}

		auth.log.Info().Str("user_id", userID.String()).Msg("user_id authorized")

		r = r.Clone(context.WithValue(r.Context(), UserIDKey, userID))

		next.ServeHTTP(w, r)
	})
}

// GetUserID extracts the user ID from the request context.
func GetUserID(ctx context.Context) (domain.UserID, bool) {
	userID, ok := ctx.Value(UserIDKey).(domain.UserID)

	return userID, ok
}
