
�
proto/experiment/command.protobucketeer.experiment"Y
CreateGoalCommand
id (	Rid
name (	Rname 
description (	Rdescription"'
RenameGoalCommand
name (	Rname"@
ChangeDescriptionGoalCommand 
description (	Rdescription"
ArchiveGoalCommand"
DeleteGoalCommand"�
CreateExperimentCommand

feature_id (	R	featureId
start_at (RstartAt
stop_at (RstopAt
goal_ids (	RgoalIds
name (	Rname 
description (	Rdescription*
base_variation_id (	RbaseVariationIdJ"S
ChangeExperimentPeriodCommand
start_at (RstartAt
stop_at (RstopAt"1
ChangeExperimentNameCommand
name (	Rname"F
"ChangeExperimentDescriptionCommand 
description (	Rdescription"
StopExperimentCommand"
ArchiveExperimentCommand"
DeleteExperimentCommand"
StartExperimentCommand"
FinishExperimentCommandB4Z2github.com/bucketeer-io/bucketeer/proto/experimentbproto3
�
proto/feature/variation.protobucketeer.feature"g
	Variation
id (	Rid
value (	Rvalue
name (	Rname 
description (	RdescriptionB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
!proto/experiment/experiment.protobucketeer.experimentproto/feature/variation.proto"�

Experiment
id (	Rid
goal_id (	BRgoalId

feature_id (	R	featureId'
feature_version (RfeatureVersion<

variations (2.bucketeer.feature.VariationR
variations
start_at (RstartAt
stop_at (RstopAt
stopped (BRstopped!

stopped_at	 (B0R	stoppedAt

created_at
 (R	createdAt

updated_at (R	updatedAt
deleted (Rdeleted
goal_ids (	RgoalIds
name (	Rname 
description (	Rdescription*
base_variation_id (	RbaseVariationId?
status (2'.bucketeer.experiment.Experiment.StatusRstatus

maintainer (	R
maintainer
archived (Rarchived"B
Status
WAITING 
RUNNING
STOPPED
FORCE_STOPPEDJ"Q
ExperimentsB
experiments (2 .bucketeer.experiment.ExperimentRexperimentsB4Z2github.com/bucketeer-io/bucketeer/proto/experimentbproto3
�
proto/experiment/goal.protobucketeer.experiment"�
Goal
id (	Rid
name (	Rname 
description (	Rdescription
deleted (Rdeleted

created_at (R	createdAt

updated_at (R	updatedAt'
is_in_use_status (RisInUseStatus
archived (RarchivedB4Z2github.com/bucketeer-io/bucketeer/proto/experimentbproto3
�
google/protobuf/wrappers.protogoogle.protobuf"#
DoubleValue
value (Rvalue""

FloatValue
value (Rvalue""

Int64Value
value (Rvalue"#
UInt64Value
value (Rvalue""

Int32Value
value (Rvalue"#
UInt32Value
value (Rvalue"!
	BoolValue
value (Rvalue"#
StringValue
value (	Rvalue""

BytesValue
value (RvalueB�
com.google.protobufBWrappersProtoPZ1google.golang.org/protobuf/types/known/wrapperspb��GPB�Google.Protobuf.WellKnownTypesbproto3
�2
proto/experiment/service.protobucketeer.experimentgoogle/protobuf/wrappers.protoproto/experiment/command.protoproto/experiment/goal.proto!proto/experiment/experiment.proto"U
GetGoalRequest
id (	Rid3
environment_namespace (	RenvironmentNamespace"A
GetGoalResponse.
goal (2.bucketeer.experiment.GoalRgoal"�
ListGoalsRequest
	page_size (RpageSize
cursor (	Rcursor3
environment_namespace (	RenvironmentNamespaceI
order_by (2..bucketeer.experiment.ListGoalsRequest.OrderByRorderBy^
order_direction (25.bucketeer.experiment.ListGoalsRequest.OrderDirectionRorderDirection%
search_keyword (	RsearchKeywordC
is_in_use_status (2.google.protobuf.BoolValueRisInUseStatus6
archived (2.google.protobuf.BoolValueRarchived"@
OrderBy
DEFAULT 
NAME

CREATED_AT

UPDATED_AT"#
OrderDirection
ASC 
DESC"~
ListGoalsResponse0
goals (2.bucketeer.experiment.GoalRgoals
cursor (	Rcursor
total_count (R
totalCount"�
CreateGoalRequestA
command (2'.bucketeer.experiment.CreateGoalCommandRcommand3
environment_namespace (	RenvironmentNamespace"
CreateGoalResponse"�
ArchiveGoalRequest
id (	RidB
command (2(.bucketeer.experiment.ArchiveGoalCommandRcommand3
environment_namespace (	RenvironmentNamespace"
ArchiveGoalResponse"�
DeleteGoalRequest
id (	RidA
command (2'.bucketeer.experiment.DeleteGoalCommandRcommand3
environment_namespace (	RenvironmentNamespace"
DeleteGoalResponse"�
UpdateGoalRequest
id (	RidN
rename_command (2'.bucketeer.experiment.RenameGoalCommandRrenameCommandp
change_description_command (22.bucketeer.experiment.ChangeDescriptionGoalCommandRchangeDescriptionCommand3
environment_namespace (	RenvironmentNamespace"
UpdateGoalResponse"[
GetExperimentRequest
id (	Rid3
environment_namespace (	RenvironmentNamespace"Y
GetExperimentResponse@

experiment (2 .bucketeer.experiment.ExperimentR
experiment"�
ListExperimentsRequest

feature_id (	R	featureIdD
feature_version (2.google.protobuf.Int32ValueRfeatureVersion
from (Rfrom
to (Rto
	page_size (RpageSize
cursor (	Rcursor3
environment_namespace (	RenvironmentNamespace3
status (2.google.protobuf.Int32ValueRstatus

maintainer	 (	R
maintainerO
order_by
 (24.bucketeer.experiment.ListExperimentsRequest.OrderByRorderByd
order_direction (2;.bucketeer.experiment.ListExperimentsRequest.OrderDirectionRorderDirection%
search_keyword (	RsearchKeyword6
archived (2.google.protobuf.BoolValueRarchivedC
statuses (2'.bucketeer.experiment.Experiment.StatusRstatuses"@
OrderBy
DEFAULT 
NAME

CREATED_AT

UPDATED_AT"#
OrderDirection
ASC 
DESC"�
ListExperimentsResponseB
experiments (2 .bucketeer.experiment.ExperimentRexperiments
cursor (	Rcursor
total_count (R
totalCount"�
CreateExperimentRequestG
command (2-.bucketeer.experiment.CreateExperimentCommandRcommand3
environment_namespace (	RenvironmentNamespace"\
CreateExperimentResponse@

experiment (2 .bucketeer.experiment.ExperimentR
experiment"�
UpdateExperimentRequest
id (	Rid3
environment_namespace (	RenvironmentNamespace|
 change_experiment_period_command (23.bucketeer.experiment.ChangeExperimentPeriodCommandRchangeExperimentPeriodCommanda
change_name_command (21.bucketeer.experiment.ChangeExperimentNameCommandRchangeNameCommandv
change_description_command (28.bucketeer.experiment.ChangeExperimentDescriptionCommandRchangeDescriptionCommandJJ"
UpdateExperimentResponse"�
StartExperimentRequest3
environment_namespace (	RenvironmentNamespace
id (	RidF
command (2,.bucketeer.experiment.StartExperimentCommandRcommand"
StartExperimentResponse"�
FinishExperimentRequest3
environment_namespace (	RenvironmentNamespace
id (	RidG
command (2-.bucketeer.experiment.FinishExperimentCommandRcommand"
FinishExperimentResponse"�
StopExperimentRequest
id (	RidE
command (2+.bucketeer.experiment.StopExperimentCommandRcommand3
environment_namespace (	RenvironmentNamespace"
StopExperimentResponse"�
ArchiveExperimentRequest
id (	RidH
command (2..bucketeer.experiment.ArchiveExperimentCommandRcommand3
environment_namespace (	RenvironmentNamespace"
ArchiveExperimentResponse"�
DeleteExperimentRequest
id (	RidG
command (2-.bucketeer.experiment.DeleteExperimentCommandRcommand3
environment_namespace (	RenvironmentNamespace"
DeleteExperimentResponse2�
ExperimentServiceX
GetGoal$.bucketeer.experiment.GetGoalRequest%.bucketeer.experiment.GetGoalResponse" ^
	ListGoals&.bucketeer.experiment.ListGoalsRequest'.bucketeer.experiment.ListGoalsResponse" a

CreateGoal'.bucketeer.experiment.CreateGoalRequest(.bucketeer.experiment.CreateGoalResponse" a

UpdateGoal'.bucketeer.experiment.UpdateGoalRequest(.bucketeer.experiment.UpdateGoalResponse" d
ArchiveGoal(.bucketeer.experiment.ArchiveGoalRequest).bucketeer.experiment.ArchiveGoalResponse" a

DeleteGoal'.bucketeer.experiment.DeleteGoalRequest(.bucketeer.experiment.DeleteGoalResponse" j
GetExperiment*.bucketeer.experiment.GetExperimentRequest+.bucketeer.experiment.GetExperimentResponse" p
ListExperiments,.bucketeer.experiment.ListExperimentsRequest-.bucketeer.experiment.ListExperimentsResponse" s
CreateExperiment-.bucketeer.experiment.CreateExperimentRequest..bucketeer.experiment.CreateExperimentResponse" s
UpdateExperiment-.bucketeer.experiment.UpdateExperimentRequest..bucketeer.experiment.UpdateExperimentResponse" p
StartExperiment,.bucketeer.experiment.StartExperimentRequest-.bucketeer.experiment.StartExperimentResponse" s
FinishExperiment-.bucketeer.experiment.FinishExperimentRequest..bucketeer.experiment.FinishExperimentResponse" m
StopExperiment+.bucketeer.experiment.StopExperimentRequest,.bucketeer.experiment.StopExperimentResponse" v
ArchiveExperiment..bucketeer.experiment.ArchiveExperimentRequest/.bucketeer.experiment.ArchiveExperimentResponse" s
DeleteExperiment-.bucketeer.experiment.DeleteExperimentRequest..bucketeer.experiment.DeleteExperimentResponse" B4Z2github.com/bucketeer-io/bucketeer/proto/experimentbproto3