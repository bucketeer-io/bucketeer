DELETE FROM
    push
WHERE
    push.id = $1 AND
    push.environment_id = $2
