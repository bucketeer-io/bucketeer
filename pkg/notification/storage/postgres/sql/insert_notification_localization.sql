INSERT INTO notification_localization (
    notification_id,
    language,
    tags,
    title,
    content
) VALUES (
    $1, $2, $3, $4, $5
)
