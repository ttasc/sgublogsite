package model

import (
	"database/sql"
	"sgublogsite/src/internal/model/repos"
)

func (m *Model) GetCategoryByID(id int32) (repos.Category, error) {
    return m.query.GetCategoryByID(m.ctx, id)
}

func (m *Model) GetCategories() ([]repos.Category, error) {
    return m.query.GetAllCategories(m.ctx)
}

func (m *Model) GetChildCategories(id int32) ([]repos.Category, error) {
    return m.query.GetChildCategories(m.ctx, sql.NullInt32{Int32: id, Valid: true})
}

func (m *Model) GetRootCategories() ([]repos.Category, error) {
    return m.query.GetRootCategories(m.ctx)
}

func (m *Model) AddCategory(category repos.Category) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.AddCategory(m.ctx, repos.AddCategoryParams{
        ParentCategoryID:   category.ParentCategoryID,
        Name:               category.Name,
        Slug:               category.Slug,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) UpdateCategory(category repos.Category) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdateCategory(m.ctx, repos.UpdateCategoryParams{
        CategoryID:         category.CategoryID,
        ParentCategoryID:   category.ParentCategoryID,
        Name:               category.Name,
        Slug:               category.Slug,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) DeleteCategory(id int32) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.DeleteCategory(m.ctx, id)

    if err != nil {
        return err
    }

    return tx.Commit()
}
