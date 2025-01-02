SELECT
    id,
    fcm_service_account,
    tags,
    deleted,
    name,
    created_at,
    updated_at,
    disabled
FROM
    push
WHERE
    id = ? AND
    environment_id = ? 