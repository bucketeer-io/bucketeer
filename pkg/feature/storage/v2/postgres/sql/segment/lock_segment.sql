SELECT
    id
FROM
    segment
WHERE
    id = $1 AND
    environment_id = $2
FOR UPDATE
