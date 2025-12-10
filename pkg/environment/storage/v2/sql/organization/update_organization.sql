UPDATE
    organization
SET
    name = ?,
    owner_email = ?,
    description = ?,
    disabled = ?,
    archived = ?,
    trial = ?,
    created_at = ?,
    updated_at = ?
WHERE
    id = ?