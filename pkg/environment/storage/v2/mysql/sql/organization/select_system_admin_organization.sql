SELECT
    id,
    name,
    owner_email,
    url_code,
    description,
    disabled,
    archived,
    trial,
    system_admin,
    authentication_settings,
    created_at,
    updated_at
FROM
    organization
WHERE
    system_admin = 1