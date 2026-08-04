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
    id = $3 AND
    status = $4 AND
    deleted = false
ON CONFLICT (notification_id, email) DO NOTHING
