-- name: GetPostByID :one
SELECT posts.*, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE post_id = ?;

-- name: GetAllPosts :many
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
