DELETE FROM
    ops_progressive_rollout
WHERE
    id = $1 AND
    environment_id = $2
