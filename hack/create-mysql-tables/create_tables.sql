-- Copyright 2025 The Bucketeer Authors.
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS bucketeer;

-- Use the bucketeer database
USE bucketeer;

-- Evaluation Event Table
CREATE TABLE IF NOT EXISTS evaluation_event (
    id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    timestamp DATETIME(6) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    feature_version INT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    user_data JSON,
    variation_id VARCHAR(255) NOT NULL,
    reason TEXT,
    tag VARCHAR(255),
    source_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_environment_id (environment_id),
    INDEX idx_timestamp (timestamp),
    INDEX idx_feature_id (feature_id),
    INDEX idx_user_id (user_id),
    INDEX idx_variation_id (variation_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Goal Event Table
CREATE TABLE IF NOT EXISTS goal_event (
    id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    timestamp DATETIME(6) NOT NULL,
    goal_id VARCHAR(255) NOT NULL,
    value FLOAT,
    user_id VARCHAR(255) NOT NULL,
    user_data JSON,
    tag VARCHAR(255),
    source_id VARCHAR(255),
    feature_id VARCHAR(255),
    feature_version INT,
    variation_id VARCHAR(255),
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_environment_id (environment_id),
    INDEX idx_timestamp (timestamp),
    INDEX idx_goal_id (goal_id),
    INDEX idx_user_id (user_id),
    INDEX idx_feature_id (feature_id),
    INDEX idx_variation_id (variation_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci; 