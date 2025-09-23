UPDATE
    organization
SET
    name = ?,
    owner_email = ?,
    description = ?,
    disabled = ?,
    archived = ?,
    trial = ?,
    authentication_settings = ?,
    created_at = ?,
    updated_at = ?
WHERE
    id = ?