SELECT
    id,
    feature_id,
    file_path,
    line_number,
    code_snippet,
    content_hash,
    aliases,
    repository_name,
    repository_owner,
    CASE repository_type
        WHEN 'GITHUB' THEN 1
        WHEN 'GITLAB' THEN 2
        WHEN 'BITBUCKET' THEN 3
        WHEN 'CUSTOM' THEN 4
        ELSE 0
    END as repository_type,
    repository_branch,
    commit_hash,
    environment_id,
    created_at,
    updated_at
FROM
    code_reference 