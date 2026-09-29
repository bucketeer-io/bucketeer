# Universal Holdout Groups

This RFC describes **Holdout Groups**: a long-lived control group that reserves a fixed share of users across *many* feature flags, so that the cumulative impact of every change shipped during a period (typically 3–6 months) can be measured against a population that never received any of them.

**Issue**: [https://github.com/bucketeer-io/bucketeer/issues/2798](https://github.com/bucketeer-io/bucketeer/issues/2798)

This revision follows the decisions recorded in [this comment](https://github.com/bucketeer-io/bucketeer/issues/2798#issuecomment-5724972912) on the issue.

## 1. Background

A single A/B test answers "did this change help?", not "did everything we shipped this half move the business metric?" — once a winner rolls out, the baseline is gone, and effects that decay or interact go unnoticed. The industry answer is a **holdout**: a fixed slice of users that keeps the *control* variation of every flag in a program for the whole period, while everyone else goes through normal targeting, rollouts and experiments.

## 2. Goals / Non-Goals

**Goals**

- One Holdout Group deterministically decides which users are held out, and the **same** users are held out on every target flag.
- Membership is configurable: a hash-based percentage, or a segment reference.
- The Holdout Group owns its list of target flags. Flags are opted in explicitly; nothing else in the environment is affected.
- Holdout serving is authoritative: nothing on the flag side — targeting, rollout, auto-ops, the kill switch — can pull a user out of the holdout.
- Users outside the holdout are evaluated by the normal path with byte-identical results.
- Aggregate measurement: holdout vs. non-holdout, over a set of goals, using the existing Bayesian experiment calculator.
- Lifecycle: scheduled start/stop *and* manual activation/deactivation, adding and removing target flags mid-period, dissolution.
- Full audit trail; MySQL and PostgreSQL support.

**Non-Goals**

- Nested or hierarchical holdouts (a holdout inside a holdout).
- More than one Holdout Group serving at the same time in an environment (see §4.6 and §6.4).
- Cross-environment or cross-project holdouts. A Holdout Group is environment-scoped, like flags and experiments.
- Automatic enrollment of newly created flags. Opt-in stays explicit (a "suggest flags not in any holdout" UI hint is future work).

## 3. Design Overview

A Holdout Group is a **new first-class entity in the evaluation engine**, not a materialization into flag configuration (§6.1). Membership is a pure function of the group's own configuration, serving happens before any flag state is consulted, and measurement reuses the existing experiment pipeline by giving the group a synthetic feature identity.

The decisions below answer the open questions raised on the issue.

| # | Question | Decision |
|---|---|---|
| 1 | User assignment model | Configurable via `assignment_type`: hash-based percentage with a group-scoped seed (default), or a segment reference. `assignment_type`, `percentage`, `sampling_seed` and `segment_id` are all fixed at creation, so membership never changes for a given user (§4.2, §4.6). |
| 2 | Data model / evaluation strategy | New entities `holdout_group` and `holdout_group_feature`; target flags are managed through the intermediate table (§4.1, §4.3.1). |
| 3 | Holdout definitions for local-evaluation SDKs | Carried on the existing `GetFeatureFlags` response rather than a new RPC (§4.3.2). |
| 4 | Lifecycle | Registered `start_at`/`stop_at` **and** manual start/stop, on the same status machine as experiments; flags may be added/removed mid-period, both timestamped; archiving a target flag removes it automatically (§4.6). |
| 5 | Multiple holdout groups | **At most one non-stopped group serves per environment at a time.** Not-yet-started and already-finished groups may coexist with it; their periods may not overlap. A flag belongs to at most one non-stopped group, and users are therefore mutually exclusive across groups (§4.6, §6.4). |
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
| `sampling_seed` | VARCHAR(255) | HASH only; issued at creation, then fixed |
| `percentage` | DECIMAL(6,3) | HASH only; 0–100, e.g. `5` means "about 5% of all users" |
| `segment_id` | VARCHAR(255) | SEGMENT only; fixed at creation, and the segment itself is locked while `RUNNING` (§4.2.2) |
| `status` | INT | `0=WAITING`, `1=RUNNING`, `2=STOPPED`, `3=FORCE_STOPPED` |
| `start_at`, `stop_at` | BIGINT | epoch seconds |
| `stopped_at` | BIGINT | manual stop; `0` = not stopped |
| `goal_ids` | JSON | measurement targets (§4.4) |
| `experiment_id` | VARCHAR(255) | internal experiment used for analysis (§4.4) |
| `created_by`, `updated_by`, `created_at`, `updated_at`, `archived`, `deleted` | | standard audit columns |

**`holdout_group_feature`** — the intermediate table holding the group's target flags.

| Column | Type | Notes |
|---|---|---|
| `id` | VARCHAR(255) | UUID surrogate key, so re-adding a flag keeps the earlier row's history |
| `holdout_group_id`, `feature_id`, `environment_id` | VARCHAR(255) | |
| `holdout_variation_id` | VARCHAR(255) | variation served to holdout users, chosen per flag |
| `added_at` | BIGINT | epoch seconds |
| `removed_at` | BIGINT | `0` = still a target; soft removal keeps the audit trail |

Indexes: `(environment_id, status)` on `holdout_group`; `(holdout_group_id, environment_id)` and `(feature_id, environment_id, removed_at)` on `holdout_group_feature`.

### 4.2 Membership

`IsHeldOut(group, user)` lives in the shared evaluation library:

- **HASH**: `bucket("holdout-{holdoutGroupID}-{userID}-{seed}") < percentage / 100`, using the existing murmur3-128 bucketer (`evaluation/go/bucketeer.go`).
- **SEGMENT**: the user is a member of the referenced segment.

#### 4.2.1 HASH mode (default)

- **Flag-independent**: the hash input contains no feature ID, so the same users are held out everywhere. This is the one property none of the existing primitives — `Audience`, rollout, targeting — can provide.

#### 4.2.2 SEGMENT mode

The group references an existing segment; its users are the holdout. This is for teams that must hold out a *known* cohort — a named client list, a regulated market, an internal population — where a random slice is not acceptable.

Its costs are real and are surfaced in the UI when the mode is selected:

- **Membership can be changed from outside the group.** Editing the segment, or a bulk user upload, moves users in and out mid-period. `UpdateSegment`, `BulkUploadSegmentUsers` and `DeleteSegment` therefore reject a segment referenced by a **`RUNNING`** group. The check is evaluated on demand rather than kept as a column on `segment`, so the status transitions of §4.6 are themselves the lock and unlock, with no state to repair if one is missed. This is the one guard the design otherwise avoids (§6.1), and it is the price of the mode.
- **The lock is a governance requirement, not only a measurement one.** Editing a segment needs `Environment_EDITOR` — the flag editor's role — so without it, gating the group behind `Organization_ADMIN` (§4.5) would be defeated through the referenced segment.
- **Membership resolution costs a lookup.** In HASH mode `IsHeldOut` is pure; in SEGMENT mode it needs the segment's user set, both at evaluation time (§4.3) and at measurement time (§4.4).

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

Three consequences are accepted, not designed away:

- **The kill switch does not reach holdout users.** Setting `enabled = false` on a target flag turns the feature off for everyone *except* the held-out share, who keep receiving `holdout_variation_id`. In the normal case that is the intent — that variation is the pre-change behavior. When it is not, the escape hatch is `RemoveHoldoutGroupFeatures` or `StopHoldoutGroup`, both effective on the next evaluation and both surfaced at the point of disabling (§4.7).
- **A holdout user can see a feature whose prerequisite is unsatisfied.** Choosing the flag's pre-change variation as `holdout_variation_id` makes this harmless in practice. `AddHoldoutGroupFeatures` warns when the selected variation differs from `off_variation` on a flag that has prerequisites.

In SEGMENT mode the evaluator additionally needs the group's segment users, which today's plumbing does not fetch — segment IDs are collected per flag. The live group's `segment_id` is added to that set on the server-evaluated `GetEvaluations` path and on `GetSegmentUsers`, so local-evaluation SDKs receive it too. In HASH mode nothing is fetched and the path is unchanged.

#### 4.3.2 Delivery and SDK compatibility

Two delivery paths already exist, and they have very different exposure:

| Path | SDKs | Work needed |
|---|---|---|
| `GetEvaluations` — evaluated server-side by `cmd/api` | Android, iOS, JavaScript, React, React Native, Flutter, and the OpenFeature providers | **None.** The gateway runs the shared evaluator, so holdouts apply the moment the feature ships. |
| `GetFeatureFlags` — evaluated locally by the SDK | `GO_SERVER`, `NODE_SERVER` | SDK release required. |

For the local-evaluation path, `GetFeatureFlagsResponse` gains `repeated HoldoutGroup holdout_groups`, always sent in full (the list is tiny and never diffed). It is repeated even though only one group serves at a time (§4.6), so lifting that restriction later costs no wire-format change. It rides on `GetFeatureFlags` rather than a new RPC because a separate endpoint would let the flag and holdout snapshots drift, leaving windows where the holdout is not yet applied or still applied after it ended.

Two consequences for change detection: `feature_flags_id` is generated from the flags alone, so it must fold in the holdout definition as well, or a holdout-only change would never reach an SDK that already holds the current flag set. Once it does, that change produces an empty flag diff, which the current response builder treats as "nothing to send"; it gains a branch that sends the groups with an empty flag diff instead of a full resend.

The honest risk: **an old server SDK ignores the unknown field and silently serves treatments to holdout users**, contaminating the control group without any error. There is no way to express the holdout in the old payload (§6.1), so this is handled by observability rather than by a fallback:

- The gateway already receives `source_id` and `sdk_version` on every `GetFeatureFlags` call. When an environment has a live group with targets, requests from a local-evaluation SDK below the minimum supported version increment a dedicated metric labelled by source and version.
- The detail page shows an **SDK compatibility** panel built from that metric, listing non-compliant source/version pairs, and the group cannot be moved to `RUNNING` while unacknowledged incompatible clients are reporting.
- Requests are never rejected. Breaking production traffic to protect a measurement is the wrong trade.

Caching mirrors the feature flags path: a `HoldoutGroupCache` with the same shape as the existing experiments cache, refreshed by the `cacher` batch job and evicted on write. SEGMENT mode reuses the existing segment-users cache.

### 4.4 Measurement

The group is given a **synthetic feature identity** so that the entire analytics pipeline — goal linking, the DWH schema, the Bayesian calculator, `experiment_result`, and the results UI — is reused without modification:

| Field | Value |
|---|---|
| `feature_id` | the holdout group's UUID |
| `feature_version` | `1`, fixed for the life of the group |
| variations | `in_holdout` (baseline) and `not_in_holdout` |
| goals | the group's `goal_ids` |
| window | the group's `start_at` / `stop_at` |

On creation, the service also creates an **internal experiment row** (`experiment.kind = HOLDOUT`, a new field defaulting to `KIND_EXPERIMENT` for every existing row) pointing at that synthetic identity, with `base_variation_id = in_holdout`. Consequences:

- The subscriber links goal events by goal ID against `listExperiments`, and the internal experiment is returned by that same call, so **goal linking needs no change at all**.
- The experiment calculator picks it up and produces a standard `experiment_result`, so the Results tab is the existing experiment-result component pointed at a different ID.
- `ListExperiments` filters `kind = KIND_EXPERIMENT` by default so internal rows never appear on the Experiments page; the internal row is created through the domain layer, bypassing the public API's feature-existence validation.

The only genuinely new piece is the **exposure denominator**. The calculator counts distinct users per variation from `evaluation_event`, and no SDK emits an event for a flag that does not exist. So the subscriber derives them: while writing evaluation events, for each event whose `feature_id` is a target of the live group, it writes one additional row against the synthetic identity, with the variation computed by the shared `IsHeldOut`. Rows are deduplicated per user, group and day through the Redis locker already used for goal events, which also keeps the SEGMENT-mode lookup to one per user per day.

- The resulting population is "users exposed to at least one target flag", which is the correct denominator for a cumulative-impact readout.
- Deriving membership server-side means **measurement needs no SDK support at all**, including from SDKs too old to serve the holdout — those users are correctly counted as `not_in_holdout`, so an outdated SDK degrades the effect size rather than corrupting the assignment.

### 4.5 API and governance

New RPCs on `FeatureService` (`proto/feature/service.proto`), console-facing via the web gateway:

```
CreateHoldoutGroup / GetHoldoutGroup / ListHoldoutGroups / UpdateHoldoutGroup
StartHoldoutGroup  / StopHoldoutGroup / ArchiveHoldoutGroup / DeleteHoldoutGroup
AddHoldoutGroupFeatures / RemoveHoldoutGroupFeatures
```

`DeleteHoldoutGroup` is accepted only while `WAITING`; a group that has ever served is archived, never deleted, so its measurement history survives.

Requirement 3 — flag editors cannot override the holdout — is satisfied by two independent layers:

- **Structural.** No field on `Feature` carries holdout state, so no flag write path (`UpdateFeature`, scheduled changes, auto-ops, progressive rollout, flag triggers) can touch it, now or in the future.
- **Permission.** Every mutating RPC above requires **`Organization_ADMIN` or above** (`AccountV2_Role_Organization_ADMIN`, `_OWNER`), checked with the existing `role.CheckOrganizationRoleByEnvironmentIDWithLog` — the same helper the organization-scoped audit log APIs use. `Environment_VIEWER` is enough to read a group; no new role is introduced, and `organization.system_admin` is not part of the check. The organization grain is what makes this work: the environment roles cannot distinguish a flag editor from a holdout owner, so an `Environment_EDITOR` deliberately cannot create, start, stop or retarget a group. The cost is that the permission cannot be delegated per environment; a dedicated holdout permission can be added later.

Every mutation writes an audit log event (`HOLDOUT_GROUP_CREATED`, `_UPDATED`, `_FEATURE_ADDED`, `_FEATURE_REMOVED`, `_STARTED`, `_STOPPED`, `_ARCHIVED`), and adding/removing a flag also writes a feature-scoped entry so the change shows up in that flag's history.

Validation:

- `CreateHoldoutGroup` / `StartHoldoutGroup`: the `[start_at, stop_at)` window must not intersect another non-stopped, non-archived group in the environment; `assignment_type` decides whether `percentage` or `segment_id` is required; `segment_id` must exist in the environment and must not be referenced by a `RUNNING` group (§4.2.2).
- `AddHoldoutGroupFeatures`: the flag exists, is not archived, and is not a current target of any other group (a `holdout_group_feature` row with `removed_at = 0`); `holdout_variation_id` is one of the flag's variations; the group is `WAITING` or `RUNNING`.

### 4.6 Lifecycle

`WAITING → RUNNING → STOPPED`, on the same status machine as experiments. Both transitions have two triggers:

- **Scheduled** — `start_at` / `stop_at`, advanced by the same batch job that advances experiment statuses (`pkg/batch/jobs/experiment`).
- **Manual** — `StartHoldoutGroup` starts a `WAITING` group immediately, clamping `start_at` to now; `StopHoldoutGroup` ends a `RUNNING` group immediately, recording `stopped_at` and status `FORCE_STOPPED`. The actor is recorded in the audit log rather than in a column. The lifecycle is hybrid rather than schedule-only because a holdout deliberately withholds functionality from real users, so it must be stoppable without waiting for `stop_at`.

| Event | Behavior |
|---|---|
| Add a flag mid-period | Allowed. `added_at` is recorded and surfaced in the results UI, since a flag added late contributes less to the cumulative effect. |
| Remove a flag mid-period | Allowed, with an explicit confirmation: released users immediately get normal evaluation, which contaminates the remaining window. `removed_at` is kept. |
| Archive a target flag | **Archiving is not rejected.** `removed_at` is stamped automatically and the flag drops out on the next evaluation, exactly as a manual removal would. Blocking the archive would let a measurement hold a flag hostage for months. |
| Edit `assignment_type`, `percentage`, `sampling_seed` or `segment_id` | **Rejected at every status.** Membership must never change for a given user: raising `percentage` would only add members and lowering it only remove them, but either re-partitions a measurement in progress. A group created with the wrong value is deleted while `WAITING` and recreated. |
| Edit anything else | `name`, `description`, `goal_ids` and the period are editable while `WAITING`; only `name`, `description` and `stop_at` while `RUNNING`. The result is deliberately asymmetric: **target flags are mutable, target users are not**. |
| Stop / dissolve | Holdout users fall back to normal evaluation on the next evaluation. The internal experiment keeps calculating for two more days, matching the existing goal-event grace window. |
| Multiple groups | At most one non-stopped group per environment, and non-stopped groups may not have overlapping periods (§4.5). Groups that have not started, or have already finished, coexist freely. Because only one group can ever be serving, **users are mutually exclusive across groups by construction** — no cross-group bookkeeping, and no group's readout confounded by another (§6.4). |

### 4.7 UI

- **Holdout Groups page** (new, environment-scoped, next to Experiments): list with status, assignment mode, size, target-flag count, period; detail page with Overview, Target flags (add/remove, per-flag `holdout_variation_id`, `added_at`), Results (existing experiment-result component), and the SDK compatibility panel from §4.3.2, plus manual Start/Stop actions. The whole page is read-only below `Organization_ADMIN`.
- **Create wizard**: a radio choice between "by percentage (hash)" and "by user segment", switching the input below it between a percentage field and a segment picker; the SEGMENT option carries the §4.2.2 warning.
- **Edit screen**: the assignment mode and everything it implies — percentage, segment, seed — are read-only, with a note that changing them means recreating the group; the target-flag list next to them stays editable (§4.6).
- **Flag detail page**: a read-only badge — "This flag is subject to a holdout" — linking to the group by name, identical in both assignment modes. Nothing about the holdout is editable from here; the targeting UI stays editable because it still governs everyone outside the holdout. Without the badge, a flag owner sees a slice of users not getting the new variation with no way to find out why — usually ending in a mistaken rollback.
- **Segment list and detail pages**: a segment referenced by a `RUNNING` group says so and has its edit and delete actions disabled, mirroring the API-level lock of §4.2.2.
- **Flag disable confirmation**: turning off a flag that belongs to a live group warns that held-out users will keep receiving `holdout_variation_id` (§4.3.1).
- **Debugger**: `DebugEvaluateFeatures` surfaces `Reason.HOLDOUT` so an operator can confirm why a user got a given variation.

### 4.8 Interaction with existing features

Auto-ops, flag triggers, prerequisites and the kill switch all change flag configuration, which the holdout short-circuits for its members — the consequences are in §4.3.1. The remaining interactions:

| Feature | Interaction |
|---|---|
| Experiments | An experiment on a target flag measures only non-holdout users. Its result stays valid (holdout users never enter any variation's exposure denominator for that flag) but its sample size shrinks by the holdout share. The experiment UI notes this. |
| Progressive rollout | Continues to operate on non-holdout traffic. A rollout reaching 100% does not release holdout users — that is the point. |
| Segments | In SEGMENT mode, `UpdateSegment`, `BulkUploadSegmentUsers` and `DeleteSegment` are rejected for a segment referenced by a `RUNNING` group (§4.2.2). |
| Auto-archive | Unaffected. An archived flag is stamped `removed_at` and leaves the group (§4.6). |
| Code references | Unaffected. |

## 5. Rollout plan

| Phase | Content |
|---|---|
| 1 | Proto, DB migrations, domain + storage + API (both assignment modes, org-admin authorization, single-active-group enforcement, segment lock), audit logs. No evaluation behavior yet. |
| 2 | Evaluation: `evaluation/go`, `evaluation/typescript`, conformance fixtures, `Reason.HOLDOUT`, holdout-first ordering, segment-user plumbing, server-side path in `cmd/api`. Every client SDK becomes holdout-capable here. |
| 3 | Delivery to local-evaluation SDKs: `GetFeatureFlagsResponse.holdout_groups`, change detection, cache, compatibility metrics. Go/Node SDK releases follow in their own repos. |
| 4 | Measurement: internal experiment, `experiment.kind`, derived evaluation rows in the subscriber, dedupe. |
| 5 | Dashboard: Holdout Groups page, flag indicator, disable confirmation, results tab, compatibility panel. |
| 6 | E2E tests (`test/e2e/gateway`), docs. |

Phases 1–2 are independently shippable and already deliver a working holdout for all client SDKs; phase 4 can be validated against a short group in a dev environment before phase 5 exposes it.

## 6. Alternatives considered

### 6.1 Control-plane expansion at write time

Materialize the holdout into each target flag when the group is saved — a system-managed rule or prerequisite generated by the holdout service. Rejected:

- **It cannot express the core requirement.** A generated rule can only select users by attribute, segment membership or a rollout strategy, and the rollout hash is per-flag. The only cross-flag-consistent construct available is a segment, which is exactly the SEGMENT mode of §4.2.2 — with all of its bias, and none of the hash mode. There is no way to express HASH-mode membership in flag configuration.
- **Governance becomes a guard instead of a property.** With state living on the flag, every flag write path needs a rule to reject edits to holdout-owned rules, forever, including future write paths. §4.2.2 shows the cost of even one such guard.
- Its one real advantage — no SDK work — applies only to the two server SDKs, since all client SDKs are evaluated server-side anyway.

### 6.2 A dedicated holdout analysis pipeline

Compute holdout vs. non-holdout aggregates in a purpose-built query instead of the synthetic-feature trick. Rejected for v1: it duplicates the winsorization, Bayesian modelling and result storage that the existing calculator already implements across three data warehouses. The synthetic identity is a small, reversible amount of cleverness (one `kind` column and a list filter) in exchange for reusing that whole stack. If holdout-specific statistics — variance reduction, segment breakdowns — become necessary, a dedicated pipeline can be added later without changing the serving design.

### 6.3 Holdout below prerequisites and the `enabled` check

Evaluate the holdout after prerequisites and the kill switch, so that a disabled flag and an unsatisfied prerequisite both win. It is operationally reassuring and avoids serving a feature whose prerequisite says it is unavailable, but it makes the control experience depend on flag state that any `Environment_EDITOR` can change, silently and with no audit trail on the group — the failure mode §4.3.1 exists to prevent. The explicit escape hatches (remove the flag, stop the group) cover the same operational need while staying auditable.

### 6.4 Multiple concurrent holdout groups

Allow several groups to serve at once, each with its own seed, on disjoint flag sets. Rejected: independent hashes make the groups overlap at the product of their rates, so each group's readout is an estimate taken under the other's presence, and the overlap population grows with every group. Making users mutually exclusive across concurrent groups would require assignment bookkeeping the hash is specifically designed to avoid. Serializing the groups (§4.6) gives clean, independently interpretable readouts at the cost of scheduling, and scheduling is a product problem rather than a statistical one.

## 7. Remaining open questions

1. **Minimum SDK versions.** The exact Go/Node SDK versions gating §4.3.2 can only be fixed once those releases are cut.
2. **Nested holdouts.** Layering a holdout inside a holdout, to measure a single program within the universal one, is a natural next step. The single-active-group rule of §4.6 forecloses it for now; lifting that rule for the nested case specifically is the likely path, but neither the assignment interaction nor the analysis has been designed. Deferred.
3. **Concurrent groups, if the restriction is lifted** (§6.4). Mutual exclusion would have to be built, and its cost depends on the modes involved: two HASH groups can share a seed and split the bucket range, which is cheap, while HASH-against-SEGMENT and SEGMENT-against-SEGMENT have no closed form and need the overlap computed per pair up front.
4. **Segment locking granularity** (§4.2.2). Blocking all user edits on a segment referenced by a live group is the safe default, but it may be too coarse if a team wants to reuse a large operational segment. An alternative is to snapshot the segment's user set into the group at start.
5. **`evaluationTotal` semantics** for the synthetic holdout feature. The per-user/per-day dedupe makes it "exposed user-days" rather than raw exposures; the Bayesian models consume user counts, but the number is still shown in the results UI and may need its own label.
6. **Default percentage and period** for the create wizard — a product choice; 5% / 3 months is the common starting point.
