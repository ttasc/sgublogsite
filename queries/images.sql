-- name: GetImageByID :one
SELECT * FROM images
WHERE image_id = ?;

-- name: GetAllImages :many
SELECT * FROM images;

-- name: AddImage :execresult
INSERT INTO images (url, name)
VALUES (?, ?);

-- name: DeleteImage :execresult
DELETE FROM images
WHERE image_id = ?;
