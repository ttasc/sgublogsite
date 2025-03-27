package model

import (
	"database/sql"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
)

type GetCategory repos.Category

func (m *Model) GetCategoryByID(id int32) (repos.Category, error) {
    return m.query.GetCategoryByID(m.ctx, id)
}

func (m *Model) GetCategories() ([]repos.Category, error) {
    categories, err := m.query.GetAllCategories(m.ctx)
    if err != nil {
        return nil, err
    }
    uncategorized := repos.Category{
        CategoryID: -1,
        ParentCategoryID: sql.NullInt32{Int32: 0, Valid: false},
        Name: "Uncategorized",
        Slug: "uncategorized",
    }
    categories = append(categories, uncategorized)
    return categories, nil
}

func (m *Model) GetChildCategories(id int32) ([]repos.Category, error) {
    return m.query.GetChildCategories(m.ctx, sql.NullInt32{Int32: id, Valid: id > 0})
}

func (m *Model) GetParentCategoryID(id int32) (int32, error) {
    res, err := m.query.GetParentCategoryID(m.ctx, id)
    if !res.Valid {
        return 0, err
    }
    return res.Int32, nil
}

func (m *Model) AddCategory(pid int32, name, slug string) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.AddCategory(m.ctx, repos.AddCategoryParams{
        ParentCategoryID:   sql.NullInt32{Int32: pid, Valid: pid > 0},
        Name:               name,
        Slug:               slug,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) UpdateCategory(id int32, name, slug string) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdateCategory(m.ctx, repos.UpdateCategoryParams{
        CategoryID:         id,
        Name:               name,
        Slug:               slug,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) UpdateCategoryParent(id, pid int32) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdateCategoryParent(m.ctx, repos.UpdateCategoryParentParams{
        CategoryID:         id,
        ParentCategoryID:   sql.NullInt32{Int32: pid, Valid: pid > 0},
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) DeleteCategory(id int32) error {
    tx, err := m.DB.Begin()
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
