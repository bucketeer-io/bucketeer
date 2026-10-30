-- Obfuscate the public API key that the audit log editor carried in clear text,
-- keeping its first and last 4 characters like the API key entity data.
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
