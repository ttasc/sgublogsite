package model

import "github.com/ttasc/sgublogsite/src/internal/model/repos"

func (m *Model) GetImageByID(id int32) (repos.Image, error) {
    return m.query.GetImageByID(m.ctx, id)
}

func (m *Model) GetImageByURL(url string) (repos.Image, error) {
    return m.query.GetImageByURL(m.ctx, url)
}

func (m *Model) GetImages() ([]repos.Image, error) {
    return m.query.GetAllImages(m.ctx)
}

func (m *Model) AddImage(image repos.Image) (int32, error) {
    tx, err := m.DB.Begin()
    if err != nil {
        return 0, err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    res, err := qtx.AddImage(m.ctx, repos.AddImageParams{
        Url:          image.Url,
        Name:         image.Name,
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
