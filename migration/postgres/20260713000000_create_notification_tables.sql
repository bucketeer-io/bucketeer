-- Create tables for the system notification center (RFC 0047).
-- notification: system-admin-authored announcements (draft/published)
-- notification_localization: per-language tags, title, and Markdown content
-- notification_read: per-user read markers; unread = published notification
-- without a marker row for the viewer's email

CREATE TABLE notification (
    id VARCHAR(255) NOT NULL,                -- UUID
    status INT NOT NULL DEFAULT 0,           -- 0: DRAFT, 1: PUBLISHED
    created_by VARCHAR(255) NOT NULL,        -- editor email
    last_edited_by VARCHAR(255) NOT NULL,
    published_by VARCHAR(255) DEFAULT NULL,
    published_at BIGINT NOT NULL DEFAULT 0,  -- epoch seconds; 0 while draft
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX idx_notification_status_published_at ON notification (status, published_at);

CREATE TABLE notification_localization (
    notification_id VARCHAR(255) NOT NULL,
    language VARCHAR(10) NOT NULL,           -- BCP 47 code: 'en', 'ja'
    tags JSONB DEFAULT NULL,                 -- [{"name": "Announcement", "color": "#3B82F6"}]
    title VARCHAR(511) NOT NULL,
    content TEXT NOT NULL,                   -- Markdown source
    PRIMARY KEY (notification_id, language),
    CONSTRAINT fk_notification_localization
        FOREIGN KEY (notification_id)
        REFERENCES notification (id)
        ON DELETE CASCADE
);

CREATE TABLE notification_read (
    notification_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,             -- viewer identity (global across orgs)
    read_at BIGINT NOT NULL,
    PRIMARY KEY (notification_id, email),
    CONSTRAINT fk_notification_read
        FOREIGN KEY (notification_id)
        REFERENCES notification (id)
        ON DELETE CASCADE
);
CREATE INDEX idx_notification_read_email ON notification_read (email);
