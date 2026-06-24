WITH grouped_by_user_evaluation AS (
    SELECT
        user_id,
        variation_id,
        COUNT(id) as event_count,
        COALESCE(SUM(value), 0) as value_sum
    FROM
        goal_event
    WHERE
        "timestamp" BETWEEN $1 AND $2
        AND environment_id = $3
        AND goal_id = $4
        AND feature_id = $5
        AND feature_version = $6
    GROUP BY
        user_id,
        variation_id
),
cap_level AS (
    -- Winsorization threshold: the configurable percentile ($7, a fraction in
    -- [0,1]) of per-user value_sum pooled across all variations.
    -- PERCENTILE_CONT is an ordered-set aggregate (SQL:2003); COALESCE guards
    -- the empty-CTE case. Values above this threshold are whales whose outsized
    -- spend would otherwise make the Normal posterior overconfident (Kohavi,
    -- Tang & Xu, Trustworthy Online Controlled Experiments, 2020, §4).
    SELECT COALESCE(
        PERCENTILE_CONT($7) WITHIN GROUP (ORDER BY value_sum),
        0
    ) AS cap
    FROM grouped_by_user_evaluation
),
capped_by_user AS (
    SELECT
        u.user_id,
        u.variation_id,
        u.event_count,
        LEAST(u.value_sum, c.cap) AS value_sum
    FROM grouped_by_user_evaluation u
    CROSS JOIN cap_level c
)
SELECT
    variation_id as variationID,
    COUNT(DISTINCT user_id) as goalUser,
    SUM(event_count) as goalTotal,
    SUM(value_sum) as goalValueTotal,
    AVG(value_sum) as goalValueMean,
    COALESCE(VAR_SAMP(value_sum), 0) as goalValueVariance
FROM
    capped_by_user
GROUP BY
    variation_id
ORDER BY
    variation_id
