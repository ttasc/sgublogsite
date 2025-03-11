package services

import (
	"database/sql"
	"sgublogsite/src/internal/model/repos"
)

func GetCategoryByID(id int32) (repos.Category, error) {
    return s.query.GetCategoryByID(s.ctx, id)
}

func GetCategories() ([]repos.Category, error) {
    return s.query.GetAllCategories(s.ctx)
}

func GetChildCategories(id int32) ([]repos.Category, error) {
    return s.query.GetChildCategories(s.ctx, sql.NullInt32{Int32: id, Valid: true})
}

func GetRootCategories() ([]repos.Category, error) {
    return s.query.GetRootCategories(s.ctx)
}

func AddCategory(category repos.Category) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.AddCategory(s.ctx, repos.AddCategoryParams{
        ParentCategoryID:   category.ParentCategoryID,
        Name:               category.Name,
        Slug:               category.Slug,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func UpdateCategory(category repos.Category) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.UpdateCategory(s.ctx, repos.UpdateCategoryParams{
        ParentCategoryID:   category.ParentCategoryID,
        Name:               category.Name,
        Slug:               category.Slug,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func DeleteCategory(id int32) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.DeleteCategory(s.ctx, id)

    if err != nil {
        return err
    }

    return tx.Commit()
}
