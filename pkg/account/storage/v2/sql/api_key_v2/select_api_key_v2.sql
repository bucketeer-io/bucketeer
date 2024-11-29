SELECT
    api_key.id,
    api_key.name,
    api_key.role,
    api_key.disabled,
    api_key.created_at,
    api_key.updated_at,
    api_key.description,
    environment_v2.name as environment_name
FROM
    api_key
LEFT JOIN environment_v2 ON api_key.environment_id = environment_v2.id
%s %s %s