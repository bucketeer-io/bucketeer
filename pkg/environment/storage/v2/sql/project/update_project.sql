UPDATE project
SET
    name = ?,
    description = ?,
    disabled = ?,
    trial = ?,
    creator_email = ?,
    created_at = ?,
    updated_at = ?
WHERE
    id = ?
