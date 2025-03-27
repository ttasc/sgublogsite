-- name: GetCategoryByID :one
SELECT *
FROM categories
WHERE category_id = ?;

-- name: GetAllCategories :many
SELECT * FROM categories;

-- name: GetChildCategories :many
SELECT *
FROM categories
WHERE parent_category_id = ?;

-- name: GetRootCategories :many
SELECT *
FROM categories
WHERE parent_category_id IS NULL;

-- name: AddCategory :execresult
INSERT INTO categories (
    parent_category_id,
    name,
    slug
) VALUES (?, ?, ?);

-- name: UpdateCategory :execresult
UPDATE categories
SET name = ?,
    slug = ?
WHERE category_id = ?;

-- name: UpdateCategoryParent :execresult
UPDATE categories
SET parent_category_id = ?
WHERE category_id = ?;

-- name: DeleteCategory :execresult
DELETE FROM categories
WHERE category_id = ?;
