SELECT
    id,
    name,
    url_code,
    description,
    disabled,
    trial,
    creator_email,
    organization_id,
    created_at,
    updated_at
FROM
    project
%s %s %s
