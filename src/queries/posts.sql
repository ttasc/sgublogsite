-- name: GetAllPosts :many
SELECT *
FROM posts
ORDER BY created_at DESC;

-- name: GetPostsByUserID :many
SELECT *
FROM posts
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: GetPostsByCategoryID :many
SELECT *
FROM posts
WHERE post_id IN (
    SELECT post_id
    FROM post_categories
    WHERE category_id = ?
)
ORDER BY created_at DESC;

-- name: GetUncategorizedPosts :many
SELECT *
FROM posts
WHERE post_id NOT IN (
    SELECT post_id
    FROM post_categories
)
ORDER BY created_at DESC;

-- name: GetPostsByCategoryName :many
SELECT *
FROM posts
WHERE post_id IN (
    SELECT post_id
    FROM post_categories
    WHERE category_id IN (
        SELECT category_id
        FROM categories
        WHERE name = ?
    )
)
ORDER BY created_at DESC;

-- name: GetPostsByTagID :many
SELECT *
FROM posts
WHERE post_id IN (
    SELECT post_id
    FROM post_tags
    WHERE tag_id = ?
)
ORDER BY created_at DESC;

-- name: GetPostsByTagName :many
SELECT *
FROM posts
WHERE post_id IN (
    SELECT post_id
    FROM post_tags
    WHERE tag_id IN (
        SELECT tag_id
        FROM tags
        WHERE name = ?
    )
)
ORDER BY created_at DESC;

-- name: GetPostsByStatus :many
SELECT *
FROM posts
WHERE status = ?
ORDER BY created_at DESC;

-- name: FindPosts :many
SELECT *
FROM posts
WHERE lower(concat(title, ' ', body)) LIKE lower(sqlc.arg(text));
-- WHERE MATCH(title, body)) AGAINST (sqlc.arg(text));

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
UPDATE posts
SET
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
