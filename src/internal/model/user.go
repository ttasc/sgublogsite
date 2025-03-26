package model

import (
	"database/sql"

	"github.com/ttasc/sgublogsite/src/internal/model/repos"
)

func (m *Model) GetUserByID(id int32) (repos.GetUserByIDRow, error) {
    return m.query.GetUserByID(m.ctx, id)
}

func (m *Model) GetUserByEmailOrPhone(emailorphone string) (repos.GetUserByEmailOrPhoneRow, error) {
    return m.query.GetUserByEmailOrPhone(
        m.ctx,
        repos.GetUserByEmailOrPhoneParams{
            Email: emailorphone,
            Phone: emailorphone,
        },
    )
}

func (m *Model) SearchUsers(text string) ([]repos.FindUsersRow, error) {
    wildcard := "%" + text + "%"
    return m.query.FindUsers(m.ctx, wildcard)
}

func (m *Model) GetUsers() ([]repos.GetAllUsersRow, error) {
    return m.query.GetAllUsers(m.ctx)
}

func (m *Model) GetUserAvatarID(id int32) (int32, error) {
    res, err := m.query.GetUserAvatarID(m.ctx, id)
    if err != nil {
        return 0, err
    }
    return res.Int32, nil
}

func (m *Model) AddUser(firstname, lastname, phone, email, password string, role repos.UsersRole) (int32, error) {
    tx, err := m.DB.Begin()
    if err != nil {
        return 0, err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    res, err := qtx.AddUser(m.ctx, repos.AddUserParams{
        Firstname:      firstname,
        Lastname:       lastname,
        Phone:          phone,
        Email:          email,
        Password:       password,
        Role:           role,
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

func (m *Model) UpdateUserInfo(id int32, firstname, lastname, phone, email string) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdateUserInfo(m.ctx, repos.UpdateUserInfoParams{
        UserID:         id,
        Firstname:      firstname,
        Lastname:       lastname,
        Phone:          phone,
        Email:          email,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) UpdateUserAvatarID(userID, imageID int32) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdateUserAvatar(m.ctx, repos.UpdateUserAvatarParams{
        UserID:         userID,
        AvatarID:       sql.NullInt32{Int32: imageID, Valid: imageID > 0},
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) UpdateUserPassword(userID int32, password string) error {
    tx, err := m.DB.Begin()
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
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdateUserRole(m.ctx, repos.UpdateUserRoleParams{
        UserID:         userID,
        Role:           repos.UsersRole(role),
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) DeleteUser(id int32) error {
    tx, err := m.DB.Begin()
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
