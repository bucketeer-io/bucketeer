UPDATE
    goal
SET
    name = ?,
    description = ?,
    archived = ?,
    deleted = ?,
    created_at = ?,
    updated_at = ?
WHERE
    id = ? AND
    environment_id = ?