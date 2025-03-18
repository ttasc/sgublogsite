-- name: GetUserByID :one
SELECT *
FROM users
WHERE user_id = ?;

-- name: GetUserByEmailOrMobile :one
SELECT *
FROM users
WHERE email = ? OR mobile = ?;

-- name: FindUsers :many
SELECT *
FROM users
WHERE lower(concat(firstname, ' ', lastname, ' ', mobile, ' ', email)) LIKE lower(sqlc.arg(text));
-- WHERE MATCH(firstname, lastname, mobile, email)) AGAINST (sqlc.arg(text));

-- name: GetAllUsers :many
SELECT *
FROM users
ORDER BY lastname;

-- name: AddUser :execresult
INSERT INTO users (
    firstname,
    lastname,
    mobile,
    email,
    password,
    role
) VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateUserInfo :execresult
UPDATE users
SET
    firstname = ?,
    lastname = ?,
    mobile = ?,
    email = ?,
    profile_pic_id = ?
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
