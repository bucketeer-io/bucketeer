INSERT INTO organization (
    id,
    name,
    owner_email,
    url_code,
    description,
    disabled,
    archived,
    trial,
    system_admin,
    password_authentication_enabled,
    created_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)