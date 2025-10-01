SELECT
    api_key_id,
    environment_id,
    last_used_at,
    created_at
FROM
    api_key_last_used_info
%s