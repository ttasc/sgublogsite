package services

import "sgublogsite/src/internal/model/repos"

func GetImageByID(id int32) (repos.Image, error) {
    return s.query.GetImageByID(s.ctx, id)
}

func GetImages() ([]repos.Image, error) {
    return s.query.GetAllImages(s.ctx)
}

func AddImage(image repos.Image) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.AddImage(s.ctx, repos.AddImageParams{
        Url:          image.Url,
        Name:         image.Name,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func DeleteImage(id int32) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.DeleteImage(s.ctx, id)

    if err != nil {
        return err
    }

    return tx.Commit()
}
