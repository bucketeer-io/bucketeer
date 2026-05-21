DELETE
FROM flag_trigger
WHERE id = $1
  AND environment_id = $2
