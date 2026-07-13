-- Create tables for the system notification center (RFC 0047).
-- notification: system-admin-authored announcements (draft/published)
-- notification_localization: per-language tags, title, and Markdown content
-- notification_read: per-user read markers; unread = published notification
-- without a marker row for the viewer's email

CREATE TABLE IF NOT EXISTS notification (
    id VARCHAR(255) NOT NULL,                -- UUID
    status INT NOT NULL DEFAULT 0,           -- 0: DRAFT, 1: PUBLISHED
    created_by VARCHAR(255) NOT NULL,        -- editor email
    last_edited_by VARCHAR(255) NOT NULL,
    published_by VARCHAR(255) DEFAULT NULL,
    published_at BIGINT NOT NULL DEFAULT 0,  -- epoch seconds; 0 while draft
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_status_published_at (status, published_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS notification_localization (
    notification_id VARCHAR(255) NOT NULL,
    language VARCHAR(10) NOT NULL,           -- BCP 47 code: 'en', 'ja'
    tags JSON DEFAULT NULL,                  -- [{"name": "Announcement", "color": "#3B82F6"}]
    title VARCHAR(511) NOT NULL,
    content MEDIUMTEXT NOT NULL,             -- Markdown source
    PRIMARY KEY (notification_id, language),
    CONSTRAINT fk_notification_localization
        FOREIGN KEY (notification_id)
        REFERENCES notification (id)
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS notification_read (
    notification_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,             -- viewer identity (global across orgs)
    read_at BIGINT NOT NULL,
    PRIMARY KEY (notification_id, email),
    INDEX idx_email (email),
    CONSTRAINT fk_notification_read
        FOREIGN KEY (notification_id)
        REFERENCES notification (id)
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
