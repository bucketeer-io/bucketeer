CREATE TABLE code_reference (
    id varchar(255) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    file_path VARCHAR(512) NOT NULL,
    line_number INT NOT NULL,
    code_snippet TEXT NOT NULL,
    content_hash VARCHAR(64) NOT NULL,
    aliases JSON,
    repository_name VARCHAR(255) NOT NULL,
    repository_owner VARCHAR(255) NOT NULL,
    repository_type TINYINT NOT NULL COMMENT '0:UNSPECIFIED, 1:GITHUB, 2:GITLAB, 3:BITBUCKET, 4:CUSTOM',
    repository_branch VARCHAR(255) NOT NULL,
    commit_hash VARCHAR(40) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_file_path (file_path),
    CONSTRAINT foreign_code_references_feature FOREIGN KEY (feature_id, environment_id) 
        REFERENCES feature (id, environment_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
