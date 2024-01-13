UPDATE
    organization
SET
    name = ?,
    description = ?,
    disabled = ?,
    archived = ?,
    trial = ?,
    created_at = ?,
    updated_at = ?
WHERE
    id = ?