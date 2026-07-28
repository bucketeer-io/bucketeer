-- Soft-delete support for the notification center: deleted notifications are
-- flagged instead of removed so their localizations and read markers stay
-- intact until the admin audit log is implemented.
ALTER TABLE notification ADD COLUMN deleted BOOLEAN NOT NULL DEFAULT FALSE;
