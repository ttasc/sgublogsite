package model

import (
	"sgublogsite/src/internal/model/repos"
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

func (m *Model) CreateSiteInfo(info repos.Siteinfo) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.CreateSiteInfo(m.ctx, repos.CreateSiteInfoParams{
        SiteTitle:      info.SiteTitle,
        SiteLogoID:     info.SiteLogoID,
        SiteName:       info.SiteName,
        SiteAbout:      info.SiteAbout,
        SiteCopyright:  info.SiteCopyright,
        ContactAddress: info.ContactAddress,
        ContactEmail:   info.ContactEmail,
        ContactPhone:   info.ContactPhone,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}

func (m *Model) UpdateSiteInfo(info repos.Siteinfo) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := m.query.WithTx(tx)

    _, err = qtx.UpdateSiteInfo(m.ctx, repos.UpdateSiteInfoParams{
        SiteTitle:      info.SiteTitle,
        SiteLogoID:     info.SiteLogoID,
        SiteName:       info.SiteName,
        SiteAbout:      info.SiteAbout,
        SiteCopyright:  info.SiteCopyright,
        ContactAddress: info.ContactAddress,
        ContactEmail:   info.ContactEmail,
        ContactPhone:   info.ContactPhone,
    })

    if err != nil {
        return err
    }

    return tx.Commit()
}
