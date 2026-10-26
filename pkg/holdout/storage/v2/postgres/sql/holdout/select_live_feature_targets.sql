SELECT
    hg.id,
    hg.name,
    hg.status,
    hgf.feature_id,
    hgf.holdout_variation_id
FROM
    holdout_group_feature hgf
    INNER JOIN holdout_group hg
        ON hg.id = hgf.holdout_group_id
        AND hg.environment_id = hgf.environment_id
WHERE
    hgf.environment_id = $1
    AND hgf.feature_id = $2
    AND hgf.removed_at = 0
    AND hg.deleted = $3
    AND hg.status IN ($4, $5, $6)
