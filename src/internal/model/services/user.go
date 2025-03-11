package services

import "sgublogsite/src/internal/model/repos"

func GetUserByID(id int32) (repos.User, error) {
    return s.query.GetUserByID(s.ctx, id)
}

func SearchUsers(text string) ([]repos.User, error) {
    wildcard := "%" + text + "%"
    return s.query.FindUsers(s.ctx, wildcard)
}

func GetUsers() ([]repos.User, error) {
    return s.query.GetAllUsers(s.ctx)
}

func CreateUser(user repos.User) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.CreateUser(s.ctx, repos.CreateUserParams{
        Firstname:      user.Firstname,
        Lastname:       user.Lastname,
        Mobile:         user.Mobile,
        Email:          user.Email,
        Password:       user.Password,
        ProfilePicID:   user.ProfilePicID,
        Role:           user.Role,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func UpdateUserInfo(user repos.User) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.UpdateUserInfo(s.ctx, repos.UpdateUserInfoParams{
        Firstname:      user.Firstname,
        Lastname:       user.Lastname,
        Mobile:         user.Mobile,
        Email:          user.Email,
        ProfilePicID:   user.ProfilePicID,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func UpdateUserPassword(user repos.User) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.UpdateUserPassword(s.ctx, repos.UpdateUserPasswordParams{
        Password:       user.Password,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func UpdateUserRole(user repos.User) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.UpdateUserRole(s.ctx, repos.UpdateUserRoleParams{
        Role:           user.Role,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func DeleteUser(id int32) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := s.query.WithTx(tx)

    _, err = qtx.DeleteUser(s.ctx, id)

    if err != nil {
        return err
    }

    return tx.Commit()
}
