-- name: GetUserByID :one
SELECT users.*, images.url AS avatar
FROM users LEFT JOIN images ON users.avatar_id = images.image_id
WHERE user_id = ?;

-- name: GetUserByEmailOrPhone :one
SELECT users.*, images.url AS avatar
FROM users LEFT JOIN images ON users.avatar_id = images.image_id
WHERE email = ? OR phone = ?;

-- name: FindUsers :many
SELECT users.*, images.url AS avatar
FROM users LEFT JOIN images ON users.avatar_id = images.image_id
WHERE lower(concat(firstname, ' ', lastname, ' ', phone, ' ', email)) LIKE lower(sqlc.arg(text));
-- WHERE MATCH(firstname, lastname, phone, email)) AGAINST (sqlc.arg(text));

-- name: GetAllUsers :many
SELECT users.*, images.url AS avatar
FROM users LEFT JOIN images ON users.avatar_id = images.image_id
ORDER BY lastname;

-- name: GetUserAvatarID :one
SELECT avatar_id
FROM users
WHERE user_id = ?;

-- name: AddUser :execresult
INSERT INTO users (
    firstname,
    lastname,
    phone,
    email,
    password,
    role
) VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateUserInfo :execresult
UPDATE users
SET
    firstname = ?,
    lastname = ?,
    phone = ?,
    email = ?
WHERE user_id = ?;

-- name: UpdateUserAvatar :execresult
UPDATE users
SET avatar_id = ?
WHERE user_id = ?;

-- name: UpdateUserPassword :execresult
UPDATE users
SET password = ?
WHERE user_id = ?;

-- name: UpdateUserRole :execresult
UPDATE users
SET role = ?
WHERE user_id = ?;

-- name: DeleteUser :execresult
DELETE FROM users
WHERE user_id = ?;
