-- Obfuscate the public API key that the audit log editor carried in clear text,
-- keeping its first and last 4 characters like the API key entity data.
-- This runs in the pre-upgrade hook, so subscriber pods still on the previous
-- version may insert raw keys until they are replaced. The statement is
-- idempotent: run it again by hand once every subscriber is upgraded.
UPDATE audit_log
SET editor = JSON_SET(
  editor,
  '$.public_api_editor.token',
  CONCAT(
    LEFT(editor->>'$.public_api_editor.token', 4),
    '....',
    RIGHT(editor->>'$.public_api_editor.token', 4)
  )
)
WHERE CHAR_LENGTH(editor->>'$.public_api_editor.token') > 8;
