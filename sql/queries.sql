-- name: GetURLMapping :one
SELECT slug, original, created_at, expires_at
FROM shortener.urlmapping
WHERE slug = $1;

-- name: AddURLMapping :one
INSERT INTO shortener.urlmapping (slug, original, created_at, expires_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (original) DO UPDATE
SET slug = shortener.urlmapping.slug,
    created_at = shortener.urlmapping.created_at,
    expires_at = shortener.urlmapping.expires_at
RETURNING slug, original, created_at, expires_at;

-- name: AddURLMappingBatchCopy :copyfrom
INSERT INTO shortener.urlmapping (slug, original, created_at, expires_at)
VALUES ($1, $2, $3, $4);