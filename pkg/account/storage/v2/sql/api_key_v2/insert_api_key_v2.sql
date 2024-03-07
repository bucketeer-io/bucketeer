INSERT INTO api_key(
    id,
    name,
    role, disabled,
    created_at,
    updated_at,
    environment_namespace)
VALUES (?, ?, ?, ?, ?, ?, ?)
