SELECT 
    ak.id AS api_key_id,
    ak.name AS api_key_name,
    ak.role AS api_key_role,
    ak.disabled AS api_key_disabled,
    ak.created_at AS api_key_created_at,
    ak.updated_at AS api_key_updated_at,
    ak.description AS api_key_description,
    ak.api_key AS api_key_key,
    ak.maintainer AS api_key_maintainer,

    env.id AS environment_id,
    env.name AS environment_name,
    env.url_code AS environment_url_code,
    env.description AS environment_description,
    env.project_id AS environment_project_id,
    env.organization_id AS environment_organization_id,
    env.archived AS environment_archived,
    env.require_comment AS environment_require_comment,
    env.created_at AS environment_created_at,
    env.updated_at AS environment_updated_at,

    proj.id AS project_id,
    proj.url_code AS project_url_code,
    proj.disabled AS project_disabled
FROM api_key ak
JOIN environment_v2 env ON ak.environment_id = env.id
JOIN project proj ON env.project_id = proj.id
