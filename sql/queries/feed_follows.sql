-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows (
    id,
    created_at,
    updated_at,
    user_id,
    feed_id
  )
  VALUES ($1, $2, $3, $4, $5)
  RETURNING *
)

SELECT
  iff.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM inserted_feed_follow iff
INNER JOIN feeds
ON feeds.id = iff.feed_id
INNER JOIN users
ON users.id = iff.user_id;

-- name: GetFeedFollowsForUser :many
SELECT
  ff.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM feed_follows ff
INNER JOIN feeds
ON feeds.id = ff.feed_id
INNER JOIN users
ON users.id = ff.user_id
WHERE ff.user_id = $1;
