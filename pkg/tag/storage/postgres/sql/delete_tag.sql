DELETE FROM
  tag
WHERE
  id = $1
  AND environment_id = $2
