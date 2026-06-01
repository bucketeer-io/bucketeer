INSERT INTO code_reference (
    id,
    feature_id,
    file_path,
    file_extension,
    line_number,
    code_snippet,
    content_hash,
    aliases,
    repository_name,
    repository_owner,
    repository_type,
    repository_branch,
    commit_hash,
    environment_id,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
