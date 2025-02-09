// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sql

import (
	"time"

	domain "github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
)

type ShortenerUrlmapping struct {
	Slug      domain.Slug        `db:"slug"`
	Original  domain.OriginalURL `db:"original"`
	UserID    domain.UserID      `db:"user_id"`
	CreatedAt time.Time          `db:"created_at"`
	ExpiresAt time.Time          `db:"expires_at"`
	Deleted   bool               `db:"deleted"`
}

type UrlmappingTmp struct {
	Slug   domain.Slug   `db:"slug"`
	UserID domain.UserID `db:"user_id"`
}
