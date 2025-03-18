package model

import (
	"database/sql"
	"sgublogsite/src/internal/model/repos"
)

func (m *Model) GetPosts() ([]repos.Post, error) {
    return m.query.GetAllPosts(m.ctx)
}

func (m *Model) GetPostsByUserID(id int32) ([]repos.Post, error) {
    return m.query.GetPostsByUserID(m.ctx, sql.NullInt32{Int32: id, Valid: true})
}

func (m *Model) GetPostsByCategoryID(id int32) ([]repos.Post, error) {
    return m.query.GetPostsByCategoryID(m.ctx, id)
}

func (m *Model) GetUncategorizedPosts() ([]repos.Post, error) {
    return m.query.GetUncategorizedPosts(m.ctx)
}

func (m *Model) GetPostsByCategoryName(name string) ([]repos.Post, error) {
    return m.query.GetPostsByCategoryName(m.ctx, name)
}

func (m *Model) GetPostsByTagID(id int32) ([]repos.Post, error) {
    return m.query.GetPostsByTagID(m.ctx, id)
}

func (m *Model) GetPostsByTagName(name string) ([]repos.Post, error) {
    return m.query.GetPostsByTagName(m.ctx, name)
}

func (m *Model) GetPostsByStatus(status string) ([]repos.Post, error) {
    return m.query.GetPostsByStatus(m.ctx, repos.PostsStatus(status))
}

func (m *Model) SearchPosts(text string) ([]repos.Post, error) {
    wildcard := "%" + text + "%"
    return m.query.FindPosts(m.ctx, wildcard)
}

func (m *Model) CreatePost(post repos.Post) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.CreatePost(m.ctx, repos.CreatePostParams{
        UserID:         post.UserID,
        Title:          post.Title,
        Slug:           post.Slug,
        ThumbnailID:   post.ThumbnailID,
        Body:           post.Body,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) AddPostToCategory(postID int32, categoryID int32) error {
    tx, err := m.db.Begin()
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
    tx, err := m.db.Begin()
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
    tx, err := m.db.Begin()
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
    tx, err := m.db.Begin()
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

    err = New().SetPostPrivate(postID, isPrivate)

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

func (m *Model) UpdatePostMetadata(post repos.Post) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdatePostMetadata(m.ctx, repos.UpdatePostMetadataParams{
        PostID:         post.PostID,
        Title:          post.Title,
        Slug:           post.Slug,
        ThumbnailID:   post.ThumbnailID,
    })

    if err != nil {
        return err
    }

    err = New().draftPost(qtx, post.PostID)

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) UpdatePostBody(post repos.Post) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdatePostBody(m.ctx, repos.UpdatePostBodyParams{
        PostID:         post.PostID,
        Body:           post.Body,
    })

    if err != nil {
        return err
    }

    err = New().draftPost(qtx, post.PostID)

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) DeletePost(id int32) error {
    tx, err := m.db.Begin()
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
