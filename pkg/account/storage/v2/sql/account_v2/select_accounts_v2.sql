SELECT 
    email,
    name,
    first_name,
    last_name,
    language,
    avatar_image_url,
    avatar_file_type,
    avatar_image,
    organization_id,
    organization_role,
    environment_roles,
    disabled,
    created_at,
    updated_at,
    last_seen,
    search_filters,
    JSON_LENGTH(environment_roles) as environment_count
FROM account_v2
%s  
%s  
%s  