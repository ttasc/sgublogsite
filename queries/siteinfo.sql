-- name: GetSiteInfo :one
SELECT * FROM siteinfo;

-- name: GetSiteAbout :one
SELECT site_about FROM siteinfo;

-- name: GetContactInfo :one
SELECT contact_address, contact_email, contact_phone FROM siteinfo;

-- name: GetSiteMeta :one
SELECT
    site_title,
    site_name,
    images.url,
    site_copyright
FROM siteinfo JOIN images
ON siteinfo.site_logo_id = images.image_id;

-- name: CreateSiteInfo :execresult
INSERT INTO siteinfo (
    site_title,
    site_logo_id,
    site_name,
    site_about,
    site_copyright,
    contact_address,
    contact_email,
    contact_phone
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateSiteInfo :execresult
UPDATE siteinfo SET
    site_title = ?,
    site_logo_id = ?,
    site_name = ?,
    site_about = ?,
    site_copyright = ?,
    contact_address = ?,
    contact_email = ?,
    contact_phone = ?
WHERE site_id = 1;
