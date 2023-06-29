# Summary

# UI

We'll introduce two kinds of UI, "Manual Setting" and "Template Setting".

## Manual Setting

Users can determine any percentage and schedule in this UI.

![manual-setting-proto-ui](./images/0041-image1.png)

## Template Setting

Manual Setting has a lot of flexibility, but it’s hard when configuring a periodic increase in the percentage of a variation.
In this case, Template Setting is useful.

![template-setting-proto-ui](./images/0041-image2.png)

## Progress bar of rollout

![progress-bar-of-rollout](./images/0041-image4.png)

# Important Notice

* Users can use Progressive Rollout when the number of variations is equal to 2.
* Users can't use same scheduled time in single auto ops rules. For example, users can not set true for 50% at 2023-01-01 00:06:00 and 80% at the same time.
* The interval of time for each scheduled time must be at least 5 minutes.
* We do not support the feature to stop Progressive Rollout temporary. We might support it in the feature.
* Users cannot the update the feature while running Progressive Rollout.
* Users can use both Progressive Rollout and Feature Flag Trigger in Event Rate mode at the same time, but not in other mode such as Datetime mode.
* Users cannot configure the Progressive Rollout when the feature flag is disabled.
* If the operation of Progressive Rollout is deleted, the operation is stopped.
* The operation of Progressive Rollout can not be modified after creating it.

# Processing flow

The following image is a processing flow in this feature.

![processing-flow](./images/0041-image3.png)

For instance, Web Client sends clause as follows:

```go
&autoopsproto.ProgressiveRolloutClause{
	// The another varition id is vid-2
	VariationId: "vid-1",
	Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
		{
			// '2023-01-01 03:00:00'
			Time: 1672509600,
			Weight: 20000,
		},
		{
			// '2023-01-01 06:00:00'
			Time: 1672520400,
			Weight: 40000,
		},
		{
			// '2023-01-01 09:00:00'
			Time: 1672531200,
			Weight: 60000,
		},
	},
}
```

1. Web Client registers the above clauses by calling `CreateProgressiveRollout` rules
2. Batch service calls `ListProgressiveRollout`, and check if the current time is a scheduled time. In this case, it checks whether the current time is 2023-01-01 00:03:00.
3. If the current time is a scheduled time and the rule is enabled, Batch service call `ExecuteProgressiveRollout`.
4. AutoOps service calls `UpdateFeatureTargeting` to update feature rules. In this case, update the weight of vid-1 to 20000 and the weight of vid-2 to 80000.

# Changes

## Table

We'll create `ops_progressive_rollout` table as follows. `ProgressiveRolloutManualScheduleClause` and `ProgressiveRolloutAutomaticScheduleClause` are converted into Any type and stored into `clause` column.

```sql
CREATE TABLE IF NOT EXISTS `ops_progressive_rollout` (
  `id` VARCHAR(255) NOT NULL,
  `feature_id` VARCHAR(255) NOT NULL,
  `clause` JSON NOT NULL,
  `status` INT(11) NOT NULL,
  `type` INT(11) NOT NULL,
  `created_at` BIGINT(20) NOT NULL,
  `updated_at` BIGINT(20) NOT NULL,
  `environment_namespace` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`, `environment_namespace`),
  CONSTRAINT `foreign_progressive_rollout_feature_id_environment_namespace`
    FOREIGN KEY (`feature_id`, `environment_namespace`)
    REFERENCES `feature` (`id`, `environment_namespace`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);
```

## Proto

### Progressive rollout

* proto/autoops/progressive_rollout.proto

`Status` represents the state of operations for ProgressiveRollout. If a operation has not started, `Status` is WAITING.
If a operation is in progress, `Status` is DOING. If a operaiton is finished, `Status` is DONE.
All fields in `ProgressiveRollout` are stored to DB as columns.
The reason of using `google.protobuf.Any` is for ease of expansion.

```proto
message ProgressiveRollout {
  enum Type {
    MANUAL_SCHEDULE = 0;
    TEMPLATE_SCHEDULE = 1;
  }
  enum Status {
    WAITING = 0;
    RUNNING = 1;
    FINISHED = 2;
  }
  string id = 1;
  string feature_id = 2;
  google.protobuf.Any clause = 3;
  Status status = 4;
  int64 created_at = 5;
  int64 updated_at = 6;
  Type type = 7;
}
```

* proto/autoops/clause.proto

`ProgressiveRolloutManualScheduleClause` is set when Manual Setting is used by users.
`ProgressiveRolloutAutomaticScheduleClause` is set when Template Setting is used by users.

```proto
message ProgressiveRolloutSchedule {
  string schedule_id = 1;
  int64 execute_at = 2;
  int32 weight = 3;
  int64 triggered_at = 4;
}

message ProgressiveRolloutManualScheduleClause {
  repeated ProgressiveRolloutSchedule schedules = 1;
  string variation_id = 2;
}

message ProgressiveRolloutTemplateScheduleClause {
  enum Interval {
    UNKNOWN = 0;
    HOURLY = 1;
    DAILY = 2;
    WEEKLY = 3;
  }
  // The reason of setting `schedules` is to save `triggered_at` in each
  // schedule.
  repeated ProgressiveRolloutSchedule schedules = 1;
  Interval interval = 2;
  int64 increments = 3;
  string variation_id = 4;
}
```

* proto/autoops/service.proto

Progressive Rollout feature has unique API for creating, updating and deleting.

```proto
message CreateProgressiveRolloutRequest {
  string environment_namespace = 1;
  CreateProgressiveRolloutCommand command = 2;
}

message CreateProgressiveRolloutResponse {}

message GetProgressiveRolloutRequest {
  string environment_namespace = 1;
  string id = 2;
}

message GetProgressiveRolloutResponse {
  ProgressiveRollout progressive_rollout = 1;
}

message DeleteProgressiveRolloutRequest {
  string environment_namespace = 1;
  string id = 2;
  DeleteProgressiveRolloutCommand command = 3;
}

message DeleteProgressiveRolloutResponse {}

message ListProgressiveRolloutsRequest {
  enum OrderBy {
    DEFAULT = 0;
    CREATED_AT = 1;
    UPDATED_AT = 2;
  }
  enum OrderDirection {
    ASC = 0;
    DESC = 1;
  }
  string environment_namespace = 1;
  int64 page_size = 2;
  string cursor = 3;
  repeated string feature_ids = 4;
  OrderBy order_by = 5;
  OrderDirection order_direction = 6;
  optional ProgressiveRollout.Status status = 7;
  optional ProgressiveRollout.Type type = 8;
}

message ListProgressiveRolloutsResponse {
  repeated ProgressiveRollout progressive_rollouts = 1;
  string cursor = 2;
  int64 total_count = 3;
}

message ExecuteProgressiveRolloutRequest {
  string environment_namespace = 1;
  string id = 2;
  ChangeProgressiveRolloutScheduleTriggeredAtCommand
      change_progressive_rollout_triggered_at_command = 3;
}

message ExecuteProgressiveRolloutResponse {}
```

* proto/autoops/command.proto

```proto
message CreateProgressiveRolloutCommand {
  string feature_id = 1;
  optional ProgressiveRolloutManualScheduleClause
      progressive_rollout_manual_schedule_clause = 2;
  optional ProgressiveRolloutTemplateScheduleClause
      progressive_rollout_template_schedule_clause = 3;
}

message DeleteProgressiveRolloutCommand {}

message AddProgressiveRolloutManualScheduleClauseCommand {
  ProgressiveRolloutManualScheduleClause clause = 1;
}

message AddProgressiveRolloutTemplateScheduleClauseCommand {
  ProgressiveRolloutTemplateScheduleClause clause = 1;
}

message ChangeProgressiveRolloutScheduleTriggeredAtCommand {
  string schedule_id = 1;
}
```

## Backend Changes

### Major changes

* batch/job/progressive_rollout_watcher.go
    * This watcher updates feature rules at the scheduled time.
* batch/executor/rollout_updater.go
	* This is the executor which sends ChangeProgressiveRolloutTriggeredAtCommand to AutoOpsRule service.
* storage/v2/progressive_rollout.go
	* This is for inserting data into `ops_progressive_rollout` table.
* pkg/autoops/api/progressive_rollout.go
* pkg/autoops/command/progressive_rollout.go
* pkg/autoops/domain/progressive_rollout.go

### Minor changes

* pkg/opsevent/batch/executor/executor.go
	* Rename executor.go to flag_triggerer.go
	* Then abstract executor.go such as https://github.com/bucketeer-io/bucketeer/blob/main/pkg/eventpersisterdwh/persister/event.go.
* pkg/feature/api/feature.go
	* We need to modify `UpdateFeatureVariations`. We need to validate changing variations. Return error if Progressive Rollout is running.
