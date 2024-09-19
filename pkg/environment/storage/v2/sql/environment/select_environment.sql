SELECT
    id,
    name,
    url_code,
    description,
    project_id,
    organization_id,
    archived,
    require_comment,
    created_at,
    updated_at
FROM
    environment_v2
WHERE
    id = ?