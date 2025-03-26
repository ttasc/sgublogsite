package model

import "github.com/ttasc/sgublogsite/src/internal/model/repos"

func (m *Model) GetTagByID(id int32) (repos.Tag, error) {
    return m.query.GetTagByID(m.ctx, id)
}

func (m *Model) GetTags() ([]repos.Tag, error) {
    return m.query.GetAllTags(m.ctx)
}

func (m *Model) GetTagNames() ([]string, error) {
    return m.query.GetAllTagNames(m.ctx)
}

func (m *Model) GetTagsByPostID(id int32) ([]repos.Tag, error) {
    return m.query.GetTagsByPostID(m.ctx, id)
}

func (m *Model) AddTag(name, slug string) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.AddTag(m.ctx, repos.AddTagParams{
        Name:           name,
        Slug:           slug,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) DeleteTag(id int32) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.DeleteTag(m.ctx, id)

    if err != nil {
        return err
    }

    return tx.Commit()
}
