-- name: GetURLMapping :one
SELECT slug, original, user_id, created_at, expires_at, deleted
FROM shortener.urlmapping
WHERE slug = $1;

-- name: GetUserURLMappings :many
SELECT slug, original, user_id, created_at, expires_at, deleted
FROM shortener.urlmapping
WHERE user_id =$1;

-- name: AddURLMapping :one
INSERT INTO shortener.urlmapping (slug, original, user_id, created_at, expires_at, deleted)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (original) DO UPDATE
SET slug = shortener.urlmapping.slug,
    user_id = shortener.urlmapping.user_id,
    created_at = shortener.urlmapping.created_at,
    expires_at = shortener.urlmapping.expires_at,
    deleted = shortener.urlmapping.deleted
RETURNING slug, original, user_id, created_at, expires_at, deleted;

-- name: AddURLMappingBatchCopy :copyfrom
INSERT INTO shortener.urlmapping (slug, original, user_id, created_at, expires_at, deleted)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateDeletedSlugTempTable :exec
CREATE TEMP TABLE urlmapping_tmp (
    slug    VARCHAR(8)  PRIMARY KEY,
    user_id UUID        NOT NULL
) ON COMMIT DROP;

-- name: FillDeletedSlugTempTable :copyfrom
INSERT INTO urlmapping_tmp (slug, user_id)
VALUES ($1, $2);

-- name: DeleteSlugsInTarget :exec
UPDATE shortener.urlmapping
SET deleted = true
FROM urlmapping_tmp
WHERE shortener.urlmapping.slug = urlmapping_tmp.slug
  AND shortener.urlmapping.user_id = urlmapping_tmp.user_id;

-- name: GetStats :one
SELECT 
  COUNT(1)::BIGINT AS CountSlugs,
  COUNT(DISTINCT user_id)::BIGINT AS CountUsers
FROM shortener.urlmapping;