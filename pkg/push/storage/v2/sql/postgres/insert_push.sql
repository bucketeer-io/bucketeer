INSERT INTO push (
    id,
    fcm_service_account,
    tags,
    deleted,
    name,
    created_at,
    updated_at,
    environment_id,
    disabled
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 