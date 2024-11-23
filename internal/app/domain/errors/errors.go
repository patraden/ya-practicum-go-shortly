package errors

import (
	"errors"
	"fmt"
)

var (
	ErrStateCreate            = errors.New("[memento] could not create state")
	ErrStateRestore           = errors.New("[memento] could not restore state")
	ErrPGEmptyPool            = errors.New("[postgres] empty database conn pool")
	ErrOriginalExists         = errors.New("[repository] original url exists")
	ErrSlugExists             = errors.New("[repository] slug exists")
	ErrSlugNotFound           = errors.New("[repository] slug not found")
	ErrUserNotFound           = errors.New("[repository] user not found")
	ErrSlugInvalid            = errors.New("[shortener] invalid slug")
	ErrSlugCollision          = errors.New("[shortener] slug collision")
	ErrShortenerInternal      = errors.New("[shortener] internal error")
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

func Wrap(msg string, err error, label string) error {
	return fmt.Errorf("[%s] %s: %w", label, msg, err)
}
