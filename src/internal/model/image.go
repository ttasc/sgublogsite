package model

import (
	"database/sql"

	"github.com/ttasc/sgublogsite/src/internal/model/repos"
)

func (m *Model) CountImages() (int64, error) {
    return m.query.CountImages(m.ctx)
}

func (m *Model) GetImageByID(id int32) (repos.Image, error) {
    return m.query.GetImageByID(m.ctx, id)
}

func (m *Model) GetImageByURL(url string) (repos.Image, error) {
    return m.query.GetImageByURL(m.ctx, url)
}

func (m *Model) GetImages(limit, offset int32) ([]repos.Image, error) {
    return m.query.GetAllImages(m.ctx, repos.GetAllImagesParams{
        Limit:          limit,
        Offset:         offset,
    })
}

func (m *Model) AddImage(name, url string) (int32, error) {
    tx, err := m.DB.Begin()
    if err != nil {
        return 0, err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    res, err := qtx.AddImage(m.ctx, repos.AddImageParams{
        Url:          url,
        Name:         sql.NullString{String: name, Valid: name != ""},
    })

    if err != nil {
        return 0, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }

    return int32(id), tx.Commit()
}

func (m *Model) UpdateImageURL(id int32, url string) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdateImageURL(m.ctx, repos.UpdateImageURLParams{
        ImageID:        id,
        Url:            url,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) DeleteImage(id int32) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.DeleteImage(m.ctx, id)

    if err != nil {
        return err
    }

    return tx.Commit()
}
