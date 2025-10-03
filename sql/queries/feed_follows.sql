-- name: CreateFeedFollow :one
WITH cff AS
(
    INSERT INTO feed_follows (id, created_at, updated_at, user_id,feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT 
    cff.id AS ffID,
    cff.created_at,
    cff.updated_at,
    cff.user_id AS userID,
    cff.feed_id AS feedID,
    users.name AS userName,
    feeds.name AS feedName
FROM 
    cff
JOIN 
    feeds ON feeds.id = cff.feed_id
JOIN 
    users ON users.id = cff.user_id;

-- name: GetFeedFollowsForUser :many
SELECT
    feeds.name AS feedName,
    users.name AS userName
FROM
    feed_follows
JOIN
    feeds ON feeds.id = feed_follows.feed_id
JOIN
    users ON users.id = feed_follows.user_id
WHERE
    users.name = $1;

-- name: UnFollowFeed :exec

DELETE FROM
    feed_follows
USING 
    feeds
WHERE
    feeds.url = $1 AND feed_follows.user_id = $2;


