-- name: CreateFeedFollow :many
WITH inserted_ff AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
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
    inserted_ff.*,
    feeds.name AS feed_nme,
    users.name AS user_name
FROM inserted_ff
INNER JOIN feeds ON inserted_ff.feed_id = feeds.id
INNER JOIN users ON inserted_ff.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*,
feeds.name AS feed_name,
users.name AS user_name
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE users.name = $1;

-- name: Unfollow :exec
DELETE FROM feed_follows 
WHERE user_id = $1 AND
feed_id = $2;