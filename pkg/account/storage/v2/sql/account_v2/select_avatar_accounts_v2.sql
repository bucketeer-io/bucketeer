SELECT
    a.email,
    a.avatar_file_type,
    a.avatar_image
FROM account_v2 AS a
INNER JOIN environment_v2 AS e
ON a.organization_id = e.organization_id
%s
