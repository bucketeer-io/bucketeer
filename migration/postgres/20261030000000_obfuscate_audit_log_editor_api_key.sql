-- Obfuscate the public API key that the audit log editor carried in clear text,
-- keeping its first and last 4 characters like the API key entity data.
-- This runs in the pre-upgrade hook, so subscriber pods still on the previous
-- version may insert raw keys until they are replaced. The statement is
-- idempotent: run it again by hand once every subscriber is upgraded.
UPDATE audit_log
SET editor = jsonb_set(
  editor,
  '{public_api_editor,token}',
  to_jsonb(
    left(editor #>> '{public_api_editor,token}', 4)
    || '....'
    || right(editor #>> '{public_api_editor,token}', 4)
  )
)
WHERE length(editor #>> '{public_api_editor,token}') > 8;
