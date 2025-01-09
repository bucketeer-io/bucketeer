SELECT 
    id,
    created_at,
    updated_at,
    entity_type,
    environment_id
FROM
    tag
%s %s %s
