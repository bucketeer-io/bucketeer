
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
�
proto/feature/variation.protobucketeer.feature"g
	Variation
id (	Rid
value (	Rvalue
name (	Rname 
description (	Rdescription"J
VariationListValue4
values (2.bucketeer.feature.VariationRvaluesB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
proto/feature/reason.protobucketeer.feature"�
Reason2
type (2.bucketeer.feature.Reason.TypeRtype
rule_id (	RruleId"�
Type

TARGET 
RULE
DEFAULT
CLIENT
OFF_VARIATION
PREREQUISITE
ERROR_NO_EVALUATIONS
ERROR_FLAG_NOT_FOUND
ERROR_WRONG_TYPE
ERROR_USER_ID_NOT_SPECIFIED'
#ERROR_FEATURE_FLAG_ID_NOT_SPECIFIED
ERROR_EXCEPTIONB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
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
�
google/api/http.proto
google.api"y
Http*
rules (2.google.api.HttpRuleRrulesE
fully_decode_reserved_expansion (RfullyDecodeReservedExpansion"�
HttpRule
selector (	Rselector
get (	H Rget
put (	H Rput
post (	H Rpost
delete (	H Rdelete
patch (	H Rpatch7
custom (2.google.api.CustomHttpPatternH Rcustom
body (	RbodyE
additional_bindings (2.google.api.HttpRuleRadditionalBindingsB	
pattern";
CustomHttpPattern
kind (	Rkind
path (	RpathBj
com.google.apiB	HttpProtoPZAgoogle.golang.org/genproto/googleapis/api/annotations;annotations��GAPIbproto3
�F
 google/protobuf/descriptor.protogoogle.protobuf"M
FileDescriptorSet8
file (2$.google.protobuf.FileDescriptorProtoRfile"�
FileDescriptorProto
name (	Rname
package (	Rpackage

dependency (	R
dependency+
public_dependency
 (RpublicDependency'
weak_dependency (RweakDependencyC
message_type (2 .google.protobuf.DescriptorProtoRmessageTypeA
	enum_type (2$.google.protobuf.EnumDescriptorProtoRenumTypeA
service (2'.google.protobuf.ServiceDescriptorProtoRserviceC
	extension (2%.google.protobuf.FieldDescriptorProtoR	extension6
options (2.google.protobuf.FileOptionsRoptionsI
source_code_info	 (2.google.protobuf.SourceCodeInfoRsourceCodeInfo
syntax (	Rsyntax
edition (	Redition"�
DescriptorProto
name (	Rname;
field (2%.google.protobuf.FieldDescriptorProtoRfieldC
	extension (2%.google.protobuf.FieldDescriptorProtoR	extensionA
nested_type (2 .google.protobuf.DescriptorProtoR
nestedTypeA
	enum_type (2$.google.protobuf.EnumDescriptorProtoRenumTypeX
extension_range (2/.google.protobuf.DescriptorProto.ExtensionRangeRextensionRangeD

oneof_decl (2%.google.protobuf.OneofDescriptorProtoR	oneofDecl9
options (2.google.protobuf.MessageOptionsRoptionsU
reserved_range	 (2..google.protobuf.DescriptorProto.ReservedRangeRreservedRange#
reserved_name
 (	RreservedNamez
ExtensionRange
start (Rstart
end (Rend@
options (2&.google.protobuf.ExtensionRangeOptionsRoptions7
ReservedRange
start (Rstart
end (Rend"�
ExtensionRangeOptionsX
uninterpreted_option� (2$.google.protobuf.UninterpretedOptionRuninterpretedOptionY
declaration (22.google.protobuf.ExtensionRangeOptions.DeclarationB�Rdeclarationh
verification (28.google.protobuf.ExtensionRangeOptions.VerificationState:
UNVERIFIEDRverification�
Declaration
number (Rnumber
	full_name (	RfullName
type (	Rtype#
is_repeated (BR
isRepeated
reserved (Rreserved
repeated (Rrepeated"4
VerificationState
DECLARATION 

UNVERIFIED*	�����"�
FieldDescriptorProto
name (	Rname
number (RnumberA
label (2+.google.protobuf.FieldDescriptorProto.LabelRlabel>
type (2*.google.protobuf.FieldDescriptorProto.TypeRtype
	type_name (	RtypeName
extendee (	Rextendee#
default_value (	RdefaultValue
oneof_index	 (R
oneofIndex
	json_name
 (	RjsonName7
options (2.google.protobuf.FieldOptionsRoptions'
proto3_optional (Rproto3Optional"�
Type
TYPE_DOUBLE

TYPE_FLOAT

TYPE_INT64
TYPE_UINT64

TYPE_INT32
TYPE_FIXED64
TYPE_FIXED32
	TYPE_BOOL
TYPE_STRING	

TYPE_GROUP

TYPE_MESSAGE

TYPE_BYTES
TYPE_UINT32
	TYPE_ENUM
TYPE_SFIXED32
TYPE_SFIXED64
TYPE_SINT32
TYPE_SINT64"C
Label
LABEL_OPTIONAL
LABEL_REQUIRED
LABEL_REPEATED"c
OneofDescriptorProto
name (	Rname7
options (2.google.protobuf.OneofOptionsRoptions"�
EnumDescriptorProto
name (	Rname?
value (2).google.protobuf.EnumValueDescriptorProtoRvalue6
options (2.google.protobuf.EnumOptionsRoptions]
reserved_range (26.google.protobuf.EnumDescriptorProto.EnumReservedRangeRreservedRange#
reserved_name (	RreservedName;
EnumReservedRange
start (Rstart
end (Rend"�
EnumValueDescriptorProto
name (	Rname
number (Rnumber;
options (2!.google.protobuf.EnumValueOptionsRoptions"�
ServiceDescriptorProto
name (	Rname>
method (2&.google.protobuf.MethodDescriptorProtoRmethod9
options (2.google.protobuf.ServiceOptionsRoptions"�
MethodDescriptorProto
name (	Rname

input_type (	R	inputType
output_type (	R
outputType8
options (2.google.protobuf.MethodOptionsRoptions0
client_streaming (:falseRclientStreaming0
server_streaming (:falseRserverStreaming"�	
FileOptions!
java_package (	RjavaPackage0
java_outer_classname (	RjavaOuterClassname5
java_multiple_files
 (:falseRjavaMultipleFilesD
java_generate_equals_and_hash (BRjavaGenerateEqualsAndHash:
java_string_check_utf8 (:falseRjavaStringCheckUtf8S
optimize_for	 (2).google.protobuf.FileOptions.OptimizeMode:SPEEDRoptimizeFor

go_package (	R	goPackage5
cc_generic_services (:falseRccGenericServices9
java_generic_services (:falseRjavaGenericServices5
py_generic_services (:falseRpyGenericServices7
php_generic_services* (:falseRphpGenericServices%

deprecated (:falseR
deprecated.
cc_enable_arenas (:trueRccEnableArenas*
objc_class_prefix$ (	RobjcClassPrefix)
csharp_namespace% (	RcsharpNamespace!
swift_prefix' (	RswiftPrefix(
php_class_prefix( (	RphpClassPrefix#
php_namespace) (	RphpNamespace4
php_metadata_namespace, (	RphpMetadataNamespace!
ruby_package- (	RrubyPackageX
uninterpreted_option� (2$.google.protobuf.UninterpretedOptionRuninterpretedOption":
OptimizeMode	
SPEED
	CODE_SIZE
LITE_RUNTIME*	�����J&'"�
MessageOptions<
message_set_wire_format (:falseRmessageSetWireFormatL
no_standard_descriptor_accessor (:falseRnoStandardDescriptorAccessor%

deprecated (:falseR
deprecated
	map_entry (RmapEntryV
&deprecated_legacy_json_field_conflicts (BR"deprecatedLegacyJsonFieldConflictsX
uninterpreted_option� (2$.google.protobuf.UninterpretedOptionRuninterpretedOption*	�����JJJJ	J	
"�	
FieldOptionsA
ctype (2#.google.protobuf.FieldOptions.CType:STRINGRctype
packed (RpackedG
jstype (2$.google.protobuf.FieldOptions.JSType:	JS_NORMALRjstype
lazy (:falseRlazy.
unverified_lazy (:falseRunverifiedLazy%

deprecated (:falseR
deprecated
weak
 (:falseRweak(
debug_redact (:falseRdebugRedactK
	retention (2-.google.protobuf.FieldOptions.OptionRetentionR	retentionJ
target (2..google.protobuf.FieldOptions.OptionTargetTypeBRtargetH
targets (2..google.protobuf.FieldOptions.OptionTargetTypeRtargetsX
uninterpreted_option� (2$.google.protobuf.UninterpretedOptionRuninterpretedOption"/
CType

STRING 
CORD
STRING_PIECE"5
JSType
	JS_NORMAL 
	JS_STRING
	JS_NUMBER"U
OptionRetention
RETENTION_UNKNOWN 
RETENTION_RUNTIME
RETENTION_SOURCE"�
OptionTargetType
TARGET_TYPE_UNKNOWN 
TARGET_TYPE_FILE
TARGET_TYPE_EXTENSION_RANGE
TARGET_TYPE_MESSAGE
TARGET_TYPE_FIELD
TARGET_TYPE_ONEOF
TARGET_TYPE_ENUM
TARGET_TYPE_ENUM_ENTRY
TARGET_TYPE_SERVICE
TARGET_TYPE_METHOD	*	�����J"s
OneofOptionsX
uninterpreted_option� (2$.google.protobuf.UninterpretedOptionRuninterpretedOption*	�����"�
EnumOptions
allow_alias (R
allowAlias%

deprecated (:falseR
deprecatedV
&deprecated_legacy_json_field_conflicts (BR"deprecatedLegacyJsonFieldConflictsX
uninterpreted_option� (2$.google.protobuf.UninterpretedOptionRuninterpretedOption*	�����J"�
EnumValueOptions%

deprecated (:falseR
deprecatedX
uninterpreted_option� (2$.google.protobuf.UninterpretedOptionRuninterpretedOption*	�����"�
ServiceOptions%

deprecated! (:falseR
deprecatedX
uninterpreted_option� (2$.google.protobuf.UninterpretedOptionRuninterpretedOption*	�����"�
MethodOptions%

deprecated! (:falseR
deprecatedq
idempotency_level" (2/.google.protobuf.MethodOptions.IdempotencyLevel:IDEMPOTENCY_UNKNOWNRidempotencyLevelX
uninterpreted_option� (2$.google.protobuf.UninterpretedOptionRuninterpretedOption"P
IdempotencyLevel
IDEMPOTENCY_UNKNOWN 
NO_SIDE_EFFECTS

IDEMPOTENT*	�����"�
UninterpretedOptionA
name (2-.google.protobuf.UninterpretedOption.NamePartRname)
identifier_value (	RidentifierValue,
positive_int_value (RpositiveIntValue,
negative_int_value (RnegativeIntValue!
double_value (RdoubleValue!
string_value (RstringValue'
aggregate_value (	RaggregateValueJ
NamePart
	name_part (	RnamePart!
is_extension (RisExtension"�
SourceCodeInfoD
location (2(.google.protobuf.SourceCodeInfo.LocationRlocation�
Location
path (BRpath
span (BRspan)
leading_comments (	RleadingComments+
trailing_comments (	RtrailingComments:
leading_detached_comments (	RleadingDetachedComments"�
GeneratedCodeInfoM

annotation (2-.google.protobuf.GeneratedCodeInfo.AnnotationR
annotation�

Annotation
path (BRpath
source_file (	R
sourceFile
begin (Rbegin
end (RendR
semantic (26.google.protobuf.GeneratedCodeInfo.Annotation.SemanticRsemantic"(
Semantic
NONE 
SET	
ALIASB~
com.google.protobufBDescriptorProtosHZ-google.golang.org/protobuf/types/descriptorpb��GPB�Google.Protobuf.Reflection
�
google/api/annotations.proto
google.apigoogle/api/http.proto google/protobuf/descriptor.proto:K
http.google.protobuf.MethodOptions�ʼ" (2.google.api.HttpRuleRhttpBn
com.google.apiBAnnotationsProtoPZAgoogle.golang.org/genproto/googleapis/api/annotations;annotations�GAPIbproto3
�
google/api/field_behavior.proto
google.api google/protobuf/descriptor.proto*�
FieldBehavior
FIELD_BEHAVIOR_UNSPECIFIED 
OPTIONAL
REQUIRED
OUTPUT_ONLY

INPUT_ONLY
	IMMUTABLE
UNORDERED_LIST
NON_EMPTY_DEFAULT

IDENTIFIER:d
field_behavior.google.protobuf.FieldOptions� (2.google.api.FieldBehaviorB RfieldBehaviorBp
com.google.apiBFieldBehaviorProtoPZAgoogle.golang.org/genproto/googleapis/api/annotations;annotations�GAPIbproto3
�
google/protobuf/struct.protogoogle.protobuf"�
Struct;
fields (2#.google.protobuf.Struct.FieldsEntryRfieldsQ
FieldsEntry
key (	Rkey,
value (2.google.protobuf.ValueRvalue:8"�
Value;

null_value (2.google.protobuf.NullValueH R	nullValue#
number_value (H RnumberValue#
string_value (	H RstringValue

bool_value (H R	boolValue<
struct_value (2.google.protobuf.StructH RstructValue;

list_value (2.google.protobuf.ListValueH R	listValueB
kind";
	ListValue.
values (2.google.protobuf.ValueRvalues*
	NullValue

NULL_VALUE B
com.google.protobufBStructProtoPZ/google.golang.org/protobuf/types/known/structpb��GPB�Google.Protobuf.WellKnownTypesbproto3
�>
,protoc-gen-openapiv2/options/openapiv2.proto)grpc.gateway.protoc_gen_openapiv2.optionsgoogle/protobuf/struct.proto"�
Swagger
swagger (	RswaggerC
info (2/.grpc.gateway.protoc_gen_openapiv2.options.InfoRinfo
host (	Rhost
	base_path (	RbasePathK
schemes (21.grpc.gateway.protoc_gen_openapiv2.options.SchemeRschemes
consumes (	Rconsumes
produces (	Rproduces_
	responses
 (2A.grpc.gateway.protoc_gen_openapiv2.options.Swagger.ResponsesEntryR	responsesq
security_definitions (2>.grpc.gateway.protoc_gen_openapiv2.options.SecurityDefinitionsRsecurityDefinitionsZ
security (2>.grpc.gateway.protoc_gen_openapiv2.options.SecurityRequirementRsecurityB
tags (2..grpc.gateway.protoc_gen_openapiv2.options.TagRtagse
external_docs (2@.grpc.gateway.protoc_gen_openapiv2.options.ExternalDocumentationRexternalDocsb

extensions (2B.grpc.gateway.protoc_gen_openapiv2.options.Swagger.ExtensionsEntryR
extensionsq
ResponsesEntry
key (	RkeyI
value (23.grpc.gateway.protoc_gen_openapiv2.options.ResponseRvalue:8U
ExtensionsEntry
key (	Rkey,
value (2.google.protobuf.ValueRvalue:8J	J	
"�
	Operation
tags (	Rtags
summary (	Rsummary 
description (	Rdescriptione
external_docs (2@.grpc.gateway.protoc_gen_openapiv2.options.ExternalDocumentationRexternalDocs!
operation_id (	RoperationId
consumes (	Rconsumes
produces (	Rproducesa
	responses	 (2C.grpc.gateway.protoc_gen_openapiv2.options.Operation.ResponsesEntryR	responsesK
schemes
 (21.grpc.gateway.protoc_gen_openapiv2.options.SchemeRschemes

deprecated (R
deprecatedZ
security (2>.grpc.gateway.protoc_gen_openapiv2.options.SecurityRequirementRsecurityd

extensions (2D.grpc.gateway.protoc_gen_openapiv2.options.Operation.ExtensionsEntryR
extensionsU

parameters (25.grpc.gateway.protoc_gen_openapiv2.options.ParametersR
parametersq
ResponsesEntry
key (	RkeyI
value (23.grpc.gateway.protoc_gen_openapiv2.options.ResponseRvalue:8U
ExtensionsEntry
key (	Rkey,
value (2.google.protobuf.ValueRvalue:8J	"b

ParametersT
headers (2:.grpc.gateway.protoc_gen_openapiv2.options.HeaderParameterRheaders"�
HeaderParameter
name (	Rname 
description (	RdescriptionS
type (2?.grpc.gateway.protoc_gen_openapiv2.options.HeaderParameter.TypeRtype
format (	Rformat
required (Rrequired"E
Type
UNKNOWN 

STRING

NUMBER
INTEGER
BOOLEANJJ"�
Header 
description (	Rdescription
type (	Rtype
format (	Rformat
default (	Rdefault
pattern (	RpatternJJJJ	J	
J
JJJJJJJ"�
Response 
description (	RdescriptionI
schema (21.grpc.gateway.protoc_gen_openapiv2.options.SchemaRschemaZ
headers (2@.grpc.gateway.protoc_gen_openapiv2.options.Response.HeadersEntryRheaders]
examples (2A.grpc.gateway.protoc_gen_openapiv2.options.Response.ExamplesEntryRexamplesc

extensions (2C.grpc.gateway.protoc_gen_openapiv2.options.Response.ExtensionsEntryR
extensionsm
HeadersEntry
key (	RkeyG
value (21.grpc.gateway.protoc_gen_openapiv2.options.HeaderRvalue:8;
ExamplesEntry
key (	Rkey
value (	Rvalue:8U
ExtensionsEntry
key (	Rkey,
value (2.google.protobuf.ValueRvalue:8"�
Info
title (	Rtitle 
description (	Rdescription(
terms_of_service (	RtermsOfServiceL
contact (22.grpc.gateway.protoc_gen_openapiv2.options.ContactRcontactL
license (22.grpc.gateway.protoc_gen_openapiv2.options.LicenseRlicense
version (	Rversion_

extensions (2?.grpc.gateway.protoc_gen_openapiv2.options.Info.ExtensionsEntryR
extensionsU
ExtensionsEntry
key (	Rkey,
value (2.google.protobuf.ValueRvalue:8"E
Contact
name (	Rname
url (	Rurl
email (	Remail"/
License
name (	Rname
url (	Rurl"K
ExternalDocumentation 
description (	Rdescription
url (	Rurl"�
SchemaV
json_schema (25.grpc.gateway.protoc_gen_openapiv2.options.JSONSchemaR
jsonSchema$
discriminator (	Rdiscriminator
	read_only (RreadOnlye
external_docs (2@.grpc.gateway.protoc_gen_openapiv2.options.ExternalDocumentationRexternalDocs
example (	RexampleJ"�


JSONSchema
ref (	Rref
title (	Rtitle 
description (	Rdescription
default (	Rdefault
	read_only (RreadOnly
example	 (	Rexample
multiple_of
 (R
multipleOf
maximum (Rmaximum+
exclusive_maximum (RexclusiveMaximum
minimum (Rminimum+
exclusive_minimum (RexclusiveMinimum

max_length (R	maxLength

min_length (R	minLength
pattern (	Rpattern
	max_items (RmaxItems
	min_items (RminItems!
unique_items (RuniqueItems%
max_properties (RmaxProperties%
min_properties (RminProperties
required (	Rrequired
array" (	Rarray_
type# (2K.grpc.gateway.protoc_gen_openapiv2.options.JSONSchema.JSONSchemaSimpleTypesRtype
format$ (	Rformat
enum. (	Renumz
field_configuration� (2H.grpc.gateway.protoc_gen_openapiv2.options.JSONSchema.FieldConfigurationRfieldConfiguratione

extensions0 (2E.grpc.gateway.protoc_gen_openapiv2.options.JSONSchema.ExtensionsEntryR
extensions<
FieldConfiguration&
path_param_name/ (	RpathParamNameU
ExtensionsEntry
key (	Rkey,
value (2.google.protobuf.ValueRvalue:8"w
JSONSchemaSimpleTypes
UNKNOWN 	
ARRAY
BOOLEAN
INTEGER
NULL

NUMBER

OBJECT

STRINGJJJJJJJJJJ"J%*J*+J+."�
Tag
name (	Rname 
description (	Rdescriptione
external_docs (2@.grpc.gateway.protoc_gen_openapiv2.options.ExternalDocumentationRexternalDocs^

extensions (2>.grpc.gateway.protoc_gen_openapiv2.options.Tag.ExtensionsEntryR
extensionsU
ExtensionsEntry
key (	Rkey,
value (2.google.protobuf.ValueRvalue:8"�
SecurityDefinitionsh
security (2L.grpc.gateway.protoc_gen_openapiv2.options.SecurityDefinitions.SecurityEntryRsecurityv
SecurityEntry
key (	RkeyO
value (29.grpc.gateway.protoc_gen_openapiv2.options.SecuritySchemeRvalue:8"�
SecuritySchemeR
type (2>.grpc.gateway.protoc_gen_openapiv2.options.SecurityScheme.TypeRtype 
description (	Rdescription
name (	RnameL
in (2<.grpc.gateway.protoc_gen_openapiv2.options.SecurityScheme.InRinR
flow (2>.grpc.gateway.protoc_gen_openapiv2.options.SecurityScheme.FlowRflow+
authorization_url (	RauthorizationUrl
	token_url (	RtokenUrlI
scopes (21.grpc.gateway.protoc_gen_openapiv2.options.ScopesRscopesi

extensions	 (2I.grpc.gateway.protoc_gen_openapiv2.options.SecurityScheme.ExtensionsEntryR
extensionsU
ExtensionsEntry
key (	Rkey,
value (2.google.protobuf.ValueRvalue:8"K
Type
TYPE_INVALID 

TYPE_BASIC
TYPE_API_KEY
TYPE_OAUTH2"1
In

IN_INVALID 
IN_QUERY
	IN_HEADER"j
Flow
FLOW_INVALID 
FLOW_IMPLICIT
FLOW_PASSWORD
FLOW_APPLICATION
FLOW_ACCESS_CODE"�
SecurityRequirement�
security_requirement (2W.grpc.gateway.protoc_gen_openapiv2.options.SecurityRequirement.SecurityRequirementEntryRsecurityRequirement0
SecurityRequirementValue
scope (	Rscope�
SecurityRequirementEntry
key (	Rkeym
value (2W.grpc.gateway.protoc_gen_openapiv2.options.SecurityRequirement.SecurityRequirementValueRvalue:8"�
ScopesR
scope (2<.grpc.gateway.protoc_gen_openapiv2.options.Scopes.ScopeEntryRscope8

ScopeEntry
key (	Rkey
value (	Rvalue:8*;
Scheme
UNKNOWN 
HTTP	
HTTPS
WS
WSSBHZFgithub.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/optionsbproto3
�
.protoc-gen-openapiv2/options/annotations.proto)grpc.gateway.protoc_gen_openapiv2.options google/protobuf/descriptor.proto,protoc-gen-openapiv2/options/openapiv2.proto:~
openapiv2_swagger.google.protobuf.FileOptions� (22.grpc.gateway.protoc_gen_openapiv2.options.SwaggerRopenapiv2Swagger:�
openapiv2_operation.google.protobuf.MethodOptions� (24.grpc.gateway.protoc_gen_openapiv2.options.OperationRopenapiv2Operation:~
openapiv2_schema.google.protobuf.MessageOptions� (21.grpc.gateway.protoc_gen_openapiv2.options.SchemaRopenapiv2Schema:u
openapiv2_tag.google.protobuf.ServiceOptions� (2..grpc.gateway.protoc_gen_openapiv2.options.TagRopenapiv2Tag:~
openapiv2_field.google.protobuf.FieldOptions� (25.grpc.gateway.protoc_gen_openapiv2.options.JSONSchemaRopenapiv2FieldBHZFgithub.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/optionsbproto3
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
�V
 proto/eventcounter/service.protobucketeer.eventcountergoogle/api/annotations.protogoogle/api/field_behavior.proto.protoc-gen-openapiv2/options/annotations.protogoogle/protobuf/wrappers.proto*proto/eventcounter/experiment_result.proto#proto/eventcounter/timeseries.proto(proto/eventcounter/variation_count.proto"�
#GetExperimentEvaluationCountRequest
start_at (B�ARstartAt
end_at (B�ARendAt"

feature_id (	B�AR	featureId,
feature_version (B�ARfeatureVersion#
variation_ids (	RvariationIds*
environment_id (	B�ARenvironmentIdJ"�
$GetExperimentEvaluationCountResponse

feature_id (	R	featureId'
feature_version (RfeatureVersionQ
variation_counts (2&.bucketeer.eventcounter.VariationCountRvariationCounts"�
#GetEvaluationTimeseriesCountRequest"

feature_id (	B�AR	featureIdi

time_range (2E.bucketeer.eventcounter.GetEvaluationTimeseriesCountRequest.TimeRangeB�AR	timeRange*
environment_id (	B�ARenvironmentId"c
	TimeRange
UNKNOWN 
TWENTY_FOUR_HOURS

SEVEN_DAYS
FOURTEEN_DAYS
THIRTY_DAYSJ"�
$GetEvaluationTimeseriesCountResponseL
user_counts (2+.bucketeer.eventcounter.VariationTimeseriesR
userCountsN
event_counts (2+.bucketeer.eventcounter.VariationTimeseriesReventCounts"x
GetExperimentResultRequest(
experiment_id (	B�ARexperimentId*
environment_id (	B�ARenvironmentIdJ"t
GetExperimentResultResponseU
experiment_result (2(.bucketeer.eventcounter.ExperimentResultRexperimentResult"�
ListExperimentResultsRequest"

feature_id (	B�AR	featureIdD
feature_version (2.google.protobuf.Int32ValueRfeatureVersion*
environment_id (	B�ARenvironmentIdJ"�
ListExperimentResultsResponse\
results (2B.bucketeer.eventcounter.ListExperimentResultsResponse.ResultsEntryRresultsd
ResultsEntry
key (	Rkey>
value (2(.bucketeer.eventcounter.ExperimentResultRvalue:8"�
GetExperimentGoalCountRequest
start_at (B�ARstartAt
end_at (B�ARendAt
goal_id (	B�ARgoalId"

feature_id (	B�AR	featureId,
feature_version (B�ARfeatureVersion#
variation_ids (	RvariationIds*
environment_id (	B�ARenvironmentIdJ"�
GetExperimentGoalCountResponse
goal_id (	RgoalIdQ
variation_counts (2&.bucketeer.eventcounter.VariationCountRvariationCounts"�
 GetOpsEvaluationUserCountRequest#
ops_rule_id (	B�AR	opsRuleId 
	clause_id (	B�ARclauseId"

feature_id (	B�AR	featureId,
feature_version (B�ARfeatureVersion&
variation_id (	B�ARvariationId*
environment_id (	B�ARenvironmentIdJ"v
!GetOpsEvaluationUserCountResponse
ops_rule_id (	R	opsRuleId
	clause_id (	RclauseId
count (Rcount"�
GetOpsGoalUserCountRequest#
ops_rule_id (	B�AR	opsRuleId 
	clause_id (	B�ARclauseId"

feature_id (	B�AR	featureId,
feature_version (B�ARfeatureVersion&
variation_id (	B�ARvariationId*
environment_id (	B�ARenvironmentIdJ"p
GetOpsGoalUserCountResponse
ops_rule_id (	R	opsRuleId
	clause_id (	RclauseId
count (Rcount"j
GetMAUCountRequest"

year_month (	B�AR	yearMonth*
environment_id (	B�ARenvironmentIdJ"U
GetMAUCountResponse
event_count (R
eventCount

user_count (R	userCount"`
SummarizeMAUCountsRequest"

year_month (	B�AR	yearMonth
is_finished (R
isFinished"
SummarizeMAUCountsResponse2�9
EventCounterService�
GetExperimentEvaluationCount;.bucketeer.eventcounter.GetExperimentEvaluationCountRequest<.bucketeer.eventcounter.GetExperimentEvaluationCountResponse"��A�
experiment_evaluation_countGet Experiment Evaluation Count"Get an experiment evaluation count*4web.v1.event_counter.experiment_evaluation_count.getJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
404�
2Returned when the requested resource is not found.
.google.rpc.Status"H
application/json4{ "code": 5, "message": "not found", "details": [] }���!/v1/experiment_evaluation_count�
GetEvaluationTimeseriesCount;.bucketeer.eventcounter.GetEvaluationTimeseriesCountRequest<.bucketeer.eventcounter.GetEvaluationTimeseriesCountResponse"��A�
evaluation_timeseries_countGet Evaluation Timeseries Count"Get an evaluation timeseries count*4web.v1.event_counter.evaluation_timeseries_count.getJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
404�
2Returned when the requested resource is not found.
.google.rpc.Status"H
application/json4{ "code": 5, "message": "not found", "details": [] }���!/v1/evaluation_timeseries_count�
GetExperimentResult2.bucketeer.eventcounter.GetExperimentResultRequest3.bucketeer.eventcounter.GetExperimentResultResponse"��A�
experiment_resultGet Experiment ResultGet an experiment result**web.v1.event_counter.experiment_result.getJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
404�
2Returned when the requested resource is not found.
.google.rpc.Status"H
application/json4{ "code": 5, "message": "not found", "details": [] }���/v1/experiment_result�
ListExperimentResults4.bucketeer.eventcounter.ListExperimentResultsRequest5.bucketeer.eventcounter.ListExperimentResultsResponse"��A�
experiment_resultList Experiment ResultsList experiment results*+web.v1.event_counter.experiment_result.listJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
404�
2Returned when the requested resource is not found.
.google.rpc.Status"H
application/json4{ "code": 5, "message": "not found", "details": [] }���/v1/experiment_results�
GetExperimentGoalCount5.bucketeer.eventcounter.GetExperimentGoalCountRequest6.bucketeer.eventcounter.GetExperimentGoalCountResponse"��A�
experiment_goal_countGet Experiment Goal CountGet an experiment goal count*.web.v1.event_counter.experiment_goal_count.getJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
404�
2Returned when the requested resource is not found.
.google.rpc.Status"H
application/json4{ "code": 5, "message": "not found", "details": [] }���/v1/experiment_goal_count�
GetMAUCount*.bucketeer.eventcounter.GetMAUCountRequest+.bucketeer.eventcounter.GetMAUCountResponse"��A�
	mau_countGet MAU CountGet MAU count*"web.v1.event_counter.mau_count.getJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
404�
2Returned when the requested resource is not found.
.google.rpc.Status"H
application/json4{ "code": 5, "message": "not found", "details": [] }���/v1/mau_count�
SummarizeMAUCounts1.bucketeer.eventcounter.SummarizeMAUCountsRequest2.bucketeer.eventcounter.SummarizeMAUCountsResponse"��A�
	mau_countSummarize MAU CountsSummarize MAU counts*-web.v1.event_counter.summarize_mau_counts.getJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
404�
2Returned when the requested resource is not found.
.google.rpc.Status"H
application/json4{ "code": 5, "message": "not found", "details": [] }���/v1/summarize_mau_counts�
GetOpsEvaluationUserCount8.bucketeer.eventcounter.GetOpsEvaluationUserCountRequest9.bucketeer.eventcounter.GetOpsEvaluationUserCountResponse"��A�
ops_evaluation_user_countGet Ops Evaluation User Count Get an ops evaluation user count*2web.v1.event_counter.ops_evaluation_user_count.getJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
404�
2Returned when the requested resource is not found.
.google.rpc.Status"H
application/json4{ "code": 5, "message": "not found", "details": [] }���/v1/ops_evaluation_user_count�
GetOpsGoalUserCount2.bucketeer.eventcounter.GetOpsGoalUserCountRequest3.bucketeer.eventcounter.GetOpsGoalUserCountResponse"��A�
ops_goal_user_countGet Ops Goal User CountGet an ops goal user count*,web.v1.event_counter.ops_goal_user_count.getJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
404�
2Returned when the requested resource is not found.
.google.rpc.Status"H
application/json4{ "code": 5, "message": "not found", "details": [] }���/v1/ops_goal_user_countB6Z4github.com/bucketeer-io/bucketeer/proto/eventcounterbproto3