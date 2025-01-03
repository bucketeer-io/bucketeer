
�
"proto/eventcounter/histogram.protobucketeer.eventcounter"3
	Histogram
hist (Rhist
bins (RbinsB6Z4github.com/bucketeer-io/bucketeer/proto/eventcounterbproto3
�
-proto/eventcounter/distribution_summary.protobucketeer.eventcounter"proto/eventcounter/histogram.proto"�
DistributionSummary
mean (Rmean
sd (Rsd
rhat (Rrhat?
	histogram (2!.bucketeer.eventcounter.HistogramR	histogram
median (Rmedian$
percentile025 (Rpercentile025$
percentile975 (Rpercentile975B6Z4github.com/bucketeer-io/bucketeer/proto/eventcounterbproto3
�
(proto/eventcounter/variation_count.protobucketeer.eventcounter"�
VariationCount!
variation_id (	RvariationId

user_count (R	userCount
event_count (R
eventCount
	value_sum (RvalueSum

created_at (R	createdAt'
variation_value (	RvariationValue4
value_sum_per_user_mean (RvalueSumPerUserMean<
value_sum_per_user_variance (RvalueSumPerUserVarianceB6Z4github.com/bucketeer-io/bucketeer/proto/eventcounterbproto3
�
)proto/eventcounter/evaluation_count.protobucketeer.eventcounter(proto/eventcounter/variation_count.proto"�
EvaluationCount
id (	Rid

feature_id (	R	featureId'
feature_version (RfeatureVersionO
realtime_counts (2&.bucketeer.eventcounter.VariationCountRrealtimeCountsI
batch_counts (2&.bucketeer.eventcounter.VariationCountRbatchCounts

updated_at (R	updatedAtB6Z4github.com/bucketeer-io/bucketeer/proto/eventcounterbproto3
�
)proto/eventcounter/experiment_count.protobucketeer.eventcounter(proto/eventcounter/variation_count.proto"�
ExperimentCount
id (	Rid

feature_id (	R	featureId'
feature_version (RfeatureVersion
goal_id (	BRgoalIdS
realtime_counts (2&.bucketeer.eventcounter.VariationCountBRrealtimeCountsM
batch_counts (2&.bucketeer.eventcounter.VariationCountBRbatchCounts

updated_at (R	updatedAtC
goal_counts (2".bucketeer.eventcounter.GoalCountsR
goalCounts"�

GoalCounts
goal_id (	RgoalIdO
realtime_counts (2&.bucketeer.eventcounter.VariationCountRrealtimeCountsM
batch_counts (2&.bucketeer.eventcounter.VariationCountBRbatchCountsB6Z4github.com/bucketeer-io/bucketeer/proto/eventcounterbproto3
�
#proto/eventcounter/timeseries.protobucketeer.eventcounter"|
VariationTimeseries!
variation_id (	RvariationIdB

timeseries (2".bucketeer.eventcounter.TimeseriesR
timeseries"�

Timeseries

timestamps (R
timestamps
values (Rvalues;
unit (2'.bucketeer.eventcounter.Timeseries.UnitRunit!
total_counts (RtotalCounts"
Unit
HOUR 
DAYB6Z4github.com/bucketeer-io/bucketeer/proto/eventcounterbproto3
�
)proto/eventcounter/variation_result.protobucketeer.eventcounter(proto/eventcounter/variation_count.proto-proto/eventcounter/distribution_summary.proto#proto/eventcounter/timeseries.proto"�
VariationResult!
variation_id (	RvariationIdQ
experiment_count (2&.bucketeer.eventcounter.VariationCountRexperimentCountQ
evaluation_count (2&.bucketeer.eventcounter.VariationCountRevaluationCountO
cvr_prob_best (2+.bucketeer.eventcounter.DistributionSummaryRcvrProbBest`
cvr_prob_beat_baseline (2+.bucketeer.eventcounter.DistributionSummaryRcvrProbBeatBaselineF
cvr_prob (2+.bucketeer.eventcounter.DistributionSummaryRcvrProbk
 evaluation_user_count_timeseries (2".bucketeer.eventcounter.TimeseriesRevaluationUserCountTimeseriesm
!evaluation_event_count_timeseries (2".bucketeer.eventcounter.TimeseriesRevaluationEventCountTimeseries_
goal_user_count_timeseries	 (2".bucketeer.eventcounter.TimeseriesRgoalUserCountTimeseriesa
goal_event_count_timeseries
 (2".bucketeer.eventcounter.TimeseriesRgoalEventCountTimeseries]
goal_value_sum_timeseries (2".bucketeer.eventcounter.TimeseriesRgoalValueSumTimeseriesV
cvr_median_timeseries (2".bucketeer.eventcounter.TimeseriesRcvrMedianTimeseriesd
cvr_percentile025_timeseries (2".bucketeer.eventcounter.TimeseriesRcvrPercentile025Timeseriesd
cvr_percentile975_timeseries (2".bucketeer.eventcounter.TimeseriesRcvrPercentile975TimeseriesI
cvr_timeseries (2".bucketeer.eventcounter.TimeseriesRcvrTimeseriesm
"goal_value_sum_per_user_timeseries (2".bucketeer.eventcounter.TimeseriesRgoalValueSumPerUserTimeseriesj
goal_value_sum_per_user_prob (2+.bucketeer.eventcounter.DistributionSummaryRgoalValueSumPerUserProbs
!goal_value_sum_per_user_prob_best (2+.bucketeer.eventcounter.DistributionSummaryRgoalValueSumPerUserProbBest�
*goal_value_sum_per_user_prob_beat_baseline (2+.bucketeer.eventcounter.DistributionSummaryR#goalValueSumPerUserProbBeatBaselinez
)goal_value_sum_per_user_median_timeseries (2".bucketeer.eventcounter.TimeseriesR#goalValueSumPerUserMedianTimeseries�
0goal_value_sum_per_user_percentile025_timeseries (2".bucketeer.eventcounter.TimeseriesR*goalValueSumPerUserPercentile025Timeseries�
0goal_value_sum_per_user_percentile975_timeseries (2".bucketeer.eventcounter.TimeseriesR*goalValueSumPerUserPercentile975TimeseriesB6Z4github.com/bucketeer-io/bucketeer/proto/eventcounterbproto3
�
$proto/eventcounter/goal_result.protobucketeer.eventcounter)proto/eventcounter/variation_result.proto"{

GoalResult
goal_id (	RgoalIdT
variation_results (2'.bucketeer.eventcounter.VariationResultRvariationResultsB6Z4github.com/bucketeer-io/bucketeer/proto/eventcounterbproto3
�
*proto/eventcounter/experiment_result.protobucketeer.eventcounter$proto/eventcounter/goal_result.proto"�
ExperimentResult
id (	Rid#
experiment_id (	RexperimentId

updated_at (R	updatedAtE
goal_results (2".bucketeer.eventcounter.GoalResultRgoalResultsB6Z4github.com/bucketeer-io/bucketeer/proto/eventcounterbproto3
�
google/protobuf/any.protogoogle.protobuf"6
Any
type_url (	RtypeUrl
value (RvalueBv
com.google.protobufBAnyProtoPZ,google.golang.org/protobuf/types/known/anypb�GPB�Google.Protobuf.WellKnownTypesbproto3
�
google/protobuf/duration.protogoogle.protobuf":
Duration
seconds (Rseconds
nanos (RnanosB�
com.google.protobufBDurationProtoPZ1google.golang.org/protobuf/types/known/durationpb��GPB�Google.Protobuf.WellKnownTypesbproto3
�
proto/feature/variation.protobucketeer.feature"g
	Variation
id (	Rid
value (	Rvalue
name (	Rname 
description (	RdescriptionB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
proto/feature/reason.protobucketeer.feature"�
Reason2
type (2.bucketeer.feature.Reason.TypeRtype
rule_id (	RruleId"Z
Type

TARGET 
RULE
DEFAULT

CLIENT
OFF_VARIATION
PREREQUISITEB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
proto/feature/evaluation.protobucketeer.featureproto/feature/variation.protoproto/feature/reason.proto"�

Evaluation
id (	Rid

feature_id (	R	featureId'
feature_version (RfeatureVersion
user_id (	RuserId!
variation_id (	RvariationId>
	variation (2.bucketeer.feature.VariationBR	variation1
reason (2.bucketeer.feature.ReasonRreason'
variation_value (	RvariationValue%
variation_name	 (	RvariationName"�
UserEvaluations
id (	BRid?
evaluations (2.bucketeer.feature.EvaluationRevaluations

created_at (R	createdAt0
archived_feature_ids (	RarchivedFeatureIds!
force_update (RforceUpdate"*
State

QUEUED 
PARTIAL
FULLB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
proto/user/user.protobucketeer.user"�
User
id (	Rid2
data (2.bucketeer.user.User.DataEntryRdataE
tagged_data (2$.bucketeer.user.User.TaggedDataEntryR
taggedData
	last_seen (RlastSeen

created_at (R	createdAt|
Data:
value (2$.bucketeer.user.User.Data.ValueEntryRvalue8

ValueEntry
key (	Rkey
value (	Rvalue:87
	DataEntry
key (	Rkey
value (	Rvalue:8X
TaggedDataEntry
key (	Rkey/
value (2.bucketeer.user.User.DataRvalue:8B.Z,github.com/bucketeer-io/bucketeer/proto/userbproto3
�4
proto/event/client/event.protobucketeer.event.clientgoogle/protobuf/any.protogoogle/protobuf/duration.protoproto/feature/evaluation.protoproto/feature/reason.protoproto/user/user.proto"p
Event
id (	Rid*
event (2.google.protobuf.AnyRevent%
environment_id (	RenvironmentIdJ"�
EvaluationEvent
	timestamp (R	timestamp

feature_id (	R	featureId'
feature_version (RfeatureVersion
user_id (	RuserId!
variation_id (	RvariationId(
user (2.bucketeer.user.UserRuser1
reason (2.bucketeer.feature.ReasonRreason
tag (	Rtag=
	source_id	 (2 .bucketeer.event.client.SourceIdRsourceId
sdk_version
 (	R
sdkVersionQ
metadata (25.bucketeer.event.client.EvaluationEvent.MetadataEntryRmetadata;
MetadataEntry
key (	Rkey
value (	Rvalue:8"�
	GoalEvent
	timestamp (R	timestamp
goal_id (	RgoalId
user_id (	RuserId
value (Rvalue(
user (2.bucketeer.user.UserRuserC
evaluations (2.bucketeer.feature.EvaluationBRevaluations
tag (	Rtag=
	source_id (2 .bucketeer.event.client.SourceIdRsourceId
sdk_version	 (	R
sdkVersionK
metadata
 (2/.bucketeer.event.client.GoalEvent.MetadataEntryRmetadata;
MetadataEntry
key (	Rkey
value (	Rvalue:8"�
MetricsEvent
	timestamp (R	timestamp*
event (2.google.protobuf.AnyRevent=
	source_id (2 .bucketeer.event.client.SourceIdRsourceId
sdk_version (	R
sdkVersionN
metadata (22.bucketeer.event.client.MetricsEvent.MetadataEntryRmetadata;
MetadataEntry
key (	Rkey
value (	Rvalue:8"�
 GetEvaluationLatencyMetricsEvent\
labels (2D.bucketeer.event.client.GetEvaluationLatencyMetricsEvent.LabelsEntryRlabels5
duration (2.google.protobuf.DurationRduration9
LabelsEntry
key (	Rkey
value (	Rvalue:8:"�
GetEvaluationSizeMetricsEventY
labels (2A.bucketeer.event.client.GetEvaluationSizeMetricsEvent.LabelsEntryRlabels
	size_byte (RsizeByte9
LabelsEntry
key (	Rkey
value (	Rvalue:8:"�
LatencyMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiIdO
labels (27.bucketeer.event.client.LatencyMetricsEvent.LabelsEntryRlabels9
duration (2.google.protobuf.DurationBRduration%
latency_second (RlatencySecond9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
SizeMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiIdL
labels (24.bucketeer.event.client.SizeMetricsEvent.LabelsEntryRlabels
	size_byte (RsizeByte9
LabelsEntry
key (	Rkey
value (	Rvalue:8"5
TimeoutErrorCountMetricsEvent
tag (	Rtag:"6
InternalErrorCountMetricsEvent
tag (	Rtag:"�
 RedirectionRequestExceptionEvent4
api_id (2.bucketeer.event.client.ApiIdRapiId\
labels (2D.bucketeer.event.client.RedirectionRequestExceptionEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
BadRequestErrorMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiIdW
labels (2?.bucketeer.event.client.BadRequestErrorMetricsEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
UnauthorizedErrorMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiIdY
labels (2A.bucketeer.event.client.UnauthorizedErrorMetricsEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
ForbiddenErrorMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiIdV
labels (2>.bucketeer.event.client.ForbiddenErrorMetricsEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
NotFoundErrorMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiIdU
labels (2=.bucketeer.event.client.NotFoundErrorMetricsEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
PayloadTooLargeExceptionEvent4
api_id (2.bucketeer.event.client.ApiIdRapiIdY
labels (2A.bucketeer.event.client.PayloadTooLargeExceptionEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
$ClientClosedRequestErrorMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiId`
labels (2H.bucketeer.event.client.ClientClosedRequestErrorMetricsEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
InternalServerErrorMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiId[
labels (2C.bucketeer.event.client.InternalServerErrorMetricsEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
#ServiceUnavailableErrorMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiId_
labels (2G.bucketeer.event.client.ServiceUnavailableErrorMetricsEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
TimeoutErrorMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiIdT
labels (2<.bucketeer.event.client.TimeoutErrorMetricsEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
InternalErrorMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiIdU
labels (2=.bucketeer.event.client.InternalErrorMetricsEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
NetworkErrorMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiIdT
labels (2<.bucketeer.event.client.NetworkErrorMetricsEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
InternalSdkErrorMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiIdX
labels (2@.bucketeer.event.client.InternalSdkErrorMetricsEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
UnknownErrorMetricsEvent4
api_id (2.bucketeer.event.client.ApiIdRapiIdT
labels (2<.bucketeer.event.client.UnknownErrorMetricsEvent.LabelsEntryRlabels9
LabelsEntry
key (	Rkey
value (	Rvalue:8"�
OpsEvent
	timestamp (R	timestamp

feature_id (	R	featureId'
feature_version (RfeatureVersion!
variation_id (	RvariationId
goal_id (	RgoalId
user_id (	RuserId"\
UserGoalEvent
	timestamp (R	timestamp
goal_id (	RgoalId
value (Rvalue*p
SourceId
UNKNOWN 
ANDROID
IOS
WEB
	GO_SERVER
NODE_SERVER

JAVASCRIPT"*�
ApiId
UNKNOWN_API 
GET_EVALUATION
GET_EVALUATIONS
REGISTER_EVENTS
GET_FEATURE_FLAGS
GET_SEGMENT_USERS
SDK_GET_VARIATIONdB6Z4github.com/bucketeer-io/bucketeer/proto/event/clientbproto3
�
$proto/eventcounter/mau_summary.protobucketeer.eventcounterproto/event/client/event.proto"�

MAUSummary
	yearmonth (	R	yearmonth%
environment_id (	RenvironmentId=
	source_id (2 .bucketeer.event.client.SourceIdRsourceId

user_count (R	userCount#
request_count (RrequestCount)
evaluation_count (RevaluationCount

goal_count (R	goalCount
is_all (RisAll
is_finished	 (R
isFinished

created_at
 (R	createdAt

updated_at (R	updatedAtB6Z4github.com/bucketeer-io/bucketeer/proto/eventcounterbproto3
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
�#
 proto/eventcounter/service.protobucketeer.eventcountergoogle/protobuf/wrappers.proto*proto/eventcounter/experiment_result.proto#proto/eventcounter/timeseries.proto(proto/eventcounter/variation_count.proto"�
#GetExperimentEvaluationCountRequest
start_at (RstartAt
end_at (RendAt

feature_id (	R	featureId'
feature_version (RfeatureVersion#
variation_ids (	RvariationIds%
environment_id (	RenvironmentIdJ"�
$GetExperimentEvaluationCountResponse

feature_id (	R	featureId'
feature_version (RfeatureVersionQ
variation_counts (2&.bucketeer.eventcounter.VariationCountRvariationCounts"�
#GetEvaluationTimeseriesCountRequest

feature_id (	R	featureIdd

time_range (2E.bucketeer.eventcounter.GetEvaluationTimeseriesCountRequest.TimeRangeR	timeRange%
environment_id (	RenvironmentId"c
	TimeRange
UNKNOWN 
TWENTY_FOUR_HOURS

SEVEN_DAYS
FOURTEEN_DAYS
THIRTY_DAYSJ"�
$GetEvaluationTimeseriesCountResponseL
user_counts (2+.bucketeer.eventcounter.VariationTimeseriesR
userCountsN
event_counts (2+.bucketeer.eventcounter.VariationTimeseriesReventCounts"n
GetExperimentResultRequest#
experiment_id (	RexperimentId%
environment_id (	RenvironmentIdJ"t
GetExperimentResultResponseU
experiment_result (2(.bucketeer.eventcounter.ExperimentResultRexperimentResult"�
ListExperimentResultsRequest

feature_id (	R	featureIdD
feature_version (2.google.protobuf.Int32ValueRfeatureVersion%
environment_id (	RenvironmentIdJ"�
ListExperimentResultsResponse\
results (2B.bucketeer.eventcounter.ListExperimentResultsResponse.ResultsEntryRresultsd
ResultsEntry
key (	Rkey>
value (2(.bucketeer.eventcounter.ExperimentResultRvalue:8"�
GetExperimentGoalCountRequest
start_at (RstartAt
end_at (RendAt
goal_id (	RgoalId

feature_id (	R	featureId'
feature_version (RfeatureVersion#
variation_ids (	RvariationIds%
environment_id (	RenvironmentIdJ"�
GetExperimentGoalCountResponse
goal_id (	RgoalIdQ
variation_counts (2&.bucketeer.eventcounter.VariationCountRvariationCounts"�
 GetOpsEvaluationUserCountRequest
ops_rule_id (	R	opsRuleId
	clause_id (	RclauseId

feature_id (	R	featureId'
feature_version (RfeatureVersion!
variation_id (	RvariationId%
environment_id (	RenvironmentIdJ"v
!GetOpsEvaluationUserCountResponse
ops_rule_id (	R	opsRuleId
	clause_id (	RclauseId
count (Rcount"�
GetOpsGoalUserCountRequest
ops_rule_id (	R	opsRuleId
	clause_id (	RclauseId

feature_id (	R	featureId'
feature_version (RfeatureVersion!
variation_id (	RvariationId%
environment_id (	RenvironmentIdJ"p
GetOpsGoalUserCountResponse
ops_rule_id (	R	opsRuleId
	clause_id (	RclauseId
count (Rcount"`
GetMAUCountRequest

year_month (	R	yearMonth%
environment_id (	RenvironmentIdJ"U
GetMAUCountResponse
event_count (R
eventCount

user_count (R	userCount"[
SummarizeMAUCountsRequest

year_month (	R	yearMonth
is_finished (R
isFinished"
SummarizeMAUCountsResponse2�	
EventCounterService�
GetExperimentEvaluationCount;.bucketeer.eventcounter.GetExperimentEvaluationCountRequest<.bucketeer.eventcounter.GetExperimentEvaluationCountResponse" �
GetEvaluationTimeseriesCount;.bucketeer.eventcounter.GetEvaluationTimeseriesCountRequest<.bucketeer.eventcounter.GetEvaluationTimeseriesCountResponse" �
GetExperimentResult2.bucketeer.eventcounter.GetExperimentResultRequest3.bucketeer.eventcounter.GetExperimentResultResponse" �
ListExperimentResults4.bucketeer.eventcounter.ListExperimentResultsRequest5.bucketeer.eventcounter.ListExperimentResultsResponse" �
GetExperimentGoalCount5.bucketeer.eventcounter.GetExperimentGoalCountRequest6.bucketeer.eventcounter.GetExperimentGoalCountResponse" h
GetMAUCount*.bucketeer.eventcounter.GetMAUCountRequest+.bucketeer.eventcounter.GetMAUCountResponse" }
SummarizeMAUCounts1.bucketeer.eventcounter.SummarizeMAUCountsRequest2.bucketeer.eventcounter.SummarizeMAUCountsResponse" �
GetOpsEvaluationUserCount8.bucketeer.eventcounter.GetOpsEvaluationUserCountRequest9.bucketeer.eventcounter.GetOpsEvaluationUserCountResponse" �
GetOpsGoalUserCount2.bucketeer.eventcounter.GetOpsGoalUserCountRequest3.bucketeer.eventcounter.GetOpsGoalUserCountResponse" B6Z4github.com/bucketeer-io/bucketeer/proto/eventcounterbproto3