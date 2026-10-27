SELECT
    id
FROM
    segment
WHERE
    id = ? AND
    environment_id = ?
FOR UPDATE
