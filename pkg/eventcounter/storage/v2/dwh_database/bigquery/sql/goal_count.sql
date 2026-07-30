WITH deduped_events AS (
    -- BigQuery has no primary keys, so an at-least-once Pub/Sub redelivery
    -- can append the same event twice. Deduplicate by event ID before any
    -- aggregation, otherwise duplicates inflate event_count and value_sum.
    -- Duplicates are identical rows, so DISTINCT over the used columns is
    -- equivalent to picking one row per ID.
    SELECT DISTINCT
        id,
        user_id,
        variation_id,
        value
    FROM
        `%s`
    WHERE
        timestamp BETWEEN @startAt AND @endAt
    AND environment_id = @environmentId
    AND goal_id = @goalID
    AND feature_id = @featureID
    AND feature_version = @featureVersion
),
grouped_by_user_evaluation AS (
    SELECT
        user_id,
        variation_id,
        COUNT(id) as event_count,
        IFNULL(SUM(value), 0) as value_sum
    FROM
        deduped_events
    GROUP BY
        user_id,
        variation_id
),
cap_level AS (
    -- Winsorization threshold: the configurable percentile
    -- (@valueCapPercentile, an integer in [1,100]) of per-user value_sum pooled
    -- across all variations. APPROX_QUANTILES(value_sum, 100) returns an array
    -- of 101 elements (quantiles 0..100). We UNNEST WITH OFFSET and select the
    -- requested element via a WHERE predicate rather than @param inside
    -- OFFSET(...) — a bound parameter is not reliably accepted as an array
    -- subscript in a static query. IFNULL/MAX guards the empty case. Values
    -- above this threshold are whales whose outsized spend would otherwise make
    -- the Normal posterior overconfident (Kohavi, Tang & Xu, Trustworthy Online
    -- Controlled Experiments, 2020, §4).
    SELECT IFNULL(MAX(q), 0) AS cap
    FROM (
        SELECT APPROX_QUANTILES(value_sum, 100) AS quantiles
        FROM grouped_by_user_evaluation
    ),
    UNNEST(quantiles) AS q WITH OFFSET off
    WHERE off = @valueCapPercentile
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
    IFNULL(VAR_SAMP(value_sum), 0) as goalValueVariance
FROM
    capped_by_user
GROUP BY
    variation_id
ORDER BY
    variation_id
