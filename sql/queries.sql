-- name: GetURLMapping :one
SELECT slug, original, created_at, expires_at
FROM shortener.urlmapping
WHERE slug = $1;

-- name: AddURLMapping :exec
INSERT INTO shortener.urlmapping (slug, original, created_at, expires_at)
VALUES ($1, $2, $3, $4);