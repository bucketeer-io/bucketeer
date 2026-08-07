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
    status = ? AND
    deleted = false
