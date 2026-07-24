UPDATE notification
SET
    last_edited_by = $1,
    updated_at = $2
WHERE
    id = $3
