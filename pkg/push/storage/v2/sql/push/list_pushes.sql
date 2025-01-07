SELECT
    id,
    fcm_service_account,
    tags,
    deleted,
    name,
    created_at,
    updated_at,
    disabled,
    environment_id
FROM
    push
%s %s %s 