CREATE TABLE scheduled_flag_update (
    id VARCHAR(255) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    scheduled_at BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    changes JSON NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_feature
      FOREIGN KEY (feature_id, environment_id)
      REFERENCES feature(id, environment_id)
      ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
