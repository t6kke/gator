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
JOIN feed_follows on posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
fetch first $2 row only;
