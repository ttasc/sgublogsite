package services

import (
	"database/sql"
	"sgublogsite/src/internal/model/repos"
)

func GetPosts() ([]repos.Post, error) {
    return s.query.GetAllPosts(s.ctx)
}

func GetPostsByUserID(id int32) ([]repos.Post, error) {
    return s.query.GetPostsByUserID(s.ctx, sql.NullInt32{Int32: id, Valid: true})
}

func GetPostsByCategoryID(id int32) ([]repos.Post, error) {
    return s.query.GetPostsByCategoryID(s.ctx, id)
}

func GetUncategorizedPosts() ([]repos.Post, error) {
    return s.query.GetUncategorizedPosts(s.ctx)
}

func GetPostsByCategoryName(name string) ([]repos.Post, error) {
    return s.query.GetPostsByCategoryName(s.ctx, name)
}

func GetPostsByTagID(id int32) ([]repos.Post, error) {
    return s.query.GetPostsByTagID(s.ctx, id)
}

func GetPostsByTagName(name string) ([]repos.Post, error) {
    return s.query.GetPostsByTagName(s.ctx, name)
}

func GetPostsByStatus(status string) ([]repos.Post, error) {
    return s.query.GetPostsByStatus(s.ctx, repos.PostsStatus(status))
}

func SearchPosts(text string) ([]repos.Post, error) {
    wildcard := "%" + text + "%"
    return s.query.FindPosts(s.ctx, wildcard)
}

func CreatePost(post repos.Post) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.CreatePost(s.ctx, repos.CreatePostParams{
        UserID:         post.UserID,
        Title:          post.Title,
        Slug:           post.Slug,
        PreviewPicID:   post.PreviewPicID,
        Body:           post.Body,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func AddPostToCategory(postID int32, categoryID int32) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.AddPostToCategory(s.ctx, repos.AddPostToCategoryParams{
        PostID:         postID,
        CategoryID:     categoryID,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func AddTagToPost(postID int32, tagID int32) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.AddTagToPost(s.ctx, repos.AddTagToPostParams{
        PostID:         postID,
        TagID:          tagID,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func PublishPost(postID int32) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.UpdatePostStatus(s.ctx, repos.UpdatePostStatusParams{
        PostID:         postID,
        Status:         repos.PostsStatusPublished,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func draftPost(qtx *repos.Queries, postID int32) error {
    _, err := qtx.UpdatePostStatus(s.ctx, repos.UpdatePostStatusParams{
        PostID:         postID,
        Status:         repos.PostsStatusDraft,
    })
    return err
}

func UpdatePostMetadata(post repos.Post) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.UpdatePostMetadata(s.ctx, repos.UpdatePostMetadataParams{
        Title:          post.Title,
        Slug:           post.Slug,
        PreviewPicID:   post.PreviewPicID,
    })

    if err != nil {
        return err
    }

    err = draftPost(qtx, post.PostID)

    if err != nil {
        return err
    }

    return tx.Commit()
}

func UpdatePostBody(post repos.Post) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.UpdatePostBody(s.ctx, repos.UpdatePostBodyParams{
        Body:           post.Body,
    })

    if err != nil {
        return err
    }

    err = draftPost(qtx, post.PostID)

    if err != nil {
        return err
    }

    return tx.Commit()
}

func DeletePost(id int32) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.DeletePost(s.ctx, id)

    if err != nil {
        return err
    }

    return tx.Commit()
}
