SELECT
    COUNT(DISTINCT notification.id)
FROM
    notification
LEFT JOIN
    notification_localization ON notification.id = notification_localization.notification_id
