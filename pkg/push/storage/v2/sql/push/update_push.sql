UPDATE push
SET
    fcm_service_account = ?,
    tags = ?,
    deleted = ?,
    name = ?,
    created_at = ?,
    updated_at = ?,
    disabled = ?
WHERE
    id = ? AND
    environment_id = ? 