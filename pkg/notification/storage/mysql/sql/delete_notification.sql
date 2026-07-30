UPDATE notification
SET
    deleted = true,
    last_edited_by = ?,
    updated_at = ?
WHERE
    id = ? AND
    deleted = false
