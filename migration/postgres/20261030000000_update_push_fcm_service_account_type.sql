-- Modify "push" table: fcm_service_account now holds either the legacy JSON
-- document or its encrypted form ("enc:v1:..."), so it is stored as text.
-- Existing rows keep their JSON content as text and stay readable.
ALTER TABLE push ALTER COLUMN fcm_service_account TYPE TEXT USING fcm_service_account::text;
