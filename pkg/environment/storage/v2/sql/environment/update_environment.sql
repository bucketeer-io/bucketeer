UPDATE environment_v2
SET
    name = ?,
    description = ?,
    archived = ?,
    require_comment = ?,
    created_at = ?,
    updated_at = ?
WHERE
    id = ?