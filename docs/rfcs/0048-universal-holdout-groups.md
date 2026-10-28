# Universal Holdout Groups

This RFC describes **Holdout Groups**: a long-lived control group that reserves a fixed share of users across *many* feature flags, so that the cumulative impact of every change shipped during a period (typically 3–6 months) can be measured against a population that never received any of them.

**Issue**: [https://github.com/bucketeer-io/bucketeer/issues/2798](https://github.com/bucketeer-io/bucketeer/issues/2798)

This revision follows the decisions recorded in [this comment](https://github.com/bucketeer-io/bucketeer/issues/2798#issuecomment-5724972912) on the issue.

## 1. Background

A single A/B test answers "did this change help?", not "did everything we shipped this half move the business metric?" — once a winner rolls out, the baseline is gone, and effects that decay or interact go unnoticed. The industry answer is a **holdout**: a fixed slice of users that keeps the *control* variation of every flag in a program for the whole period, while everyone else goes through normal targeting, rollouts and experiments.

[LaunchDarkly](https://launchdarkly.com/docs/home/holdouts), [GrowthBook](https://docs.growthbook.io/app/holdouts), Statsig and Eppo all ship this feature, and their implementations converge on the same core: a flag-independent hash with a program-scoped seed, a holdout decision taken ahead of other flag logic, and analysis reusing the existing experimentation pipeline. This design follows that convergence, and deliberately diverges where their implementations are weakest: in both LaunchDarkly and GrowthBook the holdout is materialized as a prerequisite-shaped reference carried on the flag side, so a flag editor can detach a running holdout and the kill switch silently overrides it — exactly the failure modes requirement 3 of the issue rules out (§6.1). Two of their ideas are adopted here instead of the more obvious designs: GrowthBook's **matched comparison sample** for measurement (§4.4) and its **analysis period** lifecycle stage (§4.6).

## 2. Goals / Non-Goals

**Goals**

- One Holdout Group deterministically decides which users are held out, and the **same** users are held out on every target flag.
- Membership is configurable: a hash-based percentage, or a segment reference.
- The Holdout Group owns its list of target flags. Flags are opted in explicitly; nothing else in the environment is affected.
- Holdout serving is authoritative: nothing on the flag side — targeting, rollout, auto-ops, the kill switch — can pull a user out of the holdout.
- Users outside the holdout are evaluated by the normal path with byte-identical results.
- Aggregate measurement: holdout vs. a **matched comparison sample** of equal size, over a set of goals, using the existing Bayesian experiment calculator (§4.4).
- **Measurement is optional per group.** A group with no goals configured serves the holdout while writing nothing to the analytics pipeline — no derived rows, no dedupe, no DWH cost. Teams that analyze in their own warehouse join on the served reason (`Reason.HOLDOUT` + `holdout_group_id`), which every evaluation event carries regardless (§4.4).
- Lifecycle: scheduled start/stop *and* manual activation/deactivation, an optional **analysis period** that freezes the target-flag list while continuing to hold users out (§4.6), adding and removing target flags mid-period, dissolution.
- Full audit trail; MySQL and PostgreSQL support.

**Non-Goals**

- Nested or hierarchical holdouts (a holdout inside a holdout).
- More than one Holdout Group serving at the same time in an environment (see §4.6 and §6.4). The hash design nonetheless reserves the ability to lift this later (§4.2.1).
- Cross-environment or cross-project holdouts. A Holdout Group is environment-scoped, like flags and experiments.
- Automatic enrollment of newly created flags. Opt-in stays explicit, but the detail page ships with a **coverage report** in v1 — flags created or first enabled during the period that are not targets of the group — so silent dilution of the readout is visible rather than discovered at the end (§4.7).

## 3. Design Overview

A Holdout Group is a **new first-class entity in the evaluation engine**, not a materialization into flag configuration (§6.1). Membership is a pure function of the group's own configuration, serving happens before any flag state is consulted, and measurement reuses the existing experiment pipeline by giving the group a synthetic feature identity.

The decisions below answer the open questions raised on the issue.

| # | Question | Decision |
|---|---|---|
| 1 | User assignment model | Configurable via `assignment_type`: hash-based percentage on an **environment-scoped seed with a per-group bucket range** (default), or a segment reference. `assignment_type`, `percentage`, `sampling_seed`, `bucket_start` and `segment_id` are all fixed at creation, so membership never changes for a given user (§4.2, §4.6). |
| 2 | Data model / evaluation strategy | New entities `holdout_group` and `holdout_group_feature`; target flags are managed through the intermediate table (§4.1, §4.3.1). |
| 3 | Holdout definitions for local-evaluation SDKs | Carried on the existing `GetFeatureFlags` response rather than a new RPC (§4.3.2). |
| 4 | Lifecycle | Registered `start_at`/`start_analysis_at`/`stop_at` **and** manual transitions, on the same status machine as experiments, with an optional `ANALYSIS` stage that freezes the flag list while users stay held out; flags may be added/removed mid-period, both timestamped; archiving a flag that is still a target is rejected until an org admin removes it from the group (§4.6). |
| 5 | Multiple holdout groups | **At most one serving (`RUNNING` or `ANALYSIS`) group per environment at a time.** `WAITING` and finished groups may coexist with it, but the periods of `WAITING` and serving groups may not overlap. A flag belongs to at most one `WAITING` or serving group, and users are therefore mutually exclusive across groups (§4.6, §6.4). |
| 6 | Governance | Two layers. Structural: holdout state does not live on the `Feature`, so no flag API can change it. Permission: every management RPC requires **`Organization_ADMIN` or above**, reusing the existing role model (§4.5). |
| 7 | UI | Dedicated Holdout Groups page under the environment; the flag detail page shows only a read-only "This flag is subject to a holdout" indicator (§4.7). |

## 4. Design Details

### 4.1 Data model

Two new tables, added to `migration/mysql/` and `migration/postgres/` (Atlas). Both are keyed by `(id, environment_id)`.

**`holdout_group`**

| Column | Type | Notes |
|---|---|---|
| `id` | VARCHAR(255) | UUID; also the synthetic `feature_id` in the DWH (§4.4) |
| `environment_id` | VARCHAR(255) | |
| `name`, `description` | VARCHAR(511), TEXT | |
| `assignment_type` | INT | `0=HASH`, `1=SEGMENT` |
| `sampling_seed` | VARCHAR(255) | HASH only; the **environment's holdout hash key**: issued when the environment's first group is created, copied verbatim onto every later group in the environment, then fixed (§4.2.1) |
| `bucket_start` | DECIMAL(6,3) | HASH only; start of the group's bucket range in [0,100), drawn randomly at creation, then fixed (§4.2.1) |
| `percentage` | DECIMAL(6,3) | HASH only; 0 < `percentage` ≤ 50, e.g. `5` means "about 5% of all users". Cap enforced because holdout + comparison arcs total **2 × percentage** on **[0, 1)** and must not overlap (§4.2.1); GrowthBook applies the same cap |
| `segment_id` | VARCHAR(255) | SEGMENT only; fixed at creation, and the segment itself is locked while the group is serving (§4.2.2) |
| `status` | INT | `0=WAITING`, `1=RUNNING`, `2=STOPPED`, `3=FORCE_STOPPED`, `4=ANALYSIS`. A group is **serving** while `RUNNING` or `ANALYSIS`; every guard in this RFC that protects a live group applies to both statuses (§4.6) |
| `start_at`, `stop_at` | BIGINT | epoch seconds |
| `start_analysis_at` | BIGINT | scheduled entry into the analysis period; `0` = no analysis period (§4.6) |
| `analysis_started_at` | BIGINT | when the analysis period actually began; `0` = not started |
| `stopped_at` | BIGINT | manual stop; `0` = not stopped |
| `goal_ids` | JSON | measurement targets; **empty = measurement disabled** for this group (§4.4) |
| `experiment_id` | VARCHAR(255) | internal experiment for the whole period; empty until goals are configured (§4.4) |
| `analysis_experiment_id` | VARCHAR(255) | internal experiment windowed to the analysis period; empty until the group enters `ANALYSIS` (§4.6) |
| `created_by`, `updated_by`, `created_at`, `updated_at`, `archived`, `deleted` | | standard audit columns |

**`holdout_group_feature`** — the intermediate table holding the group's target flags.

| Column | Type | Notes |
|---|---|---|
| `id` | VARCHAR(255) | UUID surrogate key, so re-adding a flag keeps the earlier row's history |
| `holdout_group_id`, `feature_id`, `environment_id` | VARCHAR(255) | |
| `holdout_variation_id` | VARCHAR(255) | variation served to holdout users, chosen per flag |
| `added_at` | BIGINT | epoch seconds |
| `removed_at` | BIGINT | `0` = not removed; soft removal keeps the audit trail. Not stamped when the group stops, so a row counts as a *current* target only while its group is `WAITING`, `RUNNING` or `ANALYSIS` (§4.5) |

Indexes: `(environment_id, status)` on `holdout_group`; `(holdout_group_id, environment_id)` and `(feature_id, environment_id, removed_at)` on `holdout_group_feature`.

### 4.2 Membership

`IsHeldOut(group, user)` lives in the shared evaluation library. Membership is flag-independent, so the evaluator resolves it at most once per group per evaluation, on the first flag the group targets, and reuses the answer for every other flag; an evaluation that touches no targeted flag never computes it.

- **HASH**: the user's bucket is `b = bucket("holdout-{environmentID}-{userID}-{seed}")`, using the existing murmur3-128 bucketer (`evaluation/go/bucketeer.go`), and the user is held out when `b` falls inside the group's range `[bucket_start, bucket_start + percentage) / 100`, wrapping at 1.0. `bucket_start`, `percentage`, and `b` are rounded to the DECIMAL(6,3) grid (thousandths of a percent on [0, 100)) before the half-open comparison, so the exclusive end does not depend on binary float error. A stored percentage of `33.333` therefore ends at `0.33333`, and `33.3334` rounds to the same endpoint.
- **SEGMENT**: the user is an included member of the referenced segment. Evaluation checks the segment users passed in for `segment_id` and state `INCLUDED`. Segment rules are not consulted. An excluded or absent user is not held out.

#### 4.2.1 HASH mode (default)

- **Flag-independent**: the hash input contains no feature ID, so the same users are held out everywhere. This is the one property none of the existing primitives — `Audience`, rollout, targeting — can provide.
- **Environment-scoped seed, per-group range.** The hash input contains no *group ID* either; the group's identity enters only through its bucket range. The seed is issued once per environment (with the first group) and copied onto every group, so all groups in an environment bucket users along the same axis, while the random `bucket_start` gives each successive group a fresh slice — no user is a permanent holdout resident across periods. This choice is made now, before v1, because membership is immutable: a group-scoped seed could never be migrated later, and it would permanently foreclose two things the environment axis gets for free — the matched comparison sample of §4.4 (the range adjacent to the holdout's), and, if the single-serving-group rule of §4.6 is ever lifted, concurrent HASH groups on explicitly disjoint ranges with no assignment bookkeeping (§6.4).
- **The comparison range.** The range `[bucket_start + percentage, bucket_start + 2 × percentage) / 100` — the same-size slice immediately above the holdout's on the circle — is reserved for measurement (§4.4). It has **no effect on serving**: users in it are evaluated by the normal path, byte-identically.
- **Wraparound at 1.0.** Both holdout and comparison ranges are intervals on **[0, 1)** with wrap: if an endpoint exceeds 1.0, it continues from 0.0. Ranges are arcs on the circle, not only contiguous segments on a number line. Example: `bucket_start = 97`, `percentage = 5` (on a 0–100 scale) → holdout covers **[97, 100) ∪ [0, 2)** (5% of users); comparison covers **[2, 7)** (the next 5%). A high `bucket_start` does **not** leave “no room” for the comparison band — wrap places the comparison slice after the holdout arc.
- **Why `percentage ≤ 50`.** Holdout and comparison each occupy an arc of length `percentage / 100` on **[0, 1)**. They are adjacent and must not overlap. That requires **2 × percentage ≤ 100** (i.e. **percentage ≤ 50**). At 5% holdout, roughly **5%** of users are held out, **5%** fall in the comparison band (derived `not_in_holdout` only), and the remaining **~90%** are in neither band — normal serving, no derived program row (§4.4). At the cap **50%**, the two arcs partition the entire circle (50% + 50%); there is no third “rest” band for measurement, which is why 50% is the maximum.
- **`bucket_start` assignment (server-only).** The UI collects **`percentage`** only; **`bucket_start` is not client-supplied**. On **`CreateHoldoutGroup`** with HASH assignment, the holdout domain/service draws **`bucket_start`** uniformly on **[0, 100)** using a CSPRNG (e.g. `crypto/rand`), persists it on `holdout_group` in the same transaction as the insert, and never changes it afterward (§4.6). The dashboard may show the stored value **read-only** on the group detail page for transparency; there is no admin control to pick or edit the start. **`sampling_seed`** for the environment is also issued in this create path (generated on the environment's first group, copied verbatim onto every later group in the environment — §4.2.1); it is not chosen by murmur3 and not sent by the SDK. **Murmur3** (`evaluation/go/bucketeer.go`) only maps each **user** to a fixed **`b ∈ [0, 1)`**; **`IsHeldOut`** compares **`b`** to the stored **`bucket_start`** arc. The subscriber and calculator **read** `bucket_start`; they do **not** randomize it.

#### 4.2.2 SEGMENT mode

The group references an existing segment; its users are the holdout. This is for teams that must hold out a *known* cohort — a named client list, a regulated market, an internal population — where a random slice is not acceptable.

Its costs are real and are surfaced in the UI when the mode is selected:

- **Membership can be changed from outside the group.** Editing the segment, or a bulk user upload, moves users in and out mid-period. `UpdateSegment`, `BulkUploadSegmentUsers` and `DeleteSegment` therefore reject a segment referenced by a **serving** (`RUNNING` or `ANALYSIS`) group. The check is evaluated on demand rather than kept as a column on `segment`, so the status transitions of §4.6 are themselves the lock and unlock, with no state to repair if one is missed. This is the one guard the design otherwise avoids (§6.1), and it is the price of the mode.
- **The lock is a governance requirement, not only a measurement one.** Editing a segment needs `Environment_EDITOR` — the flag editor's role — so without it, gating the group behind `Organization_ADMIN` (§4.5) would be defeated through the referenced segment.
- **Membership resolution costs a lookup.** In HASH mode `IsHeldOut` is pure; in SEGMENT mode it needs the segment's user set at evaluation time (§4.3). Measurement reads the served reason instead of recomputing membership (§4.4), so it needs no lookup.
- **Goal linking needs history after unlock.** `SegmentUserStorage` exposes **current** membership only. While the group is serving the segment is locked, so live membership matches any `goal_ts` in that window for evaluation and for on-time goal linking. After stop the segment is unlocked and may be edited; delayed goals with `goal_ts` during the serve window must **not** use live membership alone (§4.4). v1 persists a **membership snapshot** when the group enters `RUNNING` (refreshed on stop if needed for the grace window) keyed by `holdout_group_id`, used for permanent exclusion and any subscriber path that must resolve SEGMENT holdout at `goal_ts` post-unlock.
- **Results are observational.** A hand-picked cohort is not a randomized control, so the Bayesian readout measures cohort difference plus treatment effect, and SRM is skipped (§4.4). The results UI labels SEGMENT-mode results as observational (§4.7). There is also no hash to derive a matched comparison sample from, so SEGMENT mode falls back to the full non-holdout population as the comparison arm, with the event volume that implies (§4.4).

#### 4.2.3 Conformance

`evaluation/typescript` gets the identical function; `evaluation/testdata` gains a conformance fixture so Go, TypeScript and every SDK port agree bit-for-bit (the existing pattern from `segment_conformance_test.go`).

### 4.3 Evaluation

#### 4.3.1 Order

`assignUser` (`evaluation/go/evaluation.go:266`) gains one step, at the very top:

```
★ holdout ★  →  prerequisites  →  enabled / off-variation  →  targets  →  rules  →  default strategy
```

If the environment has a live group that targets this flag and `IsHeldOut` returns true, evaluation returns the group's `holdout_variation_id` for that flag immediately, with `Reason.Type = HOLDOUT` and the group ID attached. `Reason` gains `HOLDOUT = 7` (unused today) and a `holdout_group_id` field; proto3 enums are open, so older SDKs surface it as an unrecognized value rather than failing. The reason is carried through `EvaluationEvent` and is what makes holdout serving visible in the debug evaluator and in `DebugEvaluateFeatures`.

Why the holdout is evaluated first:

- **One invariant, with no exceptions.** "A held-out user always receives `holdout_variation_id`" is a property of the Holdout Group alone: no combination of flag state produces a different answer, in any SDK port.
- **`flagVariations` stays correct.** `evaluate` records `flagVariations[feature.Id]` for every flag and downstream flags resolve their prerequisites against it. With the holdout first, a downstream prerequisite check sees the variation the user actually received.

Two consequences are accepted, not designed away:

- **The kill switch does not reach holdout users.** Setting `enabled = false` on a target flag turns the feature off for everyone *except* the held-out share, who keep receiving `holdout_variation_id`. In the normal case that is the intent — that variation is the pre-change behavior. When it is not, the escape hatch is `RemoveHoldoutGroupFeatures` or `StopHoldoutGroup`, both effective on the next evaluation and both surfaced at the point of disabling (§4.7). Because both are gated on `Organization_ADMIN` (§4.5), this belongs in incident runbooks: a responder who disables a target flag during an incident must know that releasing the held-out share requires an org admin, and the disable confirmation names the group and its owners for that reason.
- **A holdout user can see a feature whose prerequisite is unsatisfied.** Choosing the flag's pre-change variation as `holdout_variation_id` makes this harmless in practice. `AddHoldoutGroupFeatures` warns when the selected variation differs from `off_variation` on a flag that has prerequisites.

In SEGMENT mode the evaluator additionally needs the group's segment users, which today's plumbing does not fetch — segment IDs are collected per flag. The live group's `segment_id` is added to that set on the server-evaluated `GetEvaluations` path and on `GetSegmentUsers`, so local-evaluation SDKs receive it too. In HASH mode nothing is fetched and the path is unchanged.

#### 4.3.2 Delivery and SDK compatibility

Two delivery paths already exist, and they have very different exposure:

| Path | SDKs | Work needed |
|---|---|---|
| `GetEvaluations` — evaluated server-side by `cmd/api` | Android, iOS, JavaScript, React, React Native, Flutter, and the OpenFeature providers | **Gateway only; no SDK release** (change detection, below). |
| `GetFeatureFlags` — evaluated locally by the SDK | `GO_SERVER`, `NODE_SERVER` | SDK release required. |

On the server-evaluated path, change detection (`UserEvaluationsID`, `EvaluateFeaturesByEvaluatedAt`) looks only at each flag's `UpdatedAt`, so a holdout change would never reach a client. Two changes fix this:

- `UserEvaluationsID` also hashes the live group's ID and `updated_at` (§4.6).
- An `evaluatedAt` older than the group's `updated_at` forces a full re-evaluation (`ForceUpdate`). The comparison uses the same grace window as the feature diff (`evaluatedAt - grace < updated_at`), so a holdout update in the same second as the cursor, or one that arrives just behind it, still forces a full re-evaluation. When serving stops, the empty cache snapshot keeps the `updated_at` of the newest stopped group. A cursor behind that revision forces a full re-evaluation, so clients drop cached `HOLDOUT` assignments instead of merging a partial flag diff. An empty snapshot with `updated_at` 0 means the environment has never stopped a group and does not force an update. Later cacher runs reuse that same stop time. Because `updated_at` is bumped by every edit (§4.6), including a name, description or goal change that does not affect serving, each such edit makes every active client's next request a full re-evaluation once. This is accepted for phase 2; if it proves costly, the cache snapshot can carry a serving-only revision (status transitions and target-flag changes) in place of `updated_at`, with no SDK impact since the value is only compared server-side.

For the local-evaluation path, `GetFeatureFlagsResponse` gains `repeated HoldoutGroup holdout_groups`, always sent in full (the list is tiny and never diffed). It is repeated even though only one group serves at a time (§4.6), so lifting that restriction later costs no wire-format change. It rides on `GetFeatureFlags` rather than a new RPC because a separate endpoint would let the flag and holdout snapshots drift, leaving windows where the holdout is not yet applied or still applied after it ended.

Two consequences for change detection: `feature_flags_id` is generated from the flags alone, so it must fold in the holdout definition as well, or a holdout-only change would never reach an SDK that already holds the current flag set. Once it does, that change produces an empty flag diff, which the current response builder treats as "nothing to send"; it gains a branch that sends the groups with an empty flag diff instead of a full resend.

The honest risk: **an old server SDK ignores the unknown field and silently serves treatments to holdout users**, contaminating the control group without any error. Materializing an approximation into the old payload was rejected (§6.1) — GrowthBook has the same gap, where SDK connections without the prerequisites capability simply never see the holdout rule — so this is handled by observability rather than by a fallback:

- The gateway already receives `source_id` and `sdk_version` on every `GetFeatureFlags` call. When an environment has a live group with targets, requests from a local-evaluation SDK below the minimum supported version increment a dedicated metric labelled by source and version.
- The detail page shows an **SDK compatibility** panel built from that metric, listing non-compliant source/version pairs, and the group cannot be moved to `RUNNING` while unacknowledged incompatible clients are reporting. The panel also tracks the **mixed-path overlap** — held-out users who were also served through an incompatible SDK in the same day (§4.4) — so the residual contamination is a measured number rather than a guess.
- Requests are never rejected. Breaking production traffic to protect a measurement is the wrong trade.

Caching mirrors the feature flags path: a `HoldoutGroupCache` with the same shape as the existing experiments cache, refreshed by the `cacher` batch job and on every holdout write. The gateway reads that cache only. A miss is treated as no live group and is not filled from MySQL on the request path; the cacher stores an empty snapshot when nothing is serving, so a miss is the cold-cache case rather than "no group". SEGMENT mode reuses the existing segment-users cache.

### 4.4 Measurement

**Measurement is opt-in per group.** The entire machinery of this section exists only when `goal_ids` is non-empty: a group with no goals serves the holdout with zero analytics cost — no internal experiment, no derived rows, no dedupe traffic, nothing written to any warehouse. This is a first-class mode, not a degenerate one: teams that analyze in their own warehouse already receive `Reason.HOLDOUT` and the `holdout_group_id` on every evaluation event (§4.3.1) and need nothing else from Bucketeer. Goals may be configured at creation or added later (while `WAITING`, `RUNNING` or `ANALYSIS`); measurement windows from the moment it is enabled. Removing all goals mid-period is rejected — measurement can be switched on, never silently off.

When goals are configured, the group is given a **synthetic feature identity** so that the entire analytics pipeline — goal linking, the DWH schema, the Bayesian calculator, `experiment_result`, and the results UI — is reused. The places that still need a `kind`-aware branch are listed below.

| Field | Value |
|---|---|
| `feature_id` | the holdout group's UUID |
| `feature_version` | `1`, fixed for the life of the group |
| variations | `in_holdout` (baseline) and `not_in_holdout` (the matched comparison sample, below) |
| goals | the group's `goal_ids` |
| window | the group's `start_at` / `stop_at` (`experiment_id`), and `analysis_started_at` / `stop_at` for the analysis-period readout (`analysis_experiment_id`, §4.6) |

When measurement is enabled, the service creates an **internal experiment row** (`experiment.kind = HOLDOUT`, a new field defaulting to `KIND_EXPERIMENT` for every existing row) pointing at that synthetic identity, with `base_variation_id = in_holdout`. Consequences:

- The subscriber links goal events by goal ID against `listExperiments`, and the internal experiment is returned by that same call. **`listExperiments` itself is unchanged**, but the linker in `goal_events_dwh.go` gains holdout-aware branches (§4.4, goal linking and retry).
- The experiment calculator picks it up and produces a standard `experiment_result`.
- `ListExperimentsRequest` gains `repeated Kind kinds`. An empty list means no filter, so the existing internal callers — the subscriber (`goal_events_dwh.go`, `evaluation_events_dwh.go`, `cache_refresher.go`) and the experiment calculator — keep receiving the internal row unchanged. The dashboard's Experiments page sends `kinds = [KIND_EXPERIMENT]`, so internal rows never appear there. The internal row is created through the domain layer, bypassing the public API's feature-existence validation.
- Places that would otherwise treat the internal row as an ordinary experiment branch on `kind = HOLDOUT`:

  | Place | Behavior for `kind = HOLDOUT` |
  |---|---|
  | `UpdateExperiment`, `DeleteExperiment` | Rejected with `FailedPrecondition` for any role, so an `Environment_EDITOR` cannot get around §4.5. The row changes only through the holdout service. |
  | SRM (`getFeatureForSRM`) | Variations and expected weights come from the group, not `GetFeature`. With the matched comparison sample the expected split is **50/50** — the symmetric case the SRM check handles best. SEGMENT mode: `SKIPPED` with a reason. |
  | Results tab | Variation labels come from the group, not `GetFeature`. |

- The internal experiments follow the group: every status transition and every edit of the scheduled timestamps (§4.6) is written to them in the same transaction, so their windows always match the group's. A manual stop sets the experiments to `STOPPED` with `stop_at = stopped_at` — never `FORCE_STOPPED`, which the subscriber does not list — so goal linking continues through the grace window.

The only genuinely new piece is the **exposure denominator**. The calculator counts distinct users per variation from `evaluation_event`, and no SDK emits an event for a flag that does not exist. So the subscriber derives them: while writing evaluation events, for each event whose `feature_id` is a target of the live group, it writes one additional row against the synthetic identity. The variation is taken from **what the user was actually served** for the holdout arm, and from the hash for the comparison arm:

- `in_holdout` — the event's reason is `HOLDOUT` with this `holdout_group_id`.
- `not_in_holdout` — the user's bucket (§4.2.1) falls in the **comparison range** `[bucket_start + percentage, bucket_start + 2 × percentage)`. The subscriber computes the hash *before* touching the dedupe locker, so for the ~`100 − 2p`% of users in neither range the event costs nothing — no Redis call, no row.

Rows are deduplicated per user, group and day through the Redis locker already used for goal events. A `HOLDOUT`-reason event produces only this derived row, and is not written against the target flag's own experiment (§4.8).

The comparison arm is a **matched sample, not the whole non-holdout population**, following GrowthBook's design and for the same reason: the statistical power of the readout is limited by the smaller arm, so tracking all non-holdout users balloons rows, dedupe keys and calculator scan volume — at a 5% holdout, by roughly 10× — without a meaningful reduction in uncertainty. At 100M-MAU scale that difference is the dominant infrastructure cost of the whole feature.

- The resulting population is "users exposed to at least one target flag", sampled symmetrically into two equal-size arms — the correct denominator for a cumulative-impact readout, with an expected 50/50 split that the SRM check verifies directly.
- Both arms trigger on the same condition (an evaluation event on a target flag), so exposure is symmetric by construction.
- Reading the served reason for the holdout arm means **measurement needs no SDK support**. A held-out user reached through an SDK too old to serve the holdout (§4.3.2) produces no derived row on those days — their bucket is in the holdout range, not the comparison range — so an outdated SDK thins the holdout arm's exposure rather than putting treated users into the comparison group. The residual bias (such a user's good-path days still count as `in_holdout` while they intermittently received treatments) is bounded by the mixed-path overlap metric on the compatibility panel (§4.3.2).
- SEGMENT mode has no hash, so it keeps the full non-holdout population as the comparison arm and pays the corresponding event volume; together with its observational nature (§4.2.2) this is surfaced when the mode is chosen.

**Goal linking and retry (`goal_events_dwh`).** Goal events carry no `reason`; the subscriber attributes them by finding experiment(s) whose `goal_ids` contain the event's `goal_id`, then loading the user's latest matching row from `evaluation_event` (`QueryUserEvaluation` on the experiment's `feature_id` and `feature_version`, within the experiment window). If no row exists, today's implementation enqueues a **Redis retry** (`storeRetryOnMiss`) and polls for up to `maxRetryGoalEventPeriod`, on the assumption a client evaluation may still arrive. That assumption is **false** for holdout users on a **per-flag experiment**: §4.4 and §4.8 require that `HOLDOUT`-reason evaluations are **never** written to the target flag's experiment trail, so a qualifying row for that experiment will **never** appear. Without an explicit branch, holdout users' goal events would retry until the period expires — wasted Redis traffic and misleading ops signals — even though the outcome is permanently "not in this experiment."

| Situation | `evaluation_event` row for lookup | Goal for this experiment | Retry? |
|---|---|---|---|
| Held-out user at `goal_ts` (§ below), **flag** experiment (`kind = EXPERIMENT`) on a target flag | **Never** (by design) | Not stored in DWH | **No** — permanent miss; do not call `storeRetryMessage` |
| Non-holdout user, flag experiment, eval not in DWH yet | Not yet | Not stored yet | **Yes** — existing Redis retry (eval may arrive late) |
| Non-holdout user, flag experiment, eval present | Yes | Stored (reason copied from eval) | N/A |
| **Holdout program** experiment (`kind = HOLDOUT`), measurement enabled | Derived row on synthetic `feature_id` (= group UUID); **never** produced for HASH users outside both measurement arcs (§4.2.1) | Stored when linked; skip this experiment for ineligible HASH users, continue other matching experiments | **No** for HASH users outside both arcs (determined from the group's fixed hash config and user id); **Yes** if eligible and derived row not yet visible — same race as today (eval/derive then goal) |
| Holdout serving only (`goal_ids` empty), goal id matches a **flag** experiment only | No row for holdout user | Not stored | **No** — same permanent miss as first row |
| Goal timestamp older than latest eval | Eval exists but too new | Not stored | **No** — same as today (`ErrGoalEventOlderThanEvaluation`) |

**Permanent exclusion (no retry) for flag experiments.** When linking against a **flag** experiment (`kind = EXPERIMENT`), do **not** infer permanence from today's serving snapshot alone. A miss can still be transient if the group had not started yet, the flag was not yet a target, or the user was not held out at the goal timestamp — an evaluation may still arrive (including one generated before enrollment but processed after the goal). Before `storeRetryOnMiss`, treat the miss as **permanent** only when **all** of the following hold at the **goal event timestamp** (`goal_ts`):

1. The holdout group was **serving** (`RUNNING` or `ANALYSIS`) at `goal_ts` (respect the same experiment window boundaries the linker already uses).
2. The flag was a **current target** at `goal_ts`: `added_at ≤ goal_ts` and (`removed_at = 0` or `goal_ts < removed_at`).
3. The user **would be held out** for that group: HASH — `IsHeldOut(group, user)` (fixed for the life of the group); SEGMENT — membership at `goal_ts` from a **retained snapshot** (§4.2.2) while the segment is unlocked, or from live segment-user storage while the group is still serving and the segment remains locked to it. A current lookup alone is **not** sufficient after unlock; if historical membership cannot be established, do **not** treat the miss as permanent — keep the **existing Redis retry** path.

If any check fails, keep the **existing Redis retry** path. When permanent, §4.4 and §4.8 guarantee no qualifying row will ever exist for this experiment — treat like `ErrGoalEventOlderThanEvaluation` on the fresh path (ack, no Redis) and delete any retry key on the retry processor path. HASH mode must not depend on a `HOLDOUT` row in the DWH — membership is known from group config and user id alone.

**Holdout measurement linking.** For `kind = HOLDOUT`, lookup uses the synthetic `feature_id`, not the target flag. **First** apply the eligibility rules in the table above (including HASH arc membership); ineligible users skip this experiment with no retry. For eligible users, a missing row usually means the derived exposure row has not been written yet — **retain retry**. Holdout users reach this path only when the group's `goal_ids` include the event's goal (program-level KPIs). Reusing the **same** `goal_id` on a flag experiment and the holdout program experiment lets one goal event match **both**, with attribution resolved **independently** per experiment via each experiment's `feature_id` / `feature_version` — sharing a goal id is not inherently ambiguous. v1 should either **reject** configuring overlapping goals between a serving holdout internal experiment and a flag experiment (preferred for simpler ops), or link both without suppressing flag experiments; **do not** silently drop flag experiments when a holdout experiment also matches, or valid non-holdout conversions on the flag experiment are lost.

**Evaluation persister (`evaluation_events_dwh`).** When an evaluation event's `reason` is `HOLDOUT`, do **not** append a row attributed to the target flag's experiment (even if one is running). Optionally write only the derived synthetic row described above when measurement is enabled. Non-holdout evaluations on target flags follow the existing experiment-linking path unchanged.

**Gateway event routing (`LinkedEventPublisher`).** The gateway forwards only the events *linked* to an experiment or an auto-ops rule to the dedicated `evaluation-events-dwh` / `goal-events-dwh` / `*-ops` topics, so that the DWH and OPS persisters can be moved off the firehose topics. Today an evaluation event is linked to the DWH when `feature_id:feature_version` matches an experiment returned by `listExperiments` (running, or stopped within `StoppedExperimentEventWindow`), and a goal event when its `goal_id` belongs to one. Once measurement creates the internal `kind = HOLDOUT` experiment, the **goal** rule already covers the program: `linkGoals` matches on `goal_id` alone, `listExperiments` returns the internal row (above), and the row stays listed for the grace window after stop, so program goal events reach the DWH topic with no further change. The **evaluation** rule does not: the internal row's synthetic `feature_id` never appears on an SDK event, and the target flags need not have an experiment of their own, so once the DWH persisters read only the dedicated topics the derived-row writer above would receive no exposures. The publisher therefore gains a holdout branch for evaluation events:

- **Source.** The branch is built from the environment's *measurement* groups — those with `goal_ids` that are serving **or stopped within `StoppedExperimentEventWindow`** — with their target rows including `added_at` / `removed_at`, read through a Redis-backed `listHoldoutGroups` and cached in memory like the experiment and auto-ops targets. It is deliberately **not** driven by the serving snapshot of `HoldoutGroupCache` (§4.3.2): that snapshot is cleared the moment a group stops, while §4.6 keeps the internal experiment linking for two more days so SDK batches flushed late still count. An event is routed when its `timestamp` falls inside `[start_at, stop_at)` and the flag was a target at that time (`added_at ≤ timestamp < removed_at`, or `removed_at = 0`); the subscriber re-checks the same window.
- **Reason-agnostic.** Evaluation events on a target flag are linked whatever their `reason`: the comparison arm is made of *non*-holdout users, so filtering on `Reason.HOLDOUT` would drop half the denominator. In HASH mode the gateway applies the same arc test the subscriber uses — bucket in the holdout range or in the comparison range `[bucket_start + percentage, bucket_start + 2 × percentage)` — before publishing, so the ~`100 − 2p`% of users in neither arc never leave the gateway. SEGMENT mode has no hash and forwards every event on a target flag.
- **Goal pre-filter is kind-aware.** The same HASH-mode arc test may be applied to goal events, but only when every experiment the `goal_id` matched is `kind = HOLDOUT`; a goal shared with a flag experiment is forwarded unconditionally, since that experiment needs every user's goals. The gateway therefore needs the `kind` of the experiments it caches, which `listExperiments` already returns.
- **Eviction.** The holdout targets are evicted with the environment's linked targets on `HOLDOUT_GROUP` domain events, alongside the existing `EXPERIMENT` / `AUTOOPS_RULE` eviction, so a target added mid-period, a goal enabled late or a stop is picked up without waiting for the in-memory TTL.
- A group **without goals** adds nothing to any dedicated topic; the serving-only mode stays free of analytics traffic (§4.4, first paragraph).

The subscriber keeps the final say — it recomputes the arc, the window and the target history, and reads the reason — so the gateway pre-filter is a cost optimization that can only drop events the subscriber would have ignored, never one it would have written.

**Auto-ops event rates are not filtered by reason.** It is tempting to also drop `Reason.HOLDOUT` evaluation events from the OPS topic, since the kill switch cannot reach held-out users (§4.3.1). It must not be done: the event-rate clause compares goal users against evaluation users for a variation, and the goal side (`goal_events_ops`) attributes every goal event with a matching `goal_id` to the clause's **configured** variation **without** looking up what the user was served — goal events carry no `reason`. Removing held-out users from the denominator while they stay in the numerator inflates the rate and makes a false trip more likely. Keeping OPS reason-agnostic makes held-out users behave exactly like every other user under today's counting: a user's evaluation counts toward the variation they were served, and their goal counts toward every clause whose goal matches, whichever variation that clause watches. That exposure-agnostic goal attribution is a pre-existing property of the event-rate rule, not something the holdout introduces — a treatment clause already counts control users' goals in its numerator — and holdout users only change the mix, never seeing the treatment arm. The DWH path is different by design: its goal linker (`QueryUserEvaluation`) reads the user's served reason, which is what makes the `HOLDOUT`-reason filter there consistent.

### 4.5 API and governance

New RPCs on `FeatureService` (`proto/feature/service.proto`), console-facing via the web gateway:

```
CreateHoldoutGroup / GetHoldoutGroup / ListHoldoutGroups / UpdateHoldoutGroup
StartHoldoutGroup  / StartHoldoutGroupAnalysis / StopHoldoutGroup
ArchiveHoldoutGroup / DeleteHoldoutGroup
AddHoldoutGroupFeatures / RemoveHoldoutGroupFeatures
```

`DeleteHoldoutGroup` is accepted only while `WAITING`; a group that has ever served is archived, never deleted, so its measurement history survives.

Requirement 3 — flag editors cannot override the holdout — is satisfied by two independent layers:

- **Structural.** No field on `Feature` carries holdout state, so no flag write path (`UpdateFeature`, scheduled changes, auto-ops, progressive rollout, flag triggers) can touch it, now or in the future.
- **Permission.** Every mutating RPC above requires **`Organization_ADMIN` or above** (`AccountV2_Role_Organization_ADMIN`, `_OWNER`), checked with the existing `role.CheckOrganizationRoleByEnvironmentIDWithLog` — the same helper the organization-scoped audit log APIs use. `Environment_VIEWER` is enough to read a group; no new role is introduced, and `organization.system_admin` is not part of the check. The organization grain is what makes this work: the environment roles cannot distinguish a flag editor from a holdout owner, so an `Environment_EDITOR` deliberately cannot create, start, stop or retarget a group. The cost is that the permission cannot be delegated per environment; a dedicated holdout permission can be added later.

Every mutation writes an audit log event (`HOLDOUT_GROUP_CREATED`, `_UPDATED`, `_FEATURE_ADDED`, `_FEATURE_REMOVED`, `_STARTED`, `_ANALYSIS_STARTED`, `_STOPPED`, `_ARCHIVED`), and adding/removing a flag also writes a feature-scoped entry so the change shows up in that flag's history.

Validation:

- `CreateHoldoutGroup` / `StartHoldoutGroup`, and `UpdateHoldoutGroup` when it changes `start_at`, `start_analysis_at` or `stop_at`: the resulting `[start_at, stop_at)` window must not intersect another `WAITING`, `RUNNING` or `ANALYSIS` group in the environment; `start_at < start_analysis_at < stop_at` when an analysis period is scheduled; `assignment_type` decides whether `percentage` or `segment_id` is required; `percentage` is in `(0, 50]`; `segment_id` must exist in the environment and must not be referenced by a serving group (§4.2.2). On **`CreateHoldoutGroup`**, reject any client-provided **`bucket_start`** or **`sampling_seed`**; for HASH mode the server assigns both as in §4.2.1. Concurrent creates that race on the environment's first group must still yield one canonical seed (transaction or row lock on the environment axis).
- `AddHoldoutGroupFeatures`: the flag exists, is not archived, and is not a current target of any other group (a `holdout_group_feature` row with `removed_at = 0` whose group is `WAITING`, `RUNNING` or `ANALYSIS`; rows left by stopped or archived groups do not count); `holdout_variation_id` is one of the flag's variations; the group is `WAITING` or `RUNNING` — **not `ANALYSIS`**, whose frozen flag list is the point of the stage (§4.6).
- `UpdateFeature` (and every variation write path): **removing a variation that a current target row of a `WAITING`, `RUNNING` or `ANALYSIS` group references as `holdout_variation_id` is rejected**, with an error naming the group — the same pattern as the archive guard. Without this, a variation edit needing only `Environment_EDITOR` would leave the holdout serving a dangling variation ID.

### 4.6 Lifecycle

`WAITING → RUNNING → ANALYSIS → STOPPED`, on the same status machine as experiments, where `ANALYSIS` is optional — `StopHoldoutGroup` is accepted from either serving status. Every transition has two triggers:

- **Scheduled** — `start_at` / `start_analysis_at` / `stop_at`, advanced by a dedicated batch job, `HoldoutGroupStatusUpdater` (`pkg/batch/jobs/holdout`), rather than by the experiment status updater: it has its own schedule (every 30 minutes in production, every minute in development), its own job metrics, and a failure in one does not stall the other. Each run lists the `WAITING`, `RUNNING` and `ANALYSIS` groups of every environment and advances **each group by at most one step**: stop when `stop_at` has passed (including a `WAITING` group whose whole window was missed, so it releases its target flags), otherwise start when `start_at` has passed, otherwise enter `ANALYSIS` when `start_analysis_at` has passed. Stops run before starts within a run, so an expired group never blocks the start of its successor until the next run. A consequence of the one-step rule is that a group started late (for example after batch downtime) whose `start_analysis_at` has also passed reaches `ANALYSIS` on the following run, not the same one; this is intentional.
- **Manual** — `StartHoldoutGroup` starts a `WAITING` group immediately, clamping `start_at` to now; `StartHoldoutGroupAnalysis` moves a `RUNNING` group to `ANALYSIS`, recording `analysis_started_at`; `StopHoldoutGroup` ends a serving group immediately, recording `stopped_at` and status `FORCE_STOPPED`. All triggers, and edits of the three scheduled timestamps, are mirrored onto the internal experiments (§4.4). Every transition and target-flag change bumps `updated_at` (§4.3.2). The actor is recorded in the audit log rather than in a column. The lifecycle is hybrid rather than schedule-only because a holdout deliberately withholds functionality from real users, so it must be stoppable without waiting for `stop_at`.

Both triggers go through the same three RPCs, distinguished by a `mode` field (`MANUAL`, the default, or `SCHEDULED`). `SCHEDULED` is what the batch job sends, and it is strictly stricter than `MANUAL`: it succeeds only once the saved timestamp has actually passed, it never clamps `start_at`, a cleared or postponed `start_analysis_at` leaves the group `RUNNING`, and a scheduled stop records `STOPPED` where a manual stop records `FORCE_STOPPED`. The timestamps are re-read under the environment lock inside the transition transaction, so an administrator who postpones a schedule between the job's list call and its RPC wins. Because `SCHEDULED` can only do less than `MANUAL`, it is not restricted to the batch service token; an org admin calling it by hand gets the same checks, and the audit log records the actor either way.

**Scheduled transition failures.** A scheduled transition can fail for a durable reason — for example a SEGMENT-mode group whose segment was deleted or is still uploading when `start_at` arrives, or an overlap left behind by an edit. The job logs the error, counts it in the job metrics, and retries on every run; the group stays in its current status, and a `WAITING` group that never manages to start is stopped as missed once `stop_at` passes. No failure state is stored on the group in phase 1. Surfacing the last failure reason to the administrator (on the group's detail page, and in the list) is part of the dashboard phase (§5, phase 5); the mechanism — a derived field from the audit log versus a persisted `last_transition_error` — is decided there, once the UI that consumes it exists.

**The analysis period** (adopted from GrowthBook) separates "accumulating changes" from "reading the result". While `ANALYSIS`, serving is identical to `RUNNING` — users stay held out, and every guard tied to a serving group stays in force — but the target-flag list is frozen (`AddHoldoutGroupFeatures` is rejected; `RemoveHoldoutGroupFeatures` stays available as the escape hatch of §4.3.1). On entering `ANALYSIS`, a second internal experiment (`analysis_experiment_id`) is created against the same synthetic identity, windowed `[analysis_started_at, stop_at)`. The calculator — unchanged — then produces two readouts: the whole-period result and the analysis-period result, measured after the general population has experienced all changes together. Without this stage, a flag added in month 5 of 6 dilutes the only readout there is; with it, the dilution is visible in the whole-period number and absent from the analysis-period number.

| Event | Behavior |
|---|---|
| Add a flag mid-period | Allowed while `RUNNING`. `added_at` is recorded and surfaced in the results UI, since a flag added late contributes less to the cumulative effect. **Rejected while `ANALYSIS`** — the frozen list is the point of the stage. |
| Remove a flag mid-period | Allowed at any serving status, with an explicit confirmation: released users immediately get normal evaluation, which contaminates the remaining window. `removed_at` is kept. |
| Archive a target flag | **Rejected** while the flag is a current target of a `WAITING`, `RUNNING` or `ANALYSIS` group. Archiving needs only `Environment_EDITOR`, so auto-removal would let a flag editor release held-out users. The error names the group; an org admin removes the flag first (`RemoveHoldoutGroupFeatures`), which is always possible, so no flag is held hostage. |
| Delete a variation serving as `holdout_variation_id` | **Rejected** while the flag is a current target (§4.5), with an error naming the group. The org admin either picks a different holdout variation (remove and re-add the flag) or removes the flag from the group first. |
| Edit `assignment_type`, `percentage`, `sampling_seed`, `bucket_start` or `segment_id` | **Rejected at every status.** Membership must never change for a given user: raising `percentage` would only add members and lowering it only remove them, but either re-partitions a measurement in progress. A group created with the wrong value is deleted while `WAITING` and recreated. |
| Edit `goal_ids` | Freely editable while `WAITING`. While `RUNNING` or `ANALYSIS`, goals may be **added** (measurement-additive: it re-partitions nothing) but not removed; if the group had no goals, the first addition enables measurement from that moment (§4.4). |
| Edit anything else | `name`, `description` and the three period timestamps are editable while `WAITING`; only `name`, `description`, `start_analysis_at` and `stop_at` while serving. The result is deliberately asymmetric: **target flags and goals are mutable, target users are not**. |
| Start analysis period | `RUNNING → ANALYSIS`, scheduled or manual. Serving unchanged; flag list frozen; `analysis_experiment_id` created (§4.4). |
| Stop / dissolve | Holdout users fall back to normal evaluation on the next evaluation. The internal experiments keep calculating for two more days, matching the existing goal-event grace window. |
| Multiple groups | At most one serving (`RUNNING` or `ANALYSIS`) group per environment, and `WAITING` / serving groups may not have overlapping periods (§4.5). `WAITING` and finished groups coexist freely. Because only one group can ever be serving, **users are mutually exclusive across groups by construction** — no cross-group bookkeeping, and no group's readout confounded by another (§6.4). |

### 4.7 UI

- **Holdout Groups page** (new, environment-scoped, next to Experiments): list with status, assignment mode, size, target-flag count, goal count, period; **detail page with three tabs** — **Overview** (assignment metadata including read-only **hash ring start** for HASH, goals, schedule, target flags table with add/remove, per-flag `holdout_variation_id` and `added_at`, plus the SDK compatibility panel from §4.3.2), **Coverage report** (§ below), **Results** (existing experiment-result component: whole-period and analysis-period readouts side by side when goals are configured, §4.6); manual Start / Start analysis / Stop on the detail page. The whole page is read-only below `Organization_ADMIN`.
- **Coverage report** (v1, on the detail page): flags created or first enabled in the environment during the group's period that are not targets of the group. Every such flag reaches held-out users and dilutes the cumulative readout toward zero, so the gap must be visible while it can still be fixed, not discovered at the end. Entries can be added to the group in place or dismissed with a note. (GrowthBook and Statsig solve the same problem with default enrollment; the coverage report keeps Bucketeer's explicit opt-in while removing its silent failure mode.)
- **Add-flags dialog**: the `holdout_variation_id` picker defaults to the flag's current default/off variation, labeled as recommended. Selecting any other variation requires an explicit confirmation, since choosing a treatment variation silently inverts the control group for that flag.
- **Create wizard**: a radio choice between "by percentage (hash)" and "by user segment", switching the input below it between a percentage field and a segment picker; the SEGMENT option carries the §4.2.2 warning. For HASH mode, short copy under the percentage field explains that holdout and comparison slices are placed at a **random position on the hash ring when the group is created** (§4.2.1) and cannot be changed later — no start-offset control. Measurement is a separate optional step — "measure in Bucketeer" (goal picker) or "analyze externally" (explains the `Reason.HOLDOUT` join, §4.4) — so a serving-only group is a first-class path, not a form with fields left blank.
- **Edit screen**: the assignment mode and everything it implies — percentage, segment, seed, bucket range — are read-only, with a note that changing them means recreating the group; the target-flag list next to them stays editable (§4.6).
- **Flag detail page**: a read-only badge — "This flag is subject to a holdout" — linking to the group by name, identical in both assignment modes. Nothing about the holdout is editable from here; the targeting UI stays editable because it still governs everyone outside the holdout. Without the badge, a flag owner sees a slice of users not getting the new variation with no way to find out why — usually ending in a mistaken rollback.
- **Segment list and detail pages**: a segment referenced by a serving group says so and has its edit and delete actions disabled, mirroring the API-level lock of §4.2.2.
- **SEGMENT-mode results**: labeled observational wherever they are shown (§4.2.2) — the Bayesian intervals measure cohort difference plus treatment effect, and presenting them like a randomized readout invites wrong conclusions.
- **Flag disable confirmation**: turning off a flag that belongs to a live group warns that held-out users will keep receiving `holdout_variation_id` (§4.3.1).
- **Flag archive action**: disabled for a flag that belongs to a live group, with a link to the group (§4.6).
- **Debugger**: `DebugEvaluateFeatures` surfaces `Reason.HOLDOUT` so an operator can confirm why a user got a given variation.

#### 4.7.1 UI prototype (stakeholder mock)

An interactive **dashboard mock** (no API wiring) lives under `ui/dashboard/src/pages/holdouts-prototype/` for review before Phase 5. It shows an amber **UI prototype** banner; data is hard-coded in `mock-data.ts`.

**How to view**

1. `cd ui/dashboard && yarn start`
2. Sign in and open an environment.
3. Sidebar: **Experimentation → Holdout groups** (below Experiments).

**Routes to review** (replace `{env}` with the environment URL code):

| URL | Purpose |
|---|---|
| `/{env}/holdouts` | List |
| `/{env}/holdouts/new` | Create wizard (basics → target flags → goals) |
| `/{env}/holdouts/hld-q1-2026-product` | Detail — Overview, Coverage report, Results |
| `/{env}/holdouts/hld-h2-growth` | Detail — scheduled group, goals set, no target flags yet |

Create step 2: multi-flag selection via **Select feature flags…** (searchable picker aligned with the experiment create flow). The add-many modal uses checkboxes and a per-flag holdout variation selector.

**Stakeholder checklist** (confirm before building Phase 5):

- Sidebar placement and naming
- Multi-flag opt-in (explicit list, not auto-enroll)
- Per-flag holdout variation picker (default/off recommended)
- Create flow: basics → flags → goals (the mock requires goals for review; production still supports serving-only per §4.4)
- Detail tabs: Overview (includes target flags), Coverage report, Results
- HASH create: **percentage** only; **hash ring start** read-only on Overview after create (§4.2.1)

**Suggested screenshots:** sidebar entry, list, create steps 1–2, add-flags modal, detail Overview, Results placeholder.

**Cleanup:** When the real Holdout Groups page ships, delete `pages/holdouts-prototype/` and drop this subsection (or mark it archived).

### 4.8 Interaction with existing features

Auto-ops, flag triggers, prerequisites and the kill switch all change flag configuration, which the holdout short-circuits for its members — the consequences are in §4.3.1. The remaining interactions:

| Feature | Interaction |
|---|---|
| Experiments | An experiment on a target flag measures only non-holdout users, and its sample size shrinks by the holdout share (noted in the experiment UI). The subscriber does not write `HOLDOUT`-reason events to the flag's own experiment (§4.4). Goal events for holdout users then find no evaluation row for that experiment and are **dropped without entering the Redis retry queue** — a permanent miss, not a late eval (§4.4, goal linking and retry). Non-holdout users keep today's retry-until-eval behavior. |
| Auto-ops event rate (kill switch) | A held-out user is counted like any other user under today's rule: their evaluation counts toward `holdout_variation_id`'s evaluation users, and their goal events count toward every clause whose goal matches, regardless of the variation that clause watches — goal attribution is exposure-agnostic for all users already (§4.4, auto-ops event rates). `HOLDOUT`-reason events are **not** filtered out of the OPS topic or the OPS persisters, because goal events cannot be filtered symmetrically. The rule that trips still cannot change what held-out users are served (§4.3.1). |
| Progressive rollout | Continues to operate on non-holdout traffic. A rollout reaching 100% does not release holdout users — that is the point. |
| Segments | In SEGMENT mode, `UpdateSegment`, `BulkUploadSegmentUsers` and `DeleteSegment` are rejected for a segment referenced by a serving group (§4.2.2). |
| Variations | Deleting a variation referenced as a live group's `holdout_variation_id` is rejected (§4.5, §4.6). Other variation edits (names, values) remain free. |
| Archive / auto-archive | Rejected for a current target of a live group; auto-archive skips such flags (§4.6). |
| Code references | Unaffected. |

## 5. Rollout plan

| Phase | Content |
|---|---|
| 1 | Proto (including the `ANALYSIS` status, `StartHoldoutGroupAnalysis` and the `MANUAL` / `SCHEDULED` request modes), DB migrations, domain + storage + API (both assignment modes, environment-scoped seed and bucket ranges, org-admin authorization, single-serving-group enforcement, segment lock, variation guard), audit logs, the `HoldoutGroupStatusUpdater` batch job and the auto-archive exclusion (§4.6). No evaluation behavior yet. |
| 2 | Evaluation: `evaluation/go`, `evaluation/typescript`, conformance fixtures, `Reason.HOLDOUT`, holdout-first ordering, segment-user plumbing, server-side path in `cmd/api` including holdout-aware `UserEvaluationsID` and forced re-evaluation. Every client SDK becomes holdout-capable here. |
| 3 | Delivery to local-evaluation SDKs: `GetFeatureFlagsResponse.holdout_groups`, change detection, cache, compatibility metrics. Go/Node SDK releases follow in their own repos. |
| 4 | Measurement (opt-in per group): internal experiments, `experiment.kind` and its write guard on `ExperimentService`, derived evaluation rows with the matched comparison sample and the `HOLDOUT`-reason filter in the subscriber, holdout-aware goal linking (goal-timestamp permanent discard for flag experiments, HASH arc eligibility and no-retry for holdout program experiments, SEGMENT membership snapshot for post-stop linking, §4.4), dedupe, 50/50 SRM, analysis-period experiment. Holdout branch of the gateway `LinkedEventPublisher` (target-flag evaluation events to the dedicated DWH topic from measurement groups serving or stopped within the grace window, kind-aware HASH-mode arc pre-filter, eviction on `HOLDOUT_GROUP` events, §4.4) — this must land before the DWH persisters are switched to read only the dedicated topics while any group with `goal_ids` is serving. |
| 5 | Dashboard: Holdout Groups page, coverage report, flag indicator, disable confirmation, results tab (both windows), compatibility panel with the mixed-path overlap metric, and the last scheduled-transition failure on the group (§4.6), including the decision on how it is stored. |
| 6 | E2E tests (`test/e2e/gateway`), docs. |

Phases 1–2 are independently shippable and already deliver a working holdout for all client SDKs — exactly the serving-only mode of §4.4, sufficient for teams that analyze externally; phase 4 can be validated against a short group in a dev environment before phase 5 exposes it.

## 6. Alternatives considered

### 6.1 Control-plane expansion at write time

Materialize the holdout into each target flag when the group is saved — a system-managed rule or prerequisite generated by the holdout service. This is how both LaunchDarkly and GrowthBook actually implement holdouts: GrowthBook compiles each holdout into a synthetic parent feature (`$holdout:<id>`, with its own seed) and injects a force rule referencing it as the first rule of every linked feature, and LaunchDarkly attaches the holdout as a prerequisite on each experiment's flag. So cross-flag-consistent HASH membership *is* expressible this way — a single synthetic parent carries the program-scoped hash. Rejected anyway, on the two grounds their implementations demonstrate:

- **Governance becomes a guard instead of a property.** With a holdout reference living on the flag, every flag write path needs a rule to reject edits to it, forever, including future write paths — and in practice the guard is incomplete: GrowthBook stores `feature.holdout` on the feature document, and its feature API lets a flag editor detach a running holdout, precisely what requirement 3 of the issue forbids. §4.2.2 shows the cost of even one such guard done properly.
- **The serving invariant is lost.** A prerequisite-shaped holdout evaluates inside the flag's own machinery, so flag state wins: in GrowthBook, disabling a feature's environment drops it from the payload and held-out users fall to the SDK fallback; in LaunchDarkly, the flag's off state wins over the holdout. The control experience then depends on state any `Environment_EDITOR` can change, silently and with no audit trail on the group (§6.3).
- Its one real advantage — no SDK work — applies only to the two server SDKs, since all client SDKs are evaluated server-side anyway.

### 6.2 A dedicated holdout analysis pipeline

Compute holdout vs. non-holdout aggregates in a purpose-built query instead of the synthetic-feature trick. Rejected for v1: it duplicates the winsorization, Bayesian modelling and result storage that the existing calculator already implements across three data warehouses. The synthetic identity is a small, reversible amount of cleverness (one `kind` column, a list filter, a write guard and two `kind`-aware readers, §4.4) in exchange for reusing that whole stack. If holdout-specific statistics — variance reduction, segment breakdowns — become necessary, a dedicated pipeline can be added later without changing the serving design.

### 6.3 Holdout below prerequisites and the `enabled` check

Evaluate the holdout after prerequisites and the kill switch, so that a disabled flag and an unsatisfied prerequisite both win. It is operationally reassuring and avoids serving a feature whose prerequisite says it is unavailable, but it makes the control experience depend on flag state that any `Environment_EDITOR` can change, silently and with no audit trail on the group — the failure mode §4.3.1 exists to prevent. The explicit escape hatches (remove the flag, stop the group) cover the same operational need while staying auditable.

### 6.4 Multiple concurrent holdout groups

Allow several groups to serve at once, each with its own seed, on disjoint flag sets. With independent per-group seeds this was unsalvageable: the groups overlap at the product of their rates, each readout is taken under the other's presence, and mutual exclusion would need assignment bookkeeping the hash is designed to avoid. The environment-scoped seed of §4.2.1 removes the statistical objection — concurrent HASH groups could simply be allocated disjoint bucket ranges, mutual exclusion by construction — but concurrency is still rejected for v1: it multiplies the UI, validation and measurement surface (HASH-against-SEGMENT and SEGMENT-against-SEGMENT still have no closed form), and serialized groups give clean, independently interpretable readouts at the cost of scheduling, which is a product problem rather than a statistical one. The difference from the previous revision is that this door is now open rather than welded shut.

## 7. Remaining open questions

1. **Minimum SDK versions.** The exact Go/Node SDK versions gating §4.3.2 can only be fixed once those releases are cut.
2. **Nested holdouts.** Layering a holdout inside a holdout, to measure a single program within the universal one, is a natural next step. The single-serving-group rule of §4.6 forecloses it for now; the environment-scoped seed and bucket ranges of §4.2.1 keep the assignment side tractable (a nested group is a sub-range), but the analysis has not been designed. Deferred.
3. **Segment snapshot shape** (§4.2.2). v1 requires a retained membership snapshot for post-stop goal linking; open choice is storage (dedicated table vs blob), refresh policy on stop, and retention through the goal-event grace window. Broader product question: blocking all user edits on a referenced segment may be too coarse for teams reusing large operational segments — snapshot-at-start was chosen partly so unlock-after-stop does not lose attribution history.
4. **`evaluationTotal` semantics** for the synthetic holdout feature. The per-user/per-day dedupe makes it "exposed user-days" rather than raw exposures; with the matched sample of §4.4 the number is at least symmetric across both arms, but it is still shown in the results UI and may need its own label.
5. **Default percentage and period** for the create wizard — a product choice; 5% / 3 months is the common starting point (and 50% the hard cap, §4.1).
6. **Residual mixed-path bias** (§4.4). A held-out user reached through both a holdout-aware path and an outdated server SDK is counted as `in_holdout` on good-path days while having intermittently received treatments; the matched comparison sample already keeps such users out of the comparison arm, so the former both-arms double count is gone, and the remaining bias is bounded by the mixed-path overlap metric on the compatibility panel (§4.3.2). Excluding such users entirely needs a change to the calculator's per-variation user count, which §4.4 otherwise avoids. Deferred until the observed overlap justifies it.
7. **Comparison-sample size.** v1 fixes the comparison range to the holdout's size (50/50 arms, GrowthBook's default). A team wanting a wider comparison arm for secondary analyses could be given a configurable multiplier later; the bucket-range design accommodates it without re-partitioning anyone, as long as the ranges stay disjoint.
8. **Shared `goal_id` between holdout program and flag experiments** (§4.4). Independent per-experiment linking is correct; v1 **prefers rejecting** overlap at configuration time for simpler ops, but must not blanket-filter flag experiments at link time if overlap is allowed.
