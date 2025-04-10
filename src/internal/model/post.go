package model

import (
	"database/sql"

	"github.com/ttasc/sgublogsite/src/internal/model/repos"
)

type GetPost repos.GetPostByIDRow

func (m *Model) GetPostByID(id int32) (repos.GetPostByIDRow, error) {
    return m.query.GetPostByID(m.ctx, id)
}

func (m *Model) CountPosts() (int64, error) {
    return m.query.CountPosts(m.ctx)
}

func (m *Model) GetFilteredPosts(
    limit, offset int32,
    title string,
    status repos.PostsStatus,
    isPrivate bool) ([]repos.GetFilteredPostsRow, error) {
    wildcard := "%" + title + "%"
    return m.query.GetFilteredPosts(m.ctx, repos.GetFilteredPostsParams{
        Limit:          limit,
        Offset:         offset,

        Title:          wildcard,
        Status:         repos.PostsStatus(status),
        Private:        isPrivate,
    })
}

func (m *Model) GetPosts(limit, offset int32, status string) ([]repos.GetPostsRow, error) {
    return m.query.GetPosts(m.ctx, repos.GetPostsParams{
        Limit:          limit,
        Offset:         offset,
        Status:         repos.PostsStatus(status),
    })
}

func (m *Model) GetPostsByUserID(id, limit, offset int32, status string) ([]repos.GetPostsByUserIDRow, error) {
    if id < 1 {
        return nil, nil
    }
    return m.query.GetPostsByUserID(m.ctx, repos.GetPostsByUserIDParams{
        UserID:         sql.NullInt32{Int32: id, Valid: true},
        Limit:          limit,
        Offset:         offset,
        Status:         repos.PostsStatus(status),
    })
}

func (m *Model) GetPostsByCategoryID(id, limit, offset int32, status string, getPrivate bool) ([]repos.GetPostsByCategoryIDRow, error) {
    return m.query.GetPostsByCategoryID(m.ctx, repos.GetPostsByCategoryIDParams{
        CategoryID:     id,
        Limit:          limit,
        Offset:         offset,
        Private:        getPrivate,
        Status:         repos.PostsStatus(status),
    })
}

func (m *Model) GetUncategorizedPosts(limit, offset int32, status string, getPrivate bool) ([]repos.GetUncategorizedPostsRow, error) {
    return m.query.GetUncategorizedPosts(m.ctx, repos.GetUncategorizedPostsParams{
        Limit:          limit,
        Offset:         offset,
        Status:         repos.PostsStatus(status),
        Private:        getPrivate,
    })
}

func (m *Model) GetPostsByCategorySlug(slug string, limit, offset int32, status string, getPrivate bool) ([]repos.GetPostsByCategorySlugRow, error) {
    return m.query.GetPostsByCategorySlug(m.ctx, repos.GetPostsByCategorySlugParams{
        Slug:           slug,
        Limit:          limit,
        Offset:         offset,
        Status:         repos.PostsStatus(status),
        Private:        getPrivate,
    })
}

func (m *Model) GetPostsByTagID(id, limit, offset int32, status string, getPrivate bool) ([]repos.GetPostsByTagIDRow, error) {
    return m.query.GetPostsByTagID(m.ctx, repos.GetPostsByTagIDParams{
        TagID:          id,
        Limit:          limit,
        Offset:         offset,
        Status:         repos.PostsStatus(status),
        Private:        getPrivate,
    })
}

func (m *Model) GetPostsByTagSlug(slug string, limit, offset int32, status string, getPrivate bool) ([]repos.GetPostsByTagSlugRow, error) {
    return m.query.GetPostsByTagSlug(m.ctx, repos.GetPostsByTagSlugParams{
        Slug:           slug,
        Limit:          limit,
        Offset:         offset,
        Status:         repos.PostsStatus(status),
        Private:        getPrivate,
    })
}

func (m *Model) GetPostsByStatus(status string, limit, offset int32) ([]repos.GetPostsByStatusRow, error) {
    return m.query.GetPostsByStatus(m.ctx, repos.GetPostsByStatusParams{
        Status:         repos.PostsStatus(status),
        Limit:          limit,
        Offset:         offset,
    })
}

func (m *Model) SearchPosts(text string, limit, offset int32, status string, getPrivate bool) ([]repos.FindPostsRow, error) {
    wildcard := "%" + text + "%"
    return m.query.FindPosts(m.ctx, repos.FindPostsParams{
        Text:           wildcard,
        Limit:          limit,
        Offset:         offset,
        Status:         repos.PostsStatus(status),
        Private:        getPrivate,
    })
}

func (m *Model) CreatePost(id int32, title string, slug string, thumbnailID int32, body string) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.CreatePost(m.ctx, repos.CreatePostParams{
        UserID:         sql.NullInt32{Int32: id, Valid: id > 0},
        Title:          title,
        Slug:           slug,
        ThumbnailID:    sql.NullInt32{Int32: thumbnailID, Valid: thumbnailID > 0},
        Body:           body,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) AddPostToCategory(postID int32, categoryID int32) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.AddPostToCategory(m.ctx, repos.AddPostToCategoryParams{
        PostID:         postID,
        CategoryID:     categoryID,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) AddTagToPost(postID int32, tagID int32) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.AddTagToPost(m.ctx, repos.AddTagToPostParams{
        PostID:         postID,
        TagID:          tagID,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) SetPostPrivate(postID int32, isPrivate bool) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdatePostPrivate(m.ctx, repos.UpdatePostPrivateParams{
        PostID:         postID,
        Private:        isPrivate,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) PublishPost(postID int32, isPrivate bool) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdatePostStatus(m.ctx, repos.UpdatePostStatusParams{
        PostID:         postID,
        Status:         repos.PostsStatusPublished,
    })

    if err != nil {
        return err
    }

    err = m.SetPostPrivate(postID, isPrivate)

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) draftPost(qtx *repos.Queries, postID int32) error {
    _, err := qtx.UpdatePostStatus(m.ctx, repos.UpdatePostStatusParams{
        PostID:         postID,
        Status:         repos.PostsStatusDraft,
    })
    return err
}

func (m *Model) UpdatePostMetadata(id int32, title string, slug string, thumbnailID int32) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdatePostMetadata(m.ctx, repos.UpdatePostMetadataParams{
        PostID:         id,
        Title:          title,
        Slug:           slug,
        ThumbnailID:    sql.NullInt32{Int32: thumbnailID, Valid: thumbnailID > 0},
    })

    if err != nil {
        return err
    }

    err = m.draftPost(qtx, id)

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) UpdatePostBody(id int32, body string) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdatePostBody(m.ctx, repos.UpdatePostBodyParams{
        PostID:         id,
        Body:           body,
    })

    if err != nil {
        return err
    }

    err = m.draftPost(qtx, id)

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) DeletePost(id int32) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.DeletePost(m.ctx, id)

    if err != nil {
        return err
    }

    return tx.Commit()
}
