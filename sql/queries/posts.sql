-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;


-- name: GetPostsForUser :many
SELECT posts.title, posts.url, posts.description, posts.published_at
FROM posts
LEFT JOIN feeds on posts.feed_id = feeds.id
LEFT JOIN users on feeds.user_id = users.id AND users.id = $1
ORDER BY posts.published_at DESC
fetch first $2 row only;
