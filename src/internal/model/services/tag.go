package services

import "sgublogsite/src/internal/model/repos"

func GetTagByID(id int32) (repos.Tag, error) {
    return s.query.GetTagByID(s.ctx, id)
}

func GetTags() ([]repos.Tag, error) {
    return s.query.GetAllTags(s.ctx)
}

func AddTag(tag repos.Tag) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.AddTag(s.ctx, repos.AddTagParams{
        Name:           tag.Name,
        Slug:           tag.Slug,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func DeleteTag(id int32) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.DeleteTag(s.ctx, id)

    if err != nil {
        return err
    }

    return tx.Commit()
}
