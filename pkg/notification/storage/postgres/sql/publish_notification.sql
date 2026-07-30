UPDATE notification
SET
    status = $1,
    published_by = $2,
    published_at = $3,
    updated_at = $4
WHERE
    id = $5 AND
    deleted = false
