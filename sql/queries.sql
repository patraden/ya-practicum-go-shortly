-- name: GetURLMapping :one
SELECT slug, original, user_id, created_at, expires_at
FROM shortener.urlmapping
WHERE slug = $1;

-- name: GetUserURLMappings :many
SELECT slug, original, user_id, created_at, expires_at
FROM shortener.urlmapping
WHERE user_id =$1;

-- name: AddURLMapping :one
INSERT INTO shortener.urlmapping (slug, original, user_id, created_at, expires_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (original) DO UPDATE
SET slug = shortener.urlmapping.slug,
    user_id = shortener.urlmapping.user_id,
    created_at = shortener.urlmapping.created_at,
    expires_at = shortener.urlmapping.expires_at
RETURNING slug, original, user_id, created_at, expires_at;

-- name: AddURLMappingBatchCopy :copyfrom
INSERT INTO shortener.urlmapping (slug, original, user_id, created_at, expires_at)
VALUES ($1, $2, $3, $4, $5);