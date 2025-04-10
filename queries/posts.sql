-- name: GetPostByID :one
SELECT posts.*, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE post_id = ?;

-- name: CountPosts :one
SELECT COUNT(*) FROM posts;

-- name: GetFilteredPosts :many
SELECT
    p.post_id,
    p.title,
    CONCAT(u.firstname, ' ', u.lastname) AS author_name,
    GROUP_CONCAT(c.name SEPARATOR ', ') AS categories,
    p.created_at,
    p.status
FROM posts p
LEFT JOIN users u ON p.user_id = u.user_id
LEFT JOIN post_categories pc ON p.post_id = pc.post_id
LEFT JOIN categories c ON pc.category_id = c.category_id
WHERE
    (CASE WHEN sqlc.arg(title)      != ''       THEN lower(p.title) LIKE lower(sqlc.arg(title)) ELSE TRUE END) AND
    (CASE WHEN sqlc.arg(status)     != ''       THEN p.status       =  sqlc.arg(status) ELSE TRUE END) AND
    (CASE WHEN sqlc.arg(private)    IS NOT NULL THEN p.private      =  sqlc.arg(private) ELSE TRUE END)
GROUP BY p.post_id, p.title, author_name, p.created_at, p.status
ORDER BY p.created_at DESC
LIMIT ? OFFSET ?;

-- name: GetPosts :many
SELECT posts.*, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetPostsByUserID :many
SELECT posts.*, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE user_id = ?
AND status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetPostsByCategoryID :many
SELECT posts.*, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE post_id IN (
    SELECT post_id
    FROM post_categories
    WHERE category_id = ?
)
AND (private = ? OR private = 0)
AND status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetUncategorizedPosts :many
SELECT posts.*, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE post_id NOT IN (
    SELECT post_id
    FROM post_categories
)
AND (private = ? OR private = 0)
AND status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetPostsByCategorySlug :many
SELECT posts.*, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE post_id IN (
    SELECT post_id
    FROM post_categories
    WHERE category_id IN (
        SELECT category_id
        FROM categories
        WHERE categories.slug = ?
    )
)
AND (private = ? OR private = 0)
AND status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetPostsByTagID :many
SELECT posts.*, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE post_id IN (
    SELECT post_id
    FROM post_tags
    WHERE tag_id = ?
)
AND (private = ? OR private = 0)
AND status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetPostsByTagSlug :many
SELECT posts.*, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE post_id IN (
    SELECT post_id
    FROM post_tags
    WHERE tag_id IN (
        SELECT tag_id
        FROM tags
        WHERE tags.slug = ?
    )
)
AND (private = ? OR private = 0)
AND status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetPostsByStatus :many
SELECT posts.*, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: FindPosts :many
SELECT posts.*, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE lower(concat(title, ' ', body)) LIKE lower(sqlc.arg(text))
-- WHERE MATCH(title, body) AGAINST (sqlc.arg(text));
AND (private = ? OR private = 0)
AND status = ?
LIMIT ? OFFSET ?;

-- name: CreatePost :execresult
INSERT INTO posts (
    user_id,
    title,
    slug,
    thumbnail_id,
    body
) VALUES (?, ?, ?, ?, ?);

-- name: AddPostToCategory :execresult
INSERT INTO post_categories (
    post_id,
    category_id
) VALUES (?, ?);

-- name: AddTagToPost :execresult
INSERT INTO post_tags (
    post_id,
    tag_id
) VALUES (?, ?);

-- name: UpdatePostMetadata :execresult
UPDATE posts SET
    title = ?,
    slug = ?,
    thumbnail_id = ?
WHERE post_id = ?;

-- name: UpdatePostStatus :execresult
UPDATE posts
SET status = ?
WHERE post_id = ?;

-- name: UpdatePostPrivate :execresult
UPDATE posts
SET private = ?
WHERE post_id = ?;

-- name: UpdatePostBody :execresult
UPDATE posts
SET body = ?
WHERE post_id = ?;

-- name: DeletePost :execresult
DELETE FROM posts
WHERE post_id = ?;
