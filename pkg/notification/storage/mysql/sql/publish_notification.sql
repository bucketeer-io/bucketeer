UPDATE notification
SET
    status = ?,
    published_by = ?,
    published_at = ?,
    updated_at = ?
WHERE
    id = ? AND
    deleted = false
