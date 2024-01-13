SELECT
    id,
    name,
    url_code,
    description,
    disabled,
    archived,
    trial,
    system_admin,
    created_at,
    updated_at
FROM
    organization
WHERE
    id = ?