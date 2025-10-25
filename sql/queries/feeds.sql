-- name: CreateFeed :one
INSERT INTO feeds (
  id,
  created_at,
  updated_at,
  name,
  url,
  user_id
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetAllFeeds :many
SELECT
  COALESCE(feeds.name, '') AS name,
  COALESCE(feeds.url, '') AS url,
  users.name AS user
FROM feeds
RIGHT JOIN users
ON users.id = feeds.user_id;
