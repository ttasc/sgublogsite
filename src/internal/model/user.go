package model

import "github.com/ttasc/sgublogsite/src/internal/model/repos"

func (m *Model) GetUserByID(id int32) (repos.GetUserByIDRow, error) {
    return m.query.GetUserByID(m.ctx, id)
}

func (m *Model) GetUserByEmailOrPhone(emailorphone string) (repos.GetUserByEmailOrPhoneRow, error) {
    return m.query.GetUserByEmailOrPhone(
        m.ctx,
        repos.GetUserByEmailOrPhoneParams{
            Email: emailorphone,
            Phone: emailorphone,
        })
}

func (m *Model) SearchUsers(text string) ([]repos.FindUsersRow, error) {
    wildcard := "%" + text + "%"
    return m.query.FindUsers(m.ctx, wildcard)
}

func (m *Model) GetUsers() ([]repos.GetAllUsersRow, error) {
    return m.query.GetAllUsers(m.ctx)
}

func (m *Model) AddUser(user repos.User) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.AddUser(m.ctx, repos.AddUserParams{
        Firstname:      user.Firstname,
        Lastname:       user.Lastname,
        Phone:         user.Phone,
        Email:          user.Email,
        Password:       user.Password,
        Role:           user.Role,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) UpdateUserInfo(user repos.User) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdateUserInfo(m.ctx, repos.UpdateUserInfoParams{
        UserID:         user.UserID,
        Firstname:      user.Firstname,
        Lastname:       user.Lastname,
        Phone:         user.Phone,
        Email:          user.Email,
        ProfilePicID:   user.ProfilePicID,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) UpdateUserPassword(userID int32, password string) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdateUserPassword(m.ctx, repos.UpdateUserPasswordParams{
        UserID:         userID,
        Password:       password,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) UpdateUserRole(userID int32, role repos.UsersRole) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdateUserRole(m.ctx, repos.UpdateUserRoleParams{
        UserID:         userID,
        Role:           repos.NullUsersRole{ UsersRole: role, Valid:  true},
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) DeleteUser(id int32) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.DeleteUser(m.ctx, id)

    if err != nil {
        return err
    }

    return tx.Commit()
}
