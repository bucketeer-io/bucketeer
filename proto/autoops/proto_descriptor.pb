
�
google/protobuf/any.protogoogle.protobuf"6
Any
type_url (	RtypeUrl
value (RvalueBv
com.google.protobufBAnyProtoPZ,google.golang.org/protobuf/types/known/anypb�GPB�Google.Protobuf.WellKnownTypesbproto3
�
proto/autoops/clause.protobucketeer.autoopsgoogle/protobuf/any.proto"�
Clause
id (	Rid,
clause (2.google.protobuf.AnyRclause>
action_type (2.bucketeer.autoops.ActionTypeR
actionType"�
OpsEventRateClause!
variation_id (	RvariationId
goal_id (	RgoalId
	min_count (RminCount)
threadshold_rate (RthreadsholdRateJ
operator (2..bucketeer.autoops.OpsEventRateClause.OperatorRoperator>
action_type (2.bucketeer.autoops.ActionTypeR
actionType"3
Operator
GREATER_OR_EQUAL 
LESS_OR_EQUALJ"d
DatetimeClause
time (Rtime>
action_type (2.bucketeer.autoops.ActionTypeR
actionType"�
ProgressiveRolloutSchedule
schedule_id (	R
scheduleId

execute_at (R	executeAt
weight (Rweight!
triggered_at (RtriggeredAt"�
&ProgressiveRolloutManualScheduleClauseK
	schedules (2-.bucketeer.autoops.ProgressiveRolloutScheduleR	schedules!
variation_id (	RvariationId"�
(ProgressiveRolloutTemplateScheduleClauseK
	schedules (2-.bucketeer.autoops.ProgressiveRolloutScheduleR	schedules`
interval (2D.bucketeer.autoops.ProgressiveRolloutTemplateScheduleClause.IntervalRinterval

increments (R
increments!
variation_id (	RvariationId":
Interval
UNKNOWN 

HOURLY	
DAILY

WEEKLY*2

ActionType
UNKNOWN 

ENABLE
DISABLEB1Z/github.com/bucketeer-io/bucketeer/proto/autoopsbproto3
�
!proto/autoops/auto_ops_rule.protobucketeer.autoopsproto/autoops/clause.proto"�
AutoOpsRule
id (	Rid

feature_id (	R	featureId5
ops_type (2.bucketeer.autoops.OpsTypeRopsType3
clauses (2.bucketeer.autoops.ClauseRclauses

created_at (R	createdAt

updated_at (R	updatedAt
deleted	 (RdeletedH
auto_ops_status
 (2 .bucketeer.autoops.AutoOpsStatusRautoOpsStatusJ"T
AutoOpsRulesD
auto_ops_rules (2.bucketeer.autoops.AutoOpsRuleRautoOpsRules*?
OpsType
TYPE_UNKNOWN 
SCHEDULE

EVENT_RATE"*D
AutoOpsStatus
WAITING 
RUNNING
FINISHED
STOPPEDB1Z/github.com/bucketeer-io/bucketeer/proto/autoopsbproto3
�
'proto/autoops/progressive_rollout.protobucketeer.autoopsgoogle/protobuf/any.proto"�
ProgressiveRollout
id (	Rid

feature_id (	R	featureId,
clause (2.google.protobuf.AnyRclauseD
status (2,.bucketeer.autoops.ProgressiveRollout.StatusRstatus

created_at (R	createdAt

updated_at (R	updatedAt>
type (2*.bucketeer.autoops.ProgressiveRollout.TypeRtypeN

stopped_by (2/.bucketeer.autoops.ProgressiveRollout.StoppedByR	stoppedBy

stopped_at	 (R	stoppedAt"2
Type
MANUAL_SCHEDULE 
TEMPLATE_SCHEDULE"=
Status
WAITING 
RUNNING
FINISHED
STOPPED"I
	StoppedBy
UNKNOWN 
USER
OPS_SCHEDULE
OPS_KILL_SWITCHB1Z/github.com/bucketeer-io/bucketeer/proto/autoopsbproto3
�
proto/autoops/command.protobucketeer.autoops!proto/autoops/auto_ops_rule.protoproto/autoops/clause.proto'proto/autoops/progressive_rollout.proto"�
CreateAutoOpsRuleCommand

feature_id (	R	featureId5
ops_type (2.bucketeer.autoops.OpsTypeRopsTypeZ
ops_event_rate_clauses (2%.bucketeer.autoops.OpsEventRateClauseRopsEventRateClausesL
datetime_clauses (2!.bucketeer.autoops.DatetimeClauseRdatetimeClauses"
DeleteAutoOpsRuleCommand"
StopAutoOpsRuleCommand"V
ChangeAutoOpsStatusCommand8
status (2 .bucketeer.autoops.AutoOpsStatusRstatus"8
ExecuteAutoOpsRuleCommand
	clause_id (	RclauseId"x
AddOpsEventRateClauseCommandX
ops_event_rate_clause (2%.bucketeer.autoops.OpsEventRateClauseRopsEventRateClause"�
ChangeOpsEventRateClauseCommand
id (	RidX
ops_event_rate_clause (2%.bucketeer.autoops.OpsEventRateClauseRopsEventRateClause"%
DeleteClauseCommand
id (	Rid"f
AddDatetimeClauseCommandJ
datetime_clause (2!.bucketeer.autoops.DatetimeClauseRdatetimeClause"y
ChangeDatetimeClauseCommand
id (	RidJ
datetime_clause (2!.bucketeer.autoops.DatetimeClauseRdatetimeClause"�
CreateProgressiveRolloutCommand

feature_id (	R	featureId�
*progressive_rollout_manual_schedule_clause (29.bucketeer.autoops.ProgressiveRolloutManualScheduleClauseH R&progressiveRolloutManualScheduleClause��
,progressive_rollout_template_schedule_clause (2;.bucketeer.autoops.ProgressiveRolloutTemplateScheduleClauseHR(progressiveRolloutTemplateScheduleClause�B-
+_progressive_rollout_manual_schedule_clauseB/
-_progressive_rollout_template_schedule_clause"o
StopProgressiveRolloutCommandN

stopped_by (2/.bucketeer.autoops.ProgressiveRollout.StoppedByR	stoppedBy"!
DeleteProgressiveRolloutCommand"�
0AddProgressiveRolloutManualScheduleClauseCommandQ
clause (29.bucketeer.autoops.ProgressiveRolloutManualScheduleClauseRclause"�
2AddProgressiveRolloutTemplateScheduleClauseCommandS
clause (2;.bucketeer.autoops.ProgressiveRolloutTemplateScheduleClauseRclause"U
2ChangeProgressiveRolloutScheduleTriggeredAtCommand
schedule_id (	R
scheduleIdB1Z/github.com/bucketeer-io/bucketeer/proto/autoopsbproto3
�
proto/autoops/ops_count.protobucketeer.autoops"�
OpsCount
id (	Rid'
auto_ops_rule_id (	RautoOpsRuleId
	clause_id (	RclauseId

updated_at (R	updatedAt&
ops_event_count (RopsEventCount)
evaluation_count (RevaluationCount

feature_id (	R	featureIdB1Z/github.com/bucketeer-io/bucketeer/proto/autoopsbproto3
�1
proto/autoops/service.protobucketeer.autoops!proto/autoops/auto_ops_rule.protoproto/autoops/command.protoproto/autoops/ops_count.proto'proto/autoops/progressive_rollout.proto"\
GetAutoOpsRuleRequest3
environment_namespace (	RenvironmentNamespace
id (	Rid"\
GetAutoOpsRuleResponseB
auto_ops_rule (2.bucketeer.autoops.AutoOpsRuleRautoOpsRule"�
CreateAutoOpsRuleRequest3
environment_namespace (	RenvironmentNamespaceE
command (2+.bucketeer.autoops.CreateAutoOpsRuleCommandRcommand"
CreateAutoOpsRuleResponse"�
ListAutoOpsRulesRequest3
environment_namespace (	RenvironmentNamespace
	page_size (RpageSize
cursor (	Rcursor
feature_ids (	R
featureIds"x
ListAutoOpsRulesResponseD
auto_ops_rules (2.bucketeer.autoops.AutoOpsRuleRautoOpsRules
cursor (	Rcursor"�
StopAutoOpsRuleRequest3
environment_namespace (	RenvironmentNamespace
id (	RidC
command (2).bucketeer.autoops.StopAutoOpsRuleCommandRcommand"
StopAutoOpsRuleResponse"�
DeleteAutoOpsRuleRequest3
environment_namespace (	RenvironmentNamespace
id (	RidE
command (2+.bucketeer.autoops.DeleteAutoOpsRuleCommandRcommand"
DeleteAutoOpsRuleResponse"�
UpdateAutoOpsRuleRequest3
environment_namespace (	RenvironmentNamespace
id (	Ridz
"add_ops_event_rate_clause_commands (2/.bucketeer.autoops.AddOpsEventRateClauseCommandRaddOpsEventRateClauseCommands�
%change_ops_event_rate_clause_commands (22.bucketeer.autoops.ChangeOpsEventRateClauseCommandR changeOpsEventRateClauseCommands\
delete_clause_commands (2&.bucketeer.autoops.DeleteClauseCommandRdeleteClauseCommandsl
add_datetime_clause_commands (2+.bucketeer.autoops.AddDatetimeClauseCommandRaddDatetimeClauseCommandsu
change_datetime_clause_commands (2..bucketeer.autoops.ChangeDatetimeClauseCommandRchangeDatetimeClauseCommandsJ"
UpdateAutoOpsRuleResponse"�
ExecuteAutoOpsRequest3
environment_namespace (	RenvironmentNamespace
id (	Ridn
execute_auto_ops_rule_command (2,.bucketeer.autoops.ExecuteAutoOpsRuleCommandRexecuteAutoOpsRuleCommandJ"E
ExecuteAutoOpsResponse+
already_triggered (RalreadyTriggered"�
ListOpsCountsRequest3
environment_namespace (	RenvironmentNamespace
	page_size (RpageSize
cursor (	Rcursor)
auto_ops_rule_ids (	RautoOpsRuleIds
feature_ids (	R
featureIds"k
ListOpsCountsResponse
cursor (	Rcursor:

ops_counts (2.bucketeer.autoops.OpsCountR	opsCounts"�
CreateProgressiveRolloutRequest3
environment_namespace (	RenvironmentNamespaceL
command (22.bucketeer.autoops.CreateProgressiveRolloutCommandRcommand""
 CreateProgressiveRolloutResponse"c
GetProgressiveRolloutRequest3
environment_namespace (	RenvironmentNamespace
id (	Rid"w
GetProgressiveRolloutResponseV
progressive_rollout (2%.bucketeer.autoops.ProgressiveRolloutRprogressiveRollout"�
StopProgressiveRolloutRequest3
environment_namespace (	RenvironmentNamespace
id (	RidJ
command (20.bucketeer.autoops.StopProgressiveRolloutCommandRcommand" 
StopProgressiveRolloutResponse"�
DeleteProgressiveRolloutRequest3
environment_namespace (	RenvironmentNamespace
id (	RidL
command (22.bucketeer.autoops.DeleteProgressiveRolloutCommandRcommand""
 DeleteProgressiveRolloutResponse"�
ListProgressiveRolloutsRequest3
environment_namespace (	RenvironmentNamespace
	page_size (RpageSize
cursor (	Rcursor
feature_ids (	R
featureIdsT
order_by (29.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderByRorderByi
order_direction (2@.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderDirectionRorderDirectionI
status (2,.bucketeer.autoops.ProgressiveRollout.StatusH Rstatus�C
type (2*.bucketeer.autoops.ProgressiveRollout.TypeHRtype�"6
OrderBy
DEFAULT 

CREATED_AT

UPDATED_AT"#
OrderDirection
ASC 
DESCB	
_statusB
_type"�
ListProgressiveRolloutsResponseX
progressive_rollouts (2%.bucketeer.autoops.ProgressiveRolloutRprogressiveRollouts
cursor (	Rcursor
total_count (R
totalCount"�
 ExecuteProgressiveRolloutRequest3
environment_namespace (	RenvironmentNamespace
id (	Rid�
/change_progressive_rollout_triggered_at_command (2E.bucketeer.autoops.ChangeProgressiveRolloutScheduleTriggeredAtCommandR*changeProgressiveRolloutTriggeredAtCommand"#
!ExecuteProgressiveRolloutResponse2�
AutoOpsServiceg
GetAutoOpsRule(.bucketeer.autoops.GetAutoOpsRuleRequest).bucketeer.autoops.GetAutoOpsRuleResponse" m
ListAutoOpsRules*.bucketeer.autoops.ListAutoOpsRulesRequest+.bucketeer.autoops.ListAutoOpsRulesResponse" p
CreateAutoOpsRule+.bucketeer.autoops.CreateAutoOpsRuleRequest,.bucketeer.autoops.CreateAutoOpsRuleResponse" j
StopAutoOpsRule).bucketeer.autoops.StopAutoOpsRuleRequest*.bucketeer.autoops.StopAutoOpsRuleResponse" p
DeleteAutoOpsRule+.bucketeer.autoops.DeleteAutoOpsRuleRequest,.bucketeer.autoops.DeleteAutoOpsRuleResponse" p
UpdateAutoOpsRule+.bucketeer.autoops.UpdateAutoOpsRuleRequest,.bucketeer.autoops.UpdateAutoOpsRuleResponse" g
ExecuteAutoOps(.bucketeer.autoops.ExecuteAutoOpsRequest).bucketeer.autoops.ExecuteAutoOpsResponse" d
ListOpsCounts'.bucketeer.autoops.ListOpsCountsRequest(.bucketeer.autoops.ListOpsCountsResponse" �
CreateProgressiveRollout2.bucketeer.autoops.CreateProgressiveRolloutRequest3.bucketeer.autoops.CreateProgressiveRolloutResponse" |
GetProgressiveRollout/.bucketeer.autoops.GetProgressiveRolloutRequest0.bucketeer.autoops.GetProgressiveRolloutResponse" 
StopProgressiveRollout0.bucketeer.autoops.StopProgressiveRolloutRequest1.bucketeer.autoops.StopProgressiveRolloutResponse" �
DeleteProgressiveRollout2.bucketeer.autoops.DeleteProgressiveRolloutRequest3.bucketeer.autoops.DeleteProgressiveRolloutResponse" �
ListProgressiveRollouts1.bucketeer.autoops.ListProgressiveRolloutsRequest2.bucketeer.autoops.ListProgressiveRolloutsResponse" �
ExecuteProgressiveRollout3.bucketeer.autoops.ExecuteProgressiveRolloutRequest4.bucketeer.autoops.ExecuteProgressiveRolloutResponse" B1Z/github.com/bucketeer-io/bucketeer/proto/autoopsbproto3