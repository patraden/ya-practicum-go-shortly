package errors

import (
	"errors"
	"fmt"
)

// Static errors.
var (
	ErrStateNotmplemented     = errors.New("[memento] state store/restore not implemented")
	ErrPGEmptyPool            = errors.New("[postgres] empty database conn pool")
	ErrOriginalExists         = errors.New("[repository] original url exists")
	ErrSlugExists             = errors.New("[repository] slug exists")
	ErrSlugNotFound           = errors.New("[repository] slug not found")
	ErrUserNotFound           = errors.New("[repository] user not found")
	ErrMissedJob              = errors.New("[batcher] missed output job")
	ErrMissedTask             = errors.New("[batcher] missed input task")
	ErrFailedCast             = errors.New("[batcher] failed to cast")
	ErrSlugInvalid            = errors.New("[shortener] invalid slug")
	ErrSlugDeleted            = errors.New("[shortener] slug deleted")
	ErrSlugCollision          = errors.New("[shortener] slug collision")
	ErrShortenerInternal      = errors.New("[shortener] internal error")
	ErrStatsProviderInternal  = errors.New("[statsprovider] internal error")
	ErrRemoverInternal        = errors.New("[remover] internal error")
	ErrRemoverInitBatcher     = errors.New("[remover] init batcher error")
	ErrInvalidConfig          = errors.New("[config] bad config parameters")
	ErrEnvConfigParse         = errors.New("[config] env vars parsing error")
	ErrUtilsCompEncoding      = errors.New("[utils] bad compression encoding")
	ErrUtilsDecompionEncoding = errors.New("[utils] bad decompression encoding")
	ErrUtilsEncoderOpen       = errors.New("[utils] compression encoder open error")
	ErrUtilsEncoderCast       = errors.New("[utils] compression encoder cast error")
	ErrURLGenGenerateSlug     = errors.New("[urlgenerator] slug(s) generation error")
	ErrAuthInvalidToken       = errors.New("[middleware] invalid jwt token")
	ErrAuthUnexpectedSign     = errors.New("[middleware] unexpected sign method")
	ErrAuthNoCookie           = errors.New("[middleware] no auth cookie")
	ErrTestGeneral            = errors.New("[test] test error")
)

// Wrap formats and wraps an error message with a specific label.
func Wrap(msg string, err error, label string) error {
	return fmt.Errorf("[%s] %s: %w", label, msg, err)
}
