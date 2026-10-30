-- SDK API keys are matched on a SHA-256 digest of the key instead of the key
-- itself, and a short prefix is kept for display. Both columns stay nullable so
-- existing rows keep working until they are backfilled; the api_key column is
-- still written and read as a fallback in the meantime.
ALTER TABLE api_key ADD COLUMN api_key_hash CHAR(64), ADD COLUMN api_key_prefix VARCHAR(12);
CREATE UNIQUE INDEX unique_api_key_hash ON api_key (api_key_hash);
