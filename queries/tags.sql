-- name: GetTagByID :one
SELECT *
FROM tags
WHERE tag_id = ?;

-- name: GetAllTags :many
SELECT * FROM tags;

-- name: GetAllTagNames :many
SELECT name FROM tags;

-- name: GetTagsByPostID :many
SELECT tags.* FROM tags
WHERE tag_id IN (
    SELECT tag_id
    FROM post_tags
    WHERE post_id = ?
);

-- name: AddTag :execresult
INSERT INTO tags (
    name,
    slug
) VALUES (?, ?);

-- name: UpdateTag :execresult
UPDATE tags
SET name = ?,
    slug = ?
WHERE tag_id = ?;

-- name: DeleteTag :execresult
DELETE FROM tags
WHERE tag_id = ?;
