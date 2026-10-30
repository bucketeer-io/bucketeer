-- Obfuscate the public API key that the audit log editor carried in clear text,
-- keeping its first and last 4 characters like the API key entity data.
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
