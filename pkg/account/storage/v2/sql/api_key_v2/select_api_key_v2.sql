SELECT
    id,
    name,
    role,
    disabled,
    created_at,
    updated_at
FROM
    api_key
    %s %s %s