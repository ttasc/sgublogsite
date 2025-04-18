// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: posts.sql

package repos

import (
	"context"
	"database/sql"
	"time"
)

const addPostToCategory = `-- name: AddPostToCategory :execresult
INSERT INTO post_categories (
    post_id,
    category_id
) VALUES (?, ?)
`

type AddPostToCategoryParams struct {
	PostID     int32 `json:"post_id"`
	CategoryID int32 `json:"category_id"`
}

func (q *Queries) AddPostToCategory(ctx context.Context, arg AddPostToCategoryParams) (sql.Result, error) {
	return q.exec(ctx, q.addPostToCategoryStmt, addPostToCategory, arg.PostID, arg.CategoryID)
}

const addTagToPost = `-- name: AddTagToPost :execresult
INSERT INTO post_tags (
    post_id,
    tag_id
) VALUES (?, ?)
`

type AddTagToPostParams struct {
	PostID int32 `json:"post_id"`
	TagID  int32 `json:"tag_id"`
}

func (q *Queries) AddTagToPost(ctx context.Context, arg AddTagToPostParams) (sql.Result, error) {
	return q.exec(ctx, q.addTagToPostStmt, addTagToPost, arg.PostID, arg.TagID)
}

const countPosts = `-- name: CountPosts :one
SELECT COUNT(*) FROM posts
`

func (q *Queries) CountPosts(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.countPostsStmt, countPosts)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPost = `-- name: CreatePost :execresult
INSERT INTO posts (
    user_id,
    title,
    slug,
    thumbnail_id,
    body
) VALUES (?, ?, ?, ?, ?)
`

type CreatePostParams struct {
	UserID      sql.NullInt32 `json:"user_id"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	ThumbnailID sql.NullInt32 `json:"thumbnail_id"`
	Body        string        `json:"body"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (sql.Result, error) {
	return q.exec(ctx, q.createPostStmt, createPost,
		arg.UserID,
		arg.Title,
		arg.Slug,
		arg.ThumbnailID,
		arg.Body,
	)
}

const deletePost = `-- name: DeletePost :execresult
DELETE FROM posts
WHERE post_id = ?
`

func (q *Queries) DeletePost(ctx context.Context, postID int32) (sql.Result, error) {
	return q.exec(ctx, q.deletePostStmt, deletePost, postID)
}

const findPosts = `-- name: FindPosts :many
SELECT posts.post_id, posts.user_id, posts.title, posts.slug, posts.thumbnail_id, posts.body, posts.status, posts.private, posts.created_at, posts.updated_at, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE lower(concat(title, ' ', body)) LIKE lower(?)
AND (private = ? OR private = 0)
AND status = ?
LIMIT ? OFFSET ?
`

type FindPostsParams struct {
	Text    string      `json:"text"`
	Private bool        `json:"private"`
	Status  PostsStatus `json:"status"`
	Limit   int32       `json:"limit"`
	Offset  int32       `json:"offset"`
}

type FindPostsRow struct {
	PostID      int32          `json:"post_id"`
	UserID      sql.NullInt32  `json:"user_id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	ThumbnailID sql.NullInt32  `json:"thumbnail_id"`
	Body        string         `json:"body"`
	Status      PostsStatus    `json:"status"`
	Private     bool           `json:"private"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Thumbnail   sql.NullString `json:"thumbnail"`
}

// WHERE MATCH(title, body) AGAINST (sqlc.arg(text));
func (q *Queries) FindPosts(ctx context.Context, arg FindPostsParams) ([]FindPostsRow, error) {
	rows, err := q.query(ctx, q.findPostsStmt, findPosts,
		arg.Text,
		arg.Private,
		arg.Status,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindPostsRow
	for rows.Next() {
		var i FindPostsRow
		if err := rows.Scan(
			&i.PostID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.ThumbnailID,
			&i.Body,
			&i.Status,
			&i.Private,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Thumbnail,
		); err != nil {
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

const getFilteredPosts = `-- name: GetFilteredPosts :many
SELECT
    p.post_id,
    p.title,
    CONCAT(u.firstname, ' ', u.lastname) AS author_name,
    GROUP_CONCAT(c.name SEPARATOR ', ') AS categories,
    p.created_at,
    p.status
FROM posts p
LEFT JOIN users u ON p.user_id = u.user_id
LEFT JOIN post_categories pc ON p.post_id = pc.post_id
LEFT JOIN categories c ON pc.category_id = c.category_id
WHERE
    (CASE WHEN ?      != ''       THEN lower(p.title) LIKE lower(?) ELSE TRUE END) AND
    (CASE WHEN ?     != ''       THEN p.status       =  ? ELSE TRUE END) AND
    (CASE WHEN ?    IS NOT NULL THEN p.private      =  ? ELSE TRUE END)
GROUP BY p.post_id, p.title, author_name, p.created_at, p.status
ORDER BY p.created_at DESC
LIMIT ? OFFSET ?
`

type GetFilteredPostsParams struct {
	Title   string      `json:"title"`
	Status  PostsStatus `json:"status"`
	Private bool        `json:"private"`
	Limit   int32       `json:"limit"`
	Offset  int32       `json:"offset"`
}

type GetFilteredPostsRow struct {
	PostID     int32          `json:"post_id"`
	Title      string         `json:"title"`
	AuthorName string         `json:"author_name"`
	Categories sql.NullString `json:"categories"`
	CreatedAt  time.Time      `json:"created_at"`
	Status     PostsStatus    `json:"status"`
}

func (q *Queries) GetFilteredPosts(ctx context.Context, arg GetFilteredPostsParams) ([]GetFilteredPostsRow, error) {
	rows, err := q.query(ctx, q.getFilteredPostsStmt, getFilteredPosts,
		arg.Title,
		arg.Title,
		arg.Status,
		arg.Status,
		arg.Private,
		arg.Private,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFilteredPostsRow
	for rows.Next() {
		var i GetFilteredPostsRow
		if err := rows.Scan(
			&i.PostID,
			&i.Title,
			&i.AuthorName,
			&i.Categories,
			&i.CreatedAt,
			&i.Status,
		); err != nil {
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

const getPostByID = `-- name: GetPostByID :one
SELECT posts.post_id, posts.user_id, posts.title, posts.slug, posts.thumbnail_id, posts.body, posts.status, posts.private, posts.created_at, posts.updated_at, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE post_id = ?
`

type GetPostByIDRow struct {
	PostID      int32          `json:"post_id"`
	UserID      sql.NullInt32  `json:"user_id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	ThumbnailID sql.NullInt32  `json:"thumbnail_id"`
	Body        string         `json:"body"`
	Status      PostsStatus    `json:"status"`
	Private     bool           `json:"private"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Thumbnail   sql.NullString `json:"thumbnail"`
}

func (q *Queries) GetPostByID(ctx context.Context, postID int32) (GetPostByIDRow, error) {
	row := q.queryRow(ctx, q.getPostByIDStmt, getPostByID, postID)
	var i GetPostByIDRow
	err := row.Scan(
		&i.PostID,
		&i.UserID,
		&i.Title,
		&i.Slug,
		&i.ThumbnailID,
		&i.Body,
		&i.Status,
		&i.Private,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Thumbnail,
	)
	return i, err
}

const getPosts = `-- name: GetPosts :many
SELECT posts.post_id, posts.user_id, posts.title, posts.slug, posts.thumbnail_id, posts.body, posts.status, posts.private, posts.created_at, posts.updated_at, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?
`

type GetPostsParams struct {
	Status PostsStatus `json:"status"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

type GetPostsRow struct {
	PostID      int32          `json:"post_id"`
	UserID      sql.NullInt32  `json:"user_id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	ThumbnailID sql.NullInt32  `json:"thumbnail_id"`
	Body        string         `json:"body"`
	Status      PostsStatus    `json:"status"`
	Private     bool           `json:"private"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Thumbnail   sql.NullString `json:"thumbnail"`
}

func (q *Queries) GetPosts(ctx context.Context, arg GetPostsParams) ([]GetPostsRow, error) {
	rows, err := q.query(ctx, q.getPostsStmt, getPosts, arg.Status, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsRow
	for rows.Next() {
		var i GetPostsRow
		if err := rows.Scan(
			&i.PostID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.ThumbnailID,
			&i.Body,
			&i.Status,
			&i.Private,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Thumbnail,
		); err != nil {
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

const getPostsByCategoryID = `-- name: GetPostsByCategoryID :many
SELECT posts.post_id, posts.user_id, posts.title, posts.slug, posts.thumbnail_id, posts.body, posts.status, posts.private, posts.created_at, posts.updated_at, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE post_id IN (
    SELECT post_id
    FROM post_categories
    WHERE category_id = ?
)
AND (private = ? OR private = 0)
AND status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?
`

type GetPostsByCategoryIDParams struct {
	CategoryID int32       `json:"category_id"`
	Private    bool        `json:"private"`
	Status     PostsStatus `json:"status"`
	Limit      int32       `json:"limit"`
	Offset     int32       `json:"offset"`
}

type GetPostsByCategoryIDRow struct {
	PostID      int32          `json:"post_id"`
	UserID      sql.NullInt32  `json:"user_id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	ThumbnailID sql.NullInt32  `json:"thumbnail_id"`
	Body        string         `json:"body"`
	Status      PostsStatus    `json:"status"`
	Private     bool           `json:"private"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Thumbnail   sql.NullString `json:"thumbnail"`
}

func (q *Queries) GetPostsByCategoryID(ctx context.Context, arg GetPostsByCategoryIDParams) ([]GetPostsByCategoryIDRow, error) {
	rows, err := q.query(ctx, q.getPostsByCategoryIDStmt, getPostsByCategoryID,
		arg.CategoryID,
		arg.Private,
		arg.Status,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsByCategoryIDRow
	for rows.Next() {
		var i GetPostsByCategoryIDRow
		if err := rows.Scan(
			&i.PostID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.ThumbnailID,
			&i.Body,
			&i.Status,
			&i.Private,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Thumbnail,
		); err != nil {
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

const getPostsByCategorySlug = `-- name: GetPostsByCategorySlug :many
SELECT posts.post_id, posts.user_id, posts.title, posts.slug, posts.thumbnail_id, posts.body, posts.status, posts.private, posts.created_at, posts.updated_at, images.url AS thumbnail
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
LIMIT ? OFFSET ?
`

type GetPostsByCategorySlugParams struct {
	Slug    string      `json:"slug"`
	Private bool        `json:"private"`
	Status  PostsStatus `json:"status"`
	Limit   int32       `json:"limit"`
	Offset  int32       `json:"offset"`
}

type GetPostsByCategorySlugRow struct {
	PostID      int32          `json:"post_id"`
	UserID      sql.NullInt32  `json:"user_id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	ThumbnailID sql.NullInt32  `json:"thumbnail_id"`
	Body        string         `json:"body"`
	Status      PostsStatus    `json:"status"`
	Private     bool           `json:"private"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Thumbnail   sql.NullString `json:"thumbnail"`
}

func (q *Queries) GetPostsByCategorySlug(ctx context.Context, arg GetPostsByCategorySlugParams) ([]GetPostsByCategorySlugRow, error) {
	rows, err := q.query(ctx, q.getPostsByCategorySlugStmt, getPostsByCategorySlug,
		arg.Slug,
		arg.Private,
		arg.Status,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsByCategorySlugRow
	for rows.Next() {
		var i GetPostsByCategorySlugRow
		if err := rows.Scan(
			&i.PostID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.ThumbnailID,
			&i.Body,
			&i.Status,
			&i.Private,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Thumbnail,
		); err != nil {
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

const getPostsByStatus = `-- name: GetPostsByStatus :many
SELECT posts.post_id, posts.user_id, posts.title, posts.slug, posts.thumbnail_id, posts.body, posts.status, posts.private, posts.created_at, posts.updated_at, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?
`

type GetPostsByStatusParams struct {
	Status PostsStatus `json:"status"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

type GetPostsByStatusRow struct {
	PostID      int32          `json:"post_id"`
	UserID      sql.NullInt32  `json:"user_id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	ThumbnailID sql.NullInt32  `json:"thumbnail_id"`
	Body        string         `json:"body"`
	Status      PostsStatus    `json:"status"`
	Private     bool           `json:"private"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Thumbnail   sql.NullString `json:"thumbnail"`
}

func (q *Queries) GetPostsByStatus(ctx context.Context, arg GetPostsByStatusParams) ([]GetPostsByStatusRow, error) {
	rows, err := q.query(ctx, q.getPostsByStatusStmt, getPostsByStatus, arg.Status, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsByStatusRow
	for rows.Next() {
		var i GetPostsByStatusRow
		if err := rows.Scan(
			&i.PostID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.ThumbnailID,
			&i.Body,
			&i.Status,
			&i.Private,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Thumbnail,
		); err != nil {
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

const getPostsByTagID = `-- name: GetPostsByTagID :many
SELECT posts.post_id, posts.user_id, posts.title, posts.slug, posts.thumbnail_id, posts.body, posts.status, posts.private, posts.created_at, posts.updated_at, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE post_id IN (
    SELECT post_id
    FROM post_tags
    WHERE tag_id = ?
)
AND (private = ? OR private = 0)
AND status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?
`

type GetPostsByTagIDParams struct {
	TagID   int32       `json:"tag_id"`
	Private bool        `json:"private"`
	Status  PostsStatus `json:"status"`
	Limit   int32       `json:"limit"`
	Offset  int32       `json:"offset"`
}

type GetPostsByTagIDRow struct {
	PostID      int32          `json:"post_id"`
	UserID      sql.NullInt32  `json:"user_id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	ThumbnailID sql.NullInt32  `json:"thumbnail_id"`
	Body        string         `json:"body"`
	Status      PostsStatus    `json:"status"`
	Private     bool           `json:"private"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Thumbnail   sql.NullString `json:"thumbnail"`
}

func (q *Queries) GetPostsByTagID(ctx context.Context, arg GetPostsByTagIDParams) ([]GetPostsByTagIDRow, error) {
	rows, err := q.query(ctx, q.getPostsByTagIDStmt, getPostsByTagID,
		arg.TagID,
		arg.Private,
		arg.Status,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsByTagIDRow
	for rows.Next() {
		var i GetPostsByTagIDRow
		if err := rows.Scan(
			&i.PostID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.ThumbnailID,
			&i.Body,
			&i.Status,
			&i.Private,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Thumbnail,
		); err != nil {
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

const getPostsByTagSlug = `-- name: GetPostsByTagSlug :many
SELECT posts.post_id, posts.user_id, posts.title, posts.slug, posts.thumbnail_id, posts.body, posts.status, posts.private, posts.created_at, posts.updated_at, images.url AS thumbnail
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
LIMIT ? OFFSET ?
`

type GetPostsByTagSlugParams struct {
	Slug    string      `json:"slug"`
	Private bool        `json:"private"`
	Status  PostsStatus `json:"status"`
	Limit   int32       `json:"limit"`
	Offset  int32       `json:"offset"`
}

type GetPostsByTagSlugRow struct {
	PostID      int32          `json:"post_id"`
	UserID      sql.NullInt32  `json:"user_id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	ThumbnailID sql.NullInt32  `json:"thumbnail_id"`
	Body        string         `json:"body"`
	Status      PostsStatus    `json:"status"`
	Private     bool           `json:"private"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Thumbnail   sql.NullString `json:"thumbnail"`
}

func (q *Queries) GetPostsByTagSlug(ctx context.Context, arg GetPostsByTagSlugParams) ([]GetPostsByTagSlugRow, error) {
	rows, err := q.query(ctx, q.getPostsByTagSlugStmt, getPostsByTagSlug,
		arg.Slug,
		arg.Private,
		arg.Status,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsByTagSlugRow
	for rows.Next() {
		var i GetPostsByTagSlugRow
		if err := rows.Scan(
			&i.PostID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.ThumbnailID,
			&i.Body,
			&i.Status,
			&i.Private,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Thumbnail,
		); err != nil {
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

const getPostsByUserID = `-- name: GetPostsByUserID :many
SELECT posts.post_id, posts.user_id, posts.title, posts.slug, posts.thumbnail_id, posts.body, posts.status, posts.private, posts.created_at, posts.updated_at, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE user_id = ?
AND status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?
`

type GetPostsByUserIDParams struct {
	UserID sql.NullInt32 `json:"user_id"`
	Status PostsStatus   `json:"status"`
	Limit  int32         `json:"limit"`
	Offset int32         `json:"offset"`
}

type GetPostsByUserIDRow struct {
	PostID      int32          `json:"post_id"`
	UserID      sql.NullInt32  `json:"user_id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	ThumbnailID sql.NullInt32  `json:"thumbnail_id"`
	Body        string         `json:"body"`
	Status      PostsStatus    `json:"status"`
	Private     bool           `json:"private"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Thumbnail   sql.NullString `json:"thumbnail"`
}

func (q *Queries) GetPostsByUserID(ctx context.Context, arg GetPostsByUserIDParams) ([]GetPostsByUserIDRow, error) {
	rows, err := q.query(ctx, q.getPostsByUserIDStmt, getPostsByUserID,
		arg.UserID,
		arg.Status,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsByUserIDRow
	for rows.Next() {
		var i GetPostsByUserIDRow
		if err := rows.Scan(
			&i.PostID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.ThumbnailID,
			&i.Body,
			&i.Status,
			&i.Private,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Thumbnail,
		); err != nil {
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

const getUncategorizedPosts = `-- name: GetUncategorizedPosts :many
SELECT posts.post_id, posts.user_id, posts.title, posts.slug, posts.thumbnail_id, posts.body, posts.status, posts.private, posts.created_at, posts.updated_at, images.url AS thumbnail
FROM posts LEFT JOIN images ON posts.thumbnail_id = images.image_id
WHERE post_id NOT IN (
    SELECT post_id
    FROM post_categories
)
AND (private = ? OR private = 0)
AND status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?
`

type GetUncategorizedPostsParams struct {
	Private bool        `json:"private"`
	Status  PostsStatus `json:"status"`
	Limit   int32       `json:"limit"`
	Offset  int32       `json:"offset"`
}

type GetUncategorizedPostsRow struct {
	PostID      int32          `json:"post_id"`
	UserID      sql.NullInt32  `json:"user_id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	ThumbnailID sql.NullInt32  `json:"thumbnail_id"`
	Body        string         `json:"body"`
	Status      PostsStatus    `json:"status"`
	Private     bool           `json:"private"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Thumbnail   sql.NullString `json:"thumbnail"`
}

func (q *Queries) GetUncategorizedPosts(ctx context.Context, arg GetUncategorizedPostsParams) ([]GetUncategorizedPostsRow, error) {
	rows, err := q.query(ctx, q.getUncategorizedPostsStmt, getUncategorizedPosts,
		arg.Private,
		arg.Status,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUncategorizedPostsRow
	for rows.Next() {
		var i GetUncategorizedPostsRow
		if err := rows.Scan(
			&i.PostID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.ThumbnailID,
			&i.Body,
			&i.Status,
			&i.Private,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Thumbnail,
		); err != nil {
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

const updatePostBody = `-- name: UpdatePostBody :execresult
UPDATE posts
SET body = ?
WHERE post_id = ?
`

type UpdatePostBodyParams struct {
	Body   string `json:"body"`
	PostID int32  `json:"post_id"`
}

func (q *Queries) UpdatePostBody(ctx context.Context, arg UpdatePostBodyParams) (sql.Result, error) {
	return q.exec(ctx, q.updatePostBodyStmt, updatePostBody, arg.Body, arg.PostID)
}

const updatePostMetadata = `-- name: UpdatePostMetadata :execresult
UPDATE posts SET
    title = ?,
    slug = ?,
    thumbnail_id = ?
WHERE post_id = ?
`

type UpdatePostMetadataParams struct {
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	ThumbnailID sql.NullInt32 `json:"thumbnail_id"`
	PostID      int32         `json:"post_id"`
}

func (q *Queries) UpdatePostMetadata(ctx context.Context, arg UpdatePostMetadataParams) (sql.Result, error) {
	return q.exec(ctx, q.updatePostMetadataStmt, updatePostMetadata,
		arg.Title,
		arg.Slug,
		arg.ThumbnailID,
		arg.PostID,
	)
}

const updatePostPrivate = `-- name: UpdatePostPrivate :execresult
UPDATE posts
SET private = ?
WHERE post_id = ?
`

type UpdatePostPrivateParams struct {
	Private bool  `json:"private"`
	PostID  int32 `json:"post_id"`
}

func (q *Queries) UpdatePostPrivate(ctx context.Context, arg UpdatePostPrivateParams) (sql.Result, error) {
	return q.exec(ctx, q.updatePostPrivateStmt, updatePostPrivate, arg.Private, arg.PostID)
}

const updatePostStatus = `-- name: UpdatePostStatus :execresult
UPDATE posts
SET status = ?
WHERE post_id = ?
`

type UpdatePostStatusParams struct {
	Status PostsStatus `json:"status"`
	PostID int32       `json:"post_id"`
}

func (q *Queries) UpdatePostStatus(ctx context.Context, arg UpdatePostStatusParams) (sql.Result, error) {
	return q.exec(ctx, q.updatePostStatusStmt, updatePostStatus, arg.Status, arg.PostID)
}
