SELECT
    id,
    status,
    created_by,
    last_edited_by,
    COALESCE(published_by, ''),
    published_at,
    created_at,
    updated_at
FROM
    notification
WHERE
    id = ? AND
    deleted = false
