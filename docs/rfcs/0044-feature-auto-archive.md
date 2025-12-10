# Summary

We will implement an environment-level auto-archive system that enables administrators to automatically archive unused feature flags across an entire environment. This feature addresses the challenge of managing stale feature flags at scale by providing centralized, environment-wide automated cleanup with safety guarantees through dependency checks and code reference validation.

The auto-archive system will allow administrators to configure at the environment level:

- **Enable/Disable auto-archiving**: Toggle auto-archiving for all flags in the environment
- **Unused days threshold**: Minimum number of days flags must be unused before archiving
- **Code reference requirements**: Whether to require zero code references before archiving

## Background

As feature flag systems mature, organizations accumulate hundreds or thousands of flags, many of which become obsolete over time. Managing these flags individually is time-consuming and error-prone, often leading to:

- **Technical debt**: Unused flags cluttering the codebase and dashboard
- **Maintenance burden**: Developers spending time managing obsolete flags manually
- **Inconsistent cleanup**: Different retention policies applied to different flags
- **Performance impact**: Evaluation overhead from unnecessary flags

The existing Bucketeer system provides:

- Manual archive functionality via `ArchiveFeature` API
- `Feature.archived` field (boolean) for marking archived flags
- Dependency checking via `HasFeaturesDependsOnTargets`
- Last used tracking via `feature_last_used_info` table
- Code reference tracking for flags

This RFC proposes adding environment-level auto-archive configuration, allowing administrators to set a single policy that applies to all flags in an environment, significantly reducing the management overhead.

## Goals


- Enable environment-wide automated archiving of unused feature flags
- Provide centralized configuration at the environment level
- Maintain safety through dependency validation and code reference checks
- Provide audit trail for all auto-archive operations
- Reuse existing archive infrastructure for consistency
- Simplify management compared to per-flag configuration


## Implementation

### Database Schema Changes

Add three columns to the existing `environment_v2` table to store auto-archive configuration:

```sql
ALTER TABLE environment_v2
  ADD COLUMN auto_archive_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN auto_archive_unused_days INT NOT NULL DEFAULT 90,
  ADD COLUMN auto_archive_check_code_refs BOOLEAN NOT NULL DEFAULT TRUE;

CREATE INDEX idx_environment_auto_archive_enabled
  ON environment_v2 (auto_archive_enabled);
```

#### Design Decision: Configuration Storage Approach

We evaluated multiple approaches for storing auto-archive configuration.

**Option A: Direct columns in environment_v2 (Chosen)**
- Consistent with existing pattern (`require_comment` field)
- Simple implementation with single-table queries
- High performance (no JOIN overhead)

**Option B: Separate configuration table**
- Better separation of concerns
- Prevents environment_v2 table bloat
- Trade-offs: JOIN overhead in queries (critical for batch jobs), implementation complexity (2-table transactions, NULL handling, UPSERT logic), 1:1 relationship redundancy, inconsistency with existing patterns

**Decision Rationale**: We chose Option A because three fields represent a reasonable scope for this feature, and maintaining consistency with the existing `require_comment` pattern ensures simplicity and performance. This approach minimizes implementation cost and operational complexity while delivering the required functionality.

**Future Consideration**: If environment-level configuration fields exceed 5-7 feature-specific settings, we should revisit this architecture and consider a more scalable configuration management approach, such as a unified configuration table or feature-specific config tables.

### Protobuf Changes

#### Environment Proto Extension

**File**: `proto/environment/environment.proto`

```protobuf
message EnvironmentV2 {
  string id = 1;
  string name = 2;
  string url_code = 3;
  string description = 4;
  string project_id = 5;
  bool archived = 6;
  int64 created_at = 7;
  int64 updated_at = 8;
  string organization_id = 9;
  bool require_comment = 10;
  int32 feature_flag_count = 11;

  // NEW: Auto-archive configuration
  bool auto_archive_enabled = 12;
  int32 auto_archive_unused_days = 13;
  bool auto_archive_check_code_refs = 14;
}
```

#### Service Proto Extension

**File**: `proto/environment/service.proto`

```protobuf
message UpdateEnvironmentV2Request {
  string id = 1;
  // ... existing fields ...

  // NEW: Auto-archive configuration
  google.protobuf.BoolValue auto_archive_enabled = X;
  google.protobuf.Int32Value auto_archive_unused_days = Y;
  google.protobuf.BoolValue auto_archive_check_code_refs = Z;
}
```

**Design Decision**: We extend `UpdateEnvironmentV2Request` to allow updating auto-archive settings alongside other environment properties, consistent with how `require_comment` is managed.

#### Batch Job Proto Extension

**File**: `proto/batch/service.proto`

```protobuf
enum BatchJob {
  // ... existing jobs ...
  FeatureAutoArchiver = 19;  // NEW
}
```

### System Architecture

The auto-archive feature integrates with the existing Bucketeer architecture through environment-level configuration and batch processing:

#### Environment Model Extensions

The environment domain model will be extended to support auto-archive configuration as environment-level properties. Each environment will define its own auto-archive policy that applies uniformly to all flags within that environment.

The existing `EnvironmentV2.Update()` or dedicated change methods (following the pattern of `ChangeRequireComment()`) will be extended to accept auto-archive settings, ensuring that all environment updates follow the same transaction and validation patterns.

**Configuration scope:**

- **Environment level**: Single policy per environment
- **Applied to**: All non-archived flags in the environment
- **Evaluation**: Based on each flag's usage history, code references, and dependencies

#### API Extensions

The existing `UpdateEnvironmentV2` API will be extended to accept auto-archive settings as part of the environment update request. This approach maintains consistency with how other environment properties are managed.

Validation will ensure that:

- Unused days threshold is a positive integer when auto-archive is enabled
- Settings changes are properly authorized (requires environment admin permissions)
- Configuration changes trigger appropriate domain events for audit logging

#### Storage Layer

The storage layer will be extended to support environment auto-archive settings:

- Update all environment SELECT queries to include the three new columns
- Update CREATE and UPDATE queries to handle auto-archive settings
- Provide method to retrieve environments with auto-archiving enabled


#### Automated Batch Processing

A new batch job will run daily (configurable via Kubernetes CronJob) to identify and archive eligible features across all environments. The job will:

1. **Query Environments**: Retrieve only environments with `auto_archive_enabled = true`
2. **For Each Environment**:
   a. Retrieve all non-archived flags in the environment
   b. Bulk fetch code reference counts for all flags
   c. For each flag, validate:
      - Must not be already archived
      - Must have exceeded the environment's unused days threshold
      - Must meet the environment's code reference requirements
      - Must not have any dependent features (safety check)
   d. Archive qualifying flags using existing `ArchiveFeatureCommand`
3. **Audit**: Generate events and logs for all archiving operations

**Batch Job Flow**:

```text
1. Query: SELECT * FROM environment_v2 WHERE auto_archive_enabled = true
2. For each environment:
   a. Query: SELECT * FROM feature
      WHERE environment_id = ?
      AND archived = false
   b. Bulk query code references
   c. Evaluate each flag against environment policy
   d. Archive eligible flags
3. Log results and metrics
```

The batch job integrates with the existing job infrastructure, metrics collection, and monitoring systems.

#### User Interface

The web dashboard will provide auto-archive configuration in the environment settings page. Administrators will be able to:

- Enable/disable auto-archiving for the entire environment
- Set the unused days threshold (minimum days without usage before archiving)
- Configure whether to require zero code references before archiving
- View the current configuration and understand its impact

**UI Location**: Environment Settings page (`/environment/{id}/settings`)

**UI Section**: "Auto-Archive Configuration" section with:

- Toggle switch for enabling auto-archiving
- Number input for unused days threshold (enabled when toggle is on)
- Checkbox for code reference requirement
- Help text explaining the policy applies to all flags in the environment
- Warning about the impact of enabling auto-archiving


## Configuration

### Kubernetes CronJob

**File**: `manifests/bucketeer/charts/batch/templates/feature-auto-archiver-cronjob.yaml` (NEW)

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: {{ include "batch.fullname" . }}-feature-auto-archiver
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM UTC
  concurrencyPolicy: Forbid
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: batch
            env:
              - name: BUCKETEER_BATCH_JOB
                value: "FeatureAutoArchiver"
            resources:
              requests:
                cpu: 100m
                memory: 256Mi
              limits:
                cpu: 500m
                memory: 512Mi
```
