// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sql

import (
	"time"

	domain "github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
)

type ShortenerUrlmapping struct {
	Slug      domain.Slug        `db:"slug" json:"short_url"`
	Original  domain.OriginalURL `db:"original" json:"original_url"`
	CreatedAt time.Time          `db:"created_at" json:"created_at"`
	ExpiresAt time.Time          `db:"expires_at" json:"expires_at"`
}
