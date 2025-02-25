-- name: CreateMessage :one
INSERT INTO messages (content)
VALUES ($1)
RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateMessage :one
UPDATE messages
SET content = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteMessage :exec
UPDATE messages
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListMessages :many
SELECT * FROM messages
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTotalMessages :one
SELECT COUNT(*) FROM messages
WHERE deleted_at IS NULL;
