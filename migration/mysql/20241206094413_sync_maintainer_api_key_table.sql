-- Get creator of api_key from audit_log and update api_key maintainer
UPDATE api_key
    JOIN audit_log ON api_key.id = audit_log.entity_id
    SET api_key.maintainer = (JSON_UNQUOTE(JSON_EXTRACT(editor, '$.email')))
WHERE audit_log.type = 400;
