UPDATE
    account_v2
SET
    search_filters = ?
WHERE
    email = ?
  AND organization_id = ?
