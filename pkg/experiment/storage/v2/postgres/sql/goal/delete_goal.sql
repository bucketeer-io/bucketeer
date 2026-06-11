DELETE
FROM
    goal
WHERE
    id = $1 AND environment_id = $2
