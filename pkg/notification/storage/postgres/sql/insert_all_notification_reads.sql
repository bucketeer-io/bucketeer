INSERT INTO notification_read (
    notification_id,
    email,
    read_at
)
SELECT
    id,
    $1,
    $2
FROM
    notification
WHERE
    status = $3 AND
    deleted = false
ON CONFLICT (notification_id, email) DO NOTHING
