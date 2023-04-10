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

# Important Notice

* Users can use Progressive Rollout when the number of variations is less than or equal to 2
* Users can't use same scheduled time in single auto ops rules. For example, users can not set true for 50% at 2023-01-01 00:06:00 and 80% at the same time.

# Processing flow

The following image is a processing flow in this feature.

![processing-flow](./images/0041-image3.png)

For instance, Web Client sends clause as follows:

```go
&autoopsproto.ProgressiveRolloutClause{
	// The another varition id is vid-2
	VariationId: "vid-1",
	Schedules: []*autoopsproto.ProgressiveRolloutClause_Schedule{
		{
			// '2023-01-01 00:03:00'
			Time: 1672498980,
			Weight: 20000,
		},
		{
			// '2023-01-01 00:06:00'
			Time: 1672499160,
			Weight: 40000,
		},
		{
			// '2023-01-01 00:09:00'
			Time: 1672499340,
			Weight: 60000,
		},
	},
}
```

1. Web Client registers the above clauses by calling `CreateAutopOps` rules
2. Batch service call `ListAutoOpsRules`, and check if the current time is a scheduled time. In this case, it checks whether the current time is 2023-01-01 00:03:00.
3. If the current time is a scheduled time, Batch service call `ExecuteAutoOps`.
4. AutoOps service call`UpdateFeatureTargeting` to update feature rules. In this case, update the weight of vid-1 to 20000 and the weight of vid-2 to 80000.

# Changes

## Proto

### Progressive rollout

`Interval` field in `ProgressiveRolloutClause` is filled by client side when UI is Template Setting.
When we call `ExecuteAutoOps`, we'll send `time` field. AutoOps service will change `executed` flag of the time match to `true`.

```diff
diff --git a/proto/autoops/auto_ops_rule.proto b/proto/autoops/auto_ops_rule.proto
index 4812ee1..77d727d 100644
--- a/proto/autoops/auto_ops_rule.proto
+++ b/proto/autoops/auto_ops_rule.proto
@@ -33,4 +33,5 @@ message AutoOpsRule {
 enum OpsType {
   ENABLE_FEATURE = 0;
   DISABLE_FEATURE = 1;
+  PROGRESSIVE_ROLLOUT = 2;
 }
diff --git a/proto/autoops/clause.proto b/proto/autoops/clause.proto
index 8c17f23..1624a2a 100644
--- a/proto/autoops/clause.proto
+++ b/proto/autoops/clause.proto
 message OpsEventRateClause {
@@ -58,3 +64,21 @@ message WebhookClause {
   string webhook_id = 1;
   repeated Condition conditions = 2;
 }
+
+message ProgressiveRolloutClause {
+  message Schedule {
+    int64 time = 1;
+    int32 weight = 2;
+    bool executed = 3;
+  }
+  // Interval field is optional.
+  enum Interval {
+    UNKNOWN = 0;
+    HOURLY = 1;
+    DAILY = 2;
+    WEEKLY = 3;
+  }
+  string variation_id = 1;
+  repeated Schedule schedules = 2;
+  Interval interval = 3;
+}
diff --git a/proto/autoops/command.proto b/proto/autoops/command.proto
index fda7a78..e1025b9 100644
--- a/proto/autoops/command.proto
+++ b/proto/autoops/command.proto
@@ -26,6 +26,7 @@ message CreateAutoOpsRuleCommand {
   repeated OpsEventRateClause ops_event_rate_clauses = 3;
   repeated DatetimeClause datetime_clauses = 4;
   repeated WebhookClause webhook_clauses = 5;
+  ProgressiveRolloutClause progressive_rollout_clause = 6;
 }
 
 message ChangeAutoOpsRuleOpsTypeCommand {
@@ -36,6 +37,10 @@ message DeleteAutoOpsRuleCommand {}
 
 message ChangeAutoOpsRuleTriggeredAtCommand {}
 
+message ChangeAutoOpsRuleExecutedCommand {
+  int64 time = 1;
+}
+
 message AddOpsEventRateClauseCommand {
   OpsEventRateClause ops_event_rate_clause = 1;
 }
@@ -81,3 +86,12 @@ message ChangeWebhookClauseCommand {
   string id = 1;
   WebhookClause webhook_clause = 2;
 }
+
+message AddProgressiveRolloutClauseCommand {
+  ProgressiveRolloutClause progressive_rollout_clause = 1;
+}
+
+message ChangeProgressiveRolloutClauseCommand {
+  string id = 1;
+  ProgressiveRolloutClause progressive_rollout_clause = 2;
+}
diff --git a/proto/autoops/service.proto b/proto/autoops/service.proto
index 1d870fe..25bd2c7 100644
--- a/proto/autoops/service.proto
+++ b/proto/autoops/service.proto
@@ -70,6 +70,8 @@ message UpdateAutoOpsRuleRequest {
   repeated ChangeDatetimeClauseCommand change_datetime_clause_commands = 8;
   repeated AddWebhookClauseCommand add_webhook_clause_commands = 9;
   repeated ChangeWebhookClauseCommand change_webhook_clause_commands = 10;
+  repeated AddProgressiveRolloutClauseCommand add_progressive_rollout_clause_commands = 11;
+  repeated ChangeProgressiveRolloutClauseCommand change_progressive_rollout_clause_commands = 12;
 }
 
 message UpdateAutoOpsRuleResponse {}
@@ -79,6 +81,8 @@ message ExecuteAutoOpsRequest {
   string id = 2;
   ChangeAutoOpsRuleTriggeredAtCommand
       change_auto_ops_rule_triggered_at_command = 3;
+  ChangeAutoOpsRuleExecutedCommand
+      change_auto_ops_rule_executed_command = 4;
 }
 
 message ExecuteAutoOpsResponse {
diff --git a/proto/event/domain/event.proto b/proto/event/domain/event.proto
index 29e58d9..33fdb87 100644
--- a/proto/event/domain/event.proto
+++ b/proto/event/domain/event.proto
@@ -879,3 +879,13 @@ message WebhookClauseChangedEvent {
   string clause_id = 1;
   bucketeer.autoops.WebhookClause webhook_clause = 2;
 }
+
+message ProgressiveRolloutClauseAddedEvent {
+  string clause_id = 1;
+  bucketeer.autoops.ProgressiveRolloutClause oprogressive_rollout_clause = 2;
+}
+
+message ProgressiveRolloutClauseChangedEvent {
+  string clause_id = 1;
+  bucketeer.autoops.ProgressiveRolloutClause progressive_rollout_clause = 2;
+}

```

### Other

In addition to this feature, we'll introduce `Type` field to `Clause` field. This will be used when filtering rules.

```diff
diff --git a/proto/autoops/clause.proto b/proto/autoops/clause.proto
index 8c17f23..af9dabf 100644
--- a/proto/autoops/clause.proto
+++ b/proto/autoops/clause.proto
@@ -22,6 +22,12 @@ import "google/protobuf/any.proto";
 message Clause {
   string id = 1;
   google.protobuf.Any clause = 2;
+  enum Type {
+    OPS_EVENT_RATE = 0;
+    DATE_TIME = 1;
+    WEBHOOK = 2;
+    PROGRESSIVE_ROLLOUT = 3;
+  }
 }
 ```

**NOTE**
We can't define `Type` in `AutoOpsRule` instead of `Clause` because `Clause` field in `AutoOpsRule` is an array.
There is a possibility that multiple `Clause` types are included.

## Backend Changes

### Major changes

* batch/job/progressive_rollout_watcher.go
    * This watcher updates feature rules at the scheduled time.
* batch/executor/rollout_updater.go
	* This is the executor which sends ChangeAutoOpsRuleExecutedCommand to AutoOpsRule service.

### Minor changes

* pkg/autoops/api/operation.go
    * We need to modify `ExecuteOperation`. When OpsType field in AutoOpsRule is PROGRESSIVE_ROLLOUT, we'll call `UpdateFeatureTargeting`.
* pkg/opsevent/batch/executor/executor.go
	* Rename executor.go to flag_triggerer.go
	* Then abstract executor.go such as https://github.com/bucketeer-io/bucketeer/blob/main/pkg/eventpersisterdwh/persister/event.go.
* pkg/autoops/command/auto_ops_rule.go
	* Add new commands to `Handle` func.

## Hack

The following procedure is for adding type field.

1. Get all projects through `ListProjects`.
2. Get all environment namespaces by project ids through `ListEnvironments`.
3. Run the transaction to prevent creating new tags while running the hack.
4. Get all rules by environment namespaces through `ListAutoOpsRules`.
5. Update rules through `UpdateAutoOpsRule`.
