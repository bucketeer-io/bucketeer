DELETE FROM
    push
WHERE
    push.id = ? AND
    push.environment_id = ?
