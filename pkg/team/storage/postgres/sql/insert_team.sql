INSERT INTO team (
    id,
    name,
    description,
    organization_id,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) ON CONFLICT (name, organization_id) DO UPDATE SET
    updated_at = EXCLUDED.updated_at
