-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, user_id,name,url)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name AS feedName, feeds.url, users.name AS userName 
FROM feeds 
INNER JOIN users 
ON feeds.user_id = users.id;

-- name: GetFeed :one
SELECT * FROM feeds 
WHERE url = $1;