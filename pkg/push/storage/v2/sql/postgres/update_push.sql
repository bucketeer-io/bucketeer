UPDATE push
SET
    fcm_service_account = $1,
    tags = $2,
    deleted = $3,
    name = $4,
    created_at = $5,
    updated_at = $6,
    disabled = $7
WHERE
    id = $8 AND
    environment_id = $9 