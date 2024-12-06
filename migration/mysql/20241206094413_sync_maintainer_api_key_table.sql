-- Get creator of api_key from audit_log and update api_key maintainer
update api_key
    join audit_log on api_key.id = audit_log.entity_id
    set api_key.maintainer = (JSON_UNQUOTE(json_extract(editor, '$.email')))
where audit_log.type = 400;
