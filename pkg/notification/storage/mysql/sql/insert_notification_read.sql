INSERT IGNORE INTO notification_read (
    notification_id,
    email,
    read_at
)
SELECT
    id,
    ?,
    ?
FROM
    notification
WHERE
    id = ? AND
    status = ? AND
    deleted = false
