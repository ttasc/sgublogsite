-- name: GetTagByID :one
SELECT *
FROM tags
WHERE tag_id = ?;

-- name: GetAllTags :many
SELECT * FROM tags;

-- name: AddTag :execresult
INSERT INTO tags (
    name,
    slug
) VALUES (?, ?);

-- name: DeleteTag :execresult
DELETE FROM tags
WHERE tag_id = ?;
