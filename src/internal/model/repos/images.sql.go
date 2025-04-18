// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: images.sql

package repos

import (
	"context"
	"database/sql"
)

const addImage = `-- name: AddImage :execresult
INSERT INTO images (url, name)
VALUES (?, ?)
`

type AddImageParams struct {
	Url  string         `json:"url"`
	Name sql.NullString `json:"name"`
}

func (q *Queries) AddImage(ctx context.Context, arg AddImageParams) (sql.Result, error) {
	return q.exec(ctx, q.addImageStmt, addImage, arg.Url, arg.Name)
}

const countImages = `-- name: CountImages :one
SELECT COUNT(*) FROM images
`

func (q *Queries) CountImages(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.countImagesStmt, countImages)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const deleteImage = `-- name: DeleteImage :execresult
DELETE FROM images
WHERE image_id = ?
`

func (q *Queries) DeleteImage(ctx context.Context, imageID int32) (sql.Result, error) {
	return q.exec(ctx, q.deleteImageStmt, deleteImage, imageID)
}

const getAllImages = `-- name: GetAllImages :many
SELECT image_id, url, name FROM images
LIMIT ? OFFSET ?
`

type GetAllImagesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllImages(ctx context.Context, arg GetAllImagesParams) ([]Image, error) {
	rows, err := q.query(ctx, q.getAllImagesStmt, getAllImages, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Image
	for rows.Next() {
		var i Image
		if err := rows.Scan(&i.ImageID, &i.Url, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getImageByID = `-- name: GetImageByID :one
SELECT image_id, url, name FROM images
WHERE image_id = ?
`

func (q *Queries) GetImageByID(ctx context.Context, imageID int32) (Image, error) {
	row := q.queryRow(ctx, q.getImageByIDStmt, getImageByID, imageID)
	var i Image
	err := row.Scan(&i.ImageID, &i.Url, &i.Name)
	return i, err
}

const getImageByURL = `-- name: GetImageByURL :one
SELECT image_id, url, name FROM images
WHERE url = ?
`

func (q *Queries) GetImageByURL(ctx context.Context, url string) (Image, error) {
	row := q.queryRow(ctx, q.getImageByURLStmt, getImageByURL, url)
	var i Image
	err := row.Scan(&i.ImageID, &i.Url, &i.Name)
	return i, err
}

const updateImageURL = `-- name: UpdateImageURL :execresult
UPDATE images
SET url = ?
WHERE image_id = ?
`

type UpdateImageURLParams struct {
	Url     string `json:"url"`
	ImageID int32  `json:"image_id"`
}

func (q *Queries) UpdateImageURL(ctx context.Context, arg UpdateImageURLParams) (sql.Result, error) {
	return q.exec(ctx, q.updateImageURLStmt, updateImageURL, arg.Url, arg.ImageID)
}
