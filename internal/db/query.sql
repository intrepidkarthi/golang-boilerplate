-- name: CreateMessage :one
INSERT INTO messages (
    content
) VALUES (
    $1
)
RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1;

-- name: ListMessages :many
SELECT * FROM messages
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateMessage :one
UPDATE messages
SET content = $2
WHERE id = $1
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages
WHERE id = $1;

-- name: GetTotalMessages :one
SELECT COUNT(*) FROM messages;
