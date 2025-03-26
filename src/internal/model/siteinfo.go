package model

import (
	"database/sql"

	"github.com/ttasc/sgublogsite/src/internal/model/repos"
)

func (m *Model) GetSiteInfo() (repos.Siteinfo, error) {
    return m.query.GetSiteInfo(m.ctx)
}

func (m *Model) GetSiteAbout() (string, error) {
    about, err := m.query.GetSiteAbout(m.ctx)
    if err != nil {
        return "", err
    }
    return about.String, nil
}

func (m *Model) GetContactInfo() (repos.GetContactInfoRow, error) {
    return m.query.GetContactInfo(m.ctx)
}

func (m *Model) GetSiteMeta() (repos.GetSiteMetaRow, error) {
    return m.query.GetSiteMeta(m.ctx)
}

func (m *Model) CreateSiteInfo(title string, logoID int32, name, about, copyright, address, email, phone string) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.CreateSiteInfo(m.ctx, repos.CreateSiteInfoParams{
        SiteTitle:      sql.NullString{String: title    , Valid: title != ""},
        SiteLogoID:     sql.NullInt32 {Int32:  logoID   , Valid: logoID > 0},
        SiteName:       sql.NullString{String: name     , Valid: name != "" },
        SiteAbout:      sql.NullString{String: about    , Valid: about != "" },
        SiteCopyright:  sql.NullString{String: copyright, Valid: copyright != "" },
        ContactAddress: sql.NullString{String: address  , Valid: address != "" },
        ContactEmail:   sql.NullString{String: email    , Valid: email != "" },
        ContactPhone:   sql.NullString{String: phone    , Valid: phone != "" },
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) UpdateSiteInfo(title string, logoID int32, name, about, copyright, address, email, phone string) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdateSiteInfo(m.ctx, repos.UpdateSiteInfoParams{
        SiteTitle:      sql.NullString{String: title    , Valid: title != ""},
        SiteLogoID:     sql.NullInt32 {Int32:  logoID   , Valid: logoID > 0},
        SiteName:       sql.NullString{String: name     , Valid: name != "" },
        SiteAbout:      sql.NullString{String: about    , Valid: about != "" },
        SiteCopyright:  sql.NullString{String: copyright, Valid: copyright != "" },
        ContactAddress: sql.NullString{String: address  , Valid: address != "" },
        ContactEmail:   sql.NullString{String: email    , Valid: email != "" },
        ContactPhone:   sql.NullString{String: phone    , Valid: phone != "" },
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}
