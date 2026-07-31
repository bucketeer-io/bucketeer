SELECT DISTINCT
    notification.id,
    notification.status,
    notification.created_by,
    notification.last_edited_by,
    COALESCE(notification.published_by, ''),
    notification.published_at,
    notification.created_at,
    notification.updated_at
FROM
    notification
LEFT JOIN
    notification_localization ON notification.id = notification_localization.notification_id
