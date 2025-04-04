-- name: CountImages :one
SELECT COUNT(*) FROM images;

-- name: GetImageByID :one
SELECT * FROM images
WHERE image_id = ?;

-- name: GetImageByURL :one
SELECT * FROM images
WHERE url = ?;

-- name: GetAllImages :many
SELECT * FROM images
LIMIT ? OFFSET ?;

-- name: AddImage :execresult
INSERT INTO images (url, name)
VALUES (?, ?);

-- name: UpdateImageURL :execresult
UPDATE images
SET url = ?
WHERE image_id = ?;

-- name: DeleteImage :execresult
DELETE FROM images
WHERE image_id = ?;
