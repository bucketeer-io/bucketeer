UPDATE
    organization
SET
    name = ?,
    owner_email = ?,
    description = ?,
    disabled = ?,
    archived = ?,
    trial = ?,
    password_authentication_enabled = ?,
    created_at = ?,
    updated_at = ?
WHERE
    id = ?