
�
proto/feature/clause.protobucketeer.feature"�
Clause
id (	Rid
	attribute (	R	attribute>
operator (2".bucketeer.feature.Clause.OperatorRoperator
values (	Rvalues"�
Operator

EQUALS 
IN
	ENDS_WITH
STARTS_WITH
SEGMENT
GREATER
GREATER_OR_EQUAL
LESS
LESS_OR_EQUAL

BEFORE		
AFTER

FEATURE_FLAG
PARTIALLY_MATCHB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
google/protobuf/any.protogoogle.protobuf"6
Any
type_url (	RtypeUrl
value (RvalueBv
com.google.protobufBAnyProtoPZ,google.golang.org/protobuf/types/known/anypb�GPB�Google.Protobuf.WellKnownTypesbproto3
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
�
proto/feature/strategy.protobucketeer.feature"-
FixedStrategy
	variation (	R	variation"�
RolloutStrategyL

variations (2,.bucketeer.feature.RolloutStrategy.VariationR
variationsA
	Variation
	variation (	R	variation
weight (Rweight"�
Strategy4
type (2 .bucketeer.feature.Strategy.TypeRtypeG
fixed_strategy (2 .bucketeer.feature.FixedStrategyRfixedStrategyM
rollout_strategy (2".bucketeer.feature.RolloutStrategyRrolloutStrategy"
Type	
FIXED 
ROLLOUTB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
proto/feature/rule.protobucketeer.featureproto/feature/clause.protoproto/feature/strategy.proto"�
Rule
id (	Rid7
strategy (2.bucketeer.feature.StrategyRstrategy3
clauses (2.bucketeer.feature.ClauseRclausesB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
proto/feature/target.protobucketeer.feature"<
Target
	variation (	R	variation
users (	RusersB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
proto/feature/variation.protobucketeer.feature"g
	Variation
id (	Rid
value (	Rvalue
name (	Rname 
description (	RdescriptionB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
*proto/feature/feature_last_used_info.protobucketeer.feature"�
FeatureLastUsedInfo

feature_id (	R	featureId
version (Rversion 
last_used_at (R
lastUsedAt

created_at (R	createdAt2
client_oldest_version (	RclientOldestVersion2
client_latest_version (	RclientLatestVersionB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
 proto/feature/prerequisite.protobucketeer.feature"P
Prerequisite

feature_id (	R	featureId!
variation_id (	RvariationIdB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
proto/feature/feature.protobucketeer.featureproto/feature/rule.protoproto/feature/target.protoproto/feature/variation.protoproto/feature/strategy.proto*proto/feature/feature_last_used_info.proto proto/feature/prerequisite.proto"�
Feature
id (	Rid
name (	Rname 
description (	Rdescription
enabled (Renabled
deleted (Rdeleted9
evaluation_undelayable (BRevaluationUndelayable
ttl (Rttl
version (Rversion

created_at	 (R	createdAt

updated_at
 (R	updatedAt<

variations (2.bucketeer.feature.VariationR
variations3
targets (2.bucketeer.feature.TargetRtargets-
rules (2.bucketeer.feature.RuleRrulesF
default_strategy (2.bucketeer.feature.StrategyRdefaultStrategy#
off_variation (	RoffVariation
tags (	RtagsL
last_used_info (2&.bucketeer.feature.FeatureLastUsedInfoRlastUsedInfo

maintainer (	R
maintainerO
variation_type (2(.bucketeer.feature.Feature.VariationTypeRvariationType
archived (RarchivedE
prerequisites (2.bucketeer.feature.PrerequisiteRprerequisites#
sampling_seed (	RsamplingSeed">
VariationType

STRING 
BOOLEAN

NUMBER
JSON"R
Features6
features (2.bucketeer.feature.FeatureRfeatures
id (	Rid"S
Tag
id (	Rid

created_at (R	createdAt

updated_at (R	updatedAtB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
proto/feature/segment.protobucketeer.featureproto/feature/rule.protoproto/feature/feature.proto"�
Segment
id (	Rid
name (	Rname 
description (	Rdescription-
rules (2.bucketeer.feature.RuleRrules

created_at (R	createdAt

updated_at (R	updatedAt
version (BRversion
deleted (Rdeleted.
included_user_count	 (RincludedUserCount2
excluded_user_count
 (BRexcludedUserCount9
status (2!.bucketeer.feature.Segment.StatusRstatus'
is_in_use_status (RisInUseStatus6
features (2.bucketeer.feature.FeatureRfeatures">
Status
INITIAL 
	UPLOADING
SUCEEDED

FAILED"�
SegmentUser
id (	Rid

segment_id (	R	segmentId
user_id (	RuserId:
state (2$.bucketeer.feature.SegmentUser.StateRstate
deleted (Rdeleted"'
State
INCLUDED 
EXCLUDED"�
SegmentUsers

segment_id (	R	segmentId4
users (2.bucketeer.feature.SegmentUserRusers

updated_at (R	updatedAtB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�
 proto/feature/flag_trigger.protobucketeer.feature"�
FlagTrigger
id (	Rid

feature_id (	R	featureId7
type (2#.bucketeer.feature.FlagTrigger.TypeRtype=
action (2%.bucketeer.feature.FlagTrigger.ActionRaction 
description (	Rdescription#
trigger_count (RtriggerCount*
last_triggered_at (RlastTriggeredAt
token	 (	Rtoken
disabled
 (Rdisabled

created_at (R	createdAt

updated_at (R	updatedAt%
environment_id (	RenvironmentId"*
Type
Type_UNKNOWN 
Type_WEBHOOK";
Action
Action_UNKNOWN 
	Action_ON

Action_OFFJB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
�#
proto/feature/command.protobucketeer.featuregoogle/protobuf/any.protogoogle/protobuf/wrappers.protoproto/feature/clause.protoproto/feature/feature.protoproto/feature/rule.protoproto/feature/variation.protoproto/feature/strategy.protoproto/feature/segment.proto proto/feature/prerequisite.proto proto/feature/flag_trigger.proto"9
Command.
command (2.google.protobuf.AnyRcommand"�
CreateFeatureCommand
id (	Rid
name (	Rname 
description (	Rdescription<

variations (2.bucketeer.feature.VariationR
variations
tags (	RtagsX
default_on_variation_index (2.google.protobuf.Int32ValueRdefaultOnVariationIndexZ
default_off_variation_index (2.google.protobuf.Int32ValueRdefaultOffVariationIndexO
variation_type (2(.bucketeer.feature.Feature.VariationTypeRvariationType"
ArchiveFeatureCommand"
UnarchiveFeatureCommand"
DeleteFeatureCommand"*
RenameFeatureCommand
name (	Rname"<
ChangeDescriptionCommand 
description (	Rdescription"�
)ChangeBulkUploadSegmentUsersStatusCommand9
status (2!.bucketeer.feature.Segment.StatusRstatus:
state (2$.bucketeer.feature.SegmentUser.StateRstate
count (Rcount"!
AddTagCommand
tag (	Rtag"$
RemoveTagCommand
tag (	Rtag"
EnableFeatureCommand"
DisableFeatureCommand"a
AddVariationCommand
value (	Rvalue
name (	Rname 
description (	Rdescription"(
RemoveVariationCommand
id (	Rid"C
ChangeVariationValueCommand
id (	Rid
value (	Rvalue"@
ChangeVariationNameCommand
id (	Rid
name (	Rname"U
!ChangeVariationDescriptionCommand
id (	Rid 
description (	Rdescription"+
ChangeOffVariationCommand
id (	Rid"?
AddUserToVariationCommand
id (	Rid
user (	Ruser"D
RemoveUserFromVariationCommand
id (	Rid
user (	Ruser"W
ChangeDefaultStrategyCommand7
strategy (2.bucketeer.feature.StrategyRstrategy"=
AddRuleCommand+
rule (2.bucketeer.feature.RuleRrule"}
ChangeRuleStrategyCommand
id (	Rid
rule_id (	RruleId7
strategy (2.bucketeer.feature.StrategyRstrategy"4
ChangeRulesOrderCommand
rule_ids (	RruleIds"#
DeleteRuleCommand
id (	Rid"^
AddClauseCommand
rule_id (	RruleId1
clause (2.bucketeer.feature.ClauseRclause">
DeleteClauseCommand
id (	Rid
rule_id (	RruleId"e
ChangeClauseAttributeCommand
id (	Rid
rule_id (	RruleId
	attribute (	R	attribute"�
ChangeClauseOperatorCommand
id (	Rid
rule_id (	RruleId>
operator (2".bucketeer.feature.Clause.OperatorRoperator"V
AddClauseValueCommand
id (	Rid
rule_id (	RruleId
value (	Rvalue"Y
RemoveClauseValueCommand
id (	Rid
rule_id (	RruleId
value (	Rvalue"�
ChangeFixedStrategyCommand
id (	Rid
rule_id (	RruleId<
strategy (2 .bucketeer.feature.FixedStrategyRstrategy"�
ChangeRolloutStrategyCommand
id (	Rid
rule_id (	RruleId>
strategy (2".bucketeer.feature.RolloutStrategyRstrategy"L
CreateSegmentCommand
name (	Rname 
description (	Rdescription"
DeleteSegmentCommand".
ChangeSegmentNameCommand
name (	Rname"C
ChangeSegmentDescriptionCommand 
description (	Rdescription"n
AddSegmentUserCommand
user_ids (	RuserIds:
state (2$.bucketeer.feature.SegmentUser.StateRstate"q
DeleteSegmentUserCommand
user_ids (	RuserIds:
state (2$.bucketeer.feature.SegmentUser.StateRstate"o
BulkUploadSegmentUsersCommand
data (Rdata:
state (2$.bucketeer.feature.SegmentUser.StateRstate" 
IncrementFeatureVersionCommand"B
CloneFeatureCommand%
environment_id (	RenvironmentIdJ"
ResetSamplingSeedCommand"]
AddPrerequisiteCommandC
prerequisite (2.bucketeer.feature.PrerequisiteRprerequisite":
RemovePrerequisiteCommand

feature_id (	R	featureId"i
"ChangePrerequisiteVariationCommandC
prerequisite (2.bucketeer.feature.PrerequisiteRprerequisite"�
CreateFlagTriggerCommand

feature_id (	R	featureId7
type (2#.bucketeer.feature.FlagTrigger.TypeRtype=
action (2%.bucketeer.feature.FlagTrigger.ActionRaction 
description (	Rdescription"
ResetFlagTriggerCommand"G
#ChangeFlagTriggerDescriptionCommand 
description (	Rdescription"
EnableFlagTriggerCommand"
DisableFlagTriggerCommand"
DeleteFlagTriggerCommand"
UpdateFlagTriggerUsageCommandB1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3
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
�y
proto/feature/service.protobucketeer.featuregoogle/api/annotations.protogoogle/protobuf/wrappers.protoproto/feature/command.protoproto/feature/feature.protoproto/feature/evaluation.protoproto/user/user.protoproto/feature/segment.proto proto/feature/flag_trigger.protoproto/feature/variation.proto proto/feature/prerequisite.protoproto/feature/rule.protoproto/feature/strategy.protoproto/feature/target.proto"P
GetFeatureRequest
id (	Rid%
environment_id (	RenvironmentIdJ"J
GetFeatureResponse4
feature (2.bucketeer.feature.FeatureRfeature"S
GetFeaturesRequest
ids (	Rids%
environment_id (	RenvironmentIdJ"M
GetFeaturesResponse6
features (2.bucketeer.feature.FeatureRfeatures"�
ListFeaturesRequest
	page_size (RpageSize
cursor (	Rcursor
tags (	RtagsI
order_by (2..bucketeer.feature.ListFeaturesRequest.OrderByRorderBy^
order_direction (25.bucketeer.feature.ListFeaturesRequest.OrderDirectionRorderDirection

maintainer (	R
maintainer4
enabled (2.google.protobuf.BoolValueRenabledA
has_experiment	 (2.google.protobuf.BoolValueRhasExperiment%
search_keyword
 (	RsearchKeyword6
archived (2.google.protobuf.BoolValueRarchivedG
has_prerequisites (2.google.protobuf.BoolValueRhasPrerequisites%
environment_id (	RenvironmentId"W
OrderBy
DEFAULT 
NAME

CREATED_AT

UPDATED_AT
TAGS
ENABLED"#
OrderDirection
ASC 
DESCJ"�
ListFeaturesResponse6
features (2.bucketeer.feature.FeatureRfeatures
cursor (	Rcursor
total_count (R
totalCount"�
ListEnabledFeaturesRequest
	page_size (RpageSize
cursor (	Rcursor
tags (	Rtags%
environment_id (	RenvironmentIdJ"m
ListEnabledFeaturesResponse6
features (2.bucketeer.feature.FeatureRfeatures
cursor (	Rcursor"�
CreateFeatureRequestA
command (2'.bucketeer.feature.CreateFeatureCommandRcommand%
environment_id (	RenvironmentIdJ"M
CreateFeatureResponse4
feature (2.bucketeer.feature.FeatureRfeature"�
UpdateFeatureRequest
comment (	Rcomment%
environment_id (	RenvironmentId
id (	Rid0
name (2.google.protobuf.StringValueRname>
description (2.google.protobuf.StringValueRdescription
tags (	Rtags4
enabled (2.google.protobuf.BoolValueRenabled6
archived (2.google.protobuf.BoolValueRarchived<

variations	 (2.bucketeer.feature.VariationR
variationsE
prerequisites
 (2.bucketeer.feature.PrerequisiteRprerequisites3
targets (2.bucketeer.feature.TargetRtargets-
rules (2.bucketeer.feature.RuleRrulesF
default_strategy (2.bucketeer.feature.StrategyRdefaultStrategyA
off_variation (2.google.protobuf.StringValueRoffVariation"M
UpdateFeatureResponse4
feature (2.bucketeer.feature.FeatureRfeature"�
EnableFeatureRequest
id (	RidA
command (2'.bucketeer.feature.EnableFeatureCommandRcommand
comment (	Rcomment%
environment_id (	RenvironmentIdJ"
EnableFeatureResponse"�
DisableFeatureRequest
id (	RidB
command (2(.bucketeer.feature.DisableFeatureCommandRcommand
comment (	Rcomment%
environment_id (	RenvironmentIdJ"
DisableFeatureResponse"�
ArchiveFeatureRequest
id (	RidB
command (2(.bucketeer.feature.ArchiveFeatureCommandRcommand
comment (	Rcomment%
environment_id (	RenvironmentIdJ"
ArchiveFeatureResponse"�
UnarchiveFeatureRequest
id (	RidD
command (2*.bucketeer.feature.UnarchiveFeatureCommandRcommand
comment (	Rcomment%
environment_id (	RenvironmentIdJ"
UnarchiveFeatureResponse"�
DeleteFeatureRequest
id (	RidA
command (2'.bucketeer.feature.DeleteFeatureCommandRcommand
comment (	Rcomment%
environment_id (	RenvironmentIdJ"
DeleteFeatureResponse"�
UpdateFeatureDetailsRequest
id (	Rid]
rename_feature_command (2'.bucketeer.feature.RenameFeatureCommandRrenameFeatureCommandi
change_description_command (2+.bucketeer.feature.ChangeDescriptionCommandRchangeDescriptionCommandJ
add_tag_commands (2 .bucketeer.feature.AddTagCommandRaddTagCommandsS
remove_tag_commands (2#.bucketeer.feature.RemoveTagCommandRremoveTagCommands
comment (	Rcomment%
environment_id (	RenvironmentIdJ"
UpdateFeatureDetailsResponse"�
UpdateFeatureVariationsRequest
id (	Rid6
commands (2.bucketeer.feature.CommandRcommands
comment (	Rcomment%
environment_id (	RenvironmentIdJ"!
UpdateFeatureVariationsResponse"�
UpdateFeatureTargetingRequest
id (	Rid6
commands (2.bucketeer.feature.CommandRcommands
comment (	RcommentI
from (25.bucketeer.feature.UpdateFeatureTargetingRequest.FromRfrom%
environment_id (	RenvironmentId"&
From
UNKNOWN 
USER
OPSJ" 
UpdateFeatureTargetingResponse"�
CloneFeatureRequest
id (	Rid@
command (2&.bucketeer.feature.CloneFeatureCommandRcommand%
environment_id (	RenvironmentIdJ"
CloneFeatureResponse"�
CreateSegmentRequestA
command (2'.bucketeer.feature.CreateSegmentCommandRcommand%
environment_id (	RenvironmentIdJ"M
CreateSegmentResponse4
segment (2.bucketeer.feature.SegmentRsegment"P
GetSegmentRequest
id (	Rid%
environment_id (	RenvironmentIdJ"J
GetSegmentResponse4
segment (2.bucketeer.feature.SegmentRsegment"�
ListSegmentsRequest
	page_size (RpageSize
cursor (	RcursorI
order_by (2..bucketeer.feature.ListSegmentsRequest.OrderByRorderBy^
order_direction (25.bucketeer.feature.ListSegmentsRequest.OrderDirectionRorderDirection%
search_keyword (	RsearchKeyword3
status (2.google.protobuf.Int32ValueRstatusC
is_in_use_status (2.google.protobuf.BoolValueRisInUseStatus%
environment_id	 (	RenvironmentId"@
OrderBy
DEFAULT 
NAME

CREATED_AT

UPDATED_AT"#
OrderDirection
ASC 
DESCJ"�
ListSegmentsResponse6
segments (2.bucketeer.feature.SegmentRsegments
cursor (	Rcursor
total_count (R
totalCount"�
DeleteSegmentRequest
id (	RidA
command (2'.bucketeer.feature.DeleteSegmentCommandRcommand%
environment_id (	RenvironmentIdJ"
DeleteSegmentResponse"�
UpdateSegmentRequest
id (	Rid6
commands (2.bucketeer.feature.CommandRcommands%
environment_id (	RenvironmentIdJ"
UpdateSegmentResponse"�
AddSegmentUserRequest
id (	RidB
command (2(.bucketeer.feature.AddSegmentUserCommandRcommand%
environment_id (	RenvironmentIdJ"
AddSegmentUserResponse"�
DeleteSegmentUserRequest
id (	RidE
command (2+.bucketeer.feature.DeleteSegmentUserCommandRcommand%
environment_id (	RenvironmentIdJ"
DeleteSegmentUserResponse"�
GetSegmentUserRequest

segment_id (	R	segmentId
user_id (	RuserId:
state (2$.bucketeer.feature.SegmentUser.StateRstate%
environment_id (	RenvironmentIdJ"L
GetSegmentUserResponse2
user (2.bucketeer.feature.SegmentUserRuser"�
ListSegmentUsersRequest
	page_size (RpageSize
cursor (	Rcursor

segment_id (	R	segmentId1
state (2.google.protobuf.Int32ValueRstate
user_id (	RuserId%
environment_id (	RenvironmentIdJ"h
ListSegmentUsersResponse4
users (2.bucketeer.feature.SegmentUserRusers
cursor (	Rcursor"�
BulkUploadSegmentUsersRequest

segment_id (	R	segmentIdJ
command (20.bucketeer.feature.BulkUploadSegmentUsersCommandRcommand%
environment_id (	RenvironmentIdJ" 
BulkUploadSegmentUsersResponse"�
BulkDownloadSegmentUsersRequest

segment_id (	R	segmentId:
state (2$.bucketeer.feature.SegmentUser.StateRstate%
environment_id (	RenvironmentIdJ"6
 BulkDownloadSegmentUsersResponse
data (Rdata"�
EvaluateFeaturesRequest(
user (2.bucketeer.user.UserRuser
tag (	Rtag

feature_id (	R	featureId%
environment_id (	RenvironmentIdJ"i
EvaluateFeaturesResponseM
user_evaluations (2".bucketeer.feature.UserEvaluationsRuserEvaluations"�
ListTagsRequest
	page_size (RpageSize
cursor (	RcursorE
order_by (2*.bucketeer.feature.ListTagsRequest.OrderByRorderByZ
order_direction (21.bucketeer.feature.ListTagsRequest.OrderDirectionRorderDirection%
search_keyword (	RsearchKeyword%
environment_id (	RenvironmentId">
OrderBy
DEFAULT 
ID

CREATED_AT

UPDATED_AT"#
OrderDirection
ASC 
DESCJ"w
ListTagsResponse*
tags (2.bucketeer.feature.TagRtags
cursor (	Rcursor
total_count (R
totalCount"�
CreateFlagTriggerRequestj
create_flag_trigger_command (2+.bucketeer.feature.CreateFlagTriggerCommandRcreateFlagTriggerCommand%
environment_id (	RenvironmentIdJ"p
CreateFlagTriggerResponseA
flag_trigger (2.bucketeer.feature.FlagTriggerRflagTrigger
url (	Rurl"�
DeleteFlagTriggerRequest
id (	Ridj
delete_flag_trigger_command (2+.bucketeer.feature.DeleteFlagTriggerCommandRdeleteFlagTriggerCommand%
environment_id (	RenvironmentIdJ"
DeleteFlagTriggerResponse"�
UpdateFlagTriggerRequest
id (	Rid�
'change_flag_trigger_description_command (26.bucketeer.feature.ChangeFlagTriggerDescriptionCommandR#changeFlagTriggerDescriptionCommand%
environment_id (	RenvironmentIdJ"
UpdateFlagTriggerResponse"�
EnableFlagTriggerRequest
id (	Ridj
enable_flag_trigger_command (2+.bucketeer.feature.EnableFlagTriggerCommandRenableFlagTriggerCommand%
environment_id (	RenvironmentIdJ"
EnableFlagTriggerResponse"�
DisableFlagTriggerRequest
id (	Ridm
disable_flag_trigger_command (2,.bucketeer.feature.DisableFlagTriggerCommandRdisableFlagTriggerCommand%
environment_id (	RenvironmentIdJ"
DisableFlagTriggerResponse"�
ResetFlagTriggerRequest
id (	Ridg
reset_flag_trigger_command (2*.bucketeer.feature.ResetFlagTriggerCommandRresetFlagTriggerCommand%
environment_id (	RenvironmentIdJ"o
ResetFlagTriggerResponseA
flag_trigger (2.bucketeer.feature.FlagTriggerRflagTrigger
url (	Rurl"T
GetFlagTriggerRequest
id (	Rid%
environment_id (	RenvironmentIdJ"m
GetFlagTriggerResponseA
flag_trigger (2.bucketeer.feature.FlagTriggerRflagTrigger
url (	Rurl"�
ListFlagTriggersRequest

feature_id (	R	featureId
cursor (	Rcursor
	page_size (RpageSizeM
order_by (22.bucketeer.feature.ListFlagTriggersRequest.OrderByRorderByb
order_direction (29.bucketeer.feature.ListFlagTriggersRequest.OrderDirectionRorderDirection%
environment_id (	RenvironmentId"6
OrderBy
DEFAULT 

CREATED_AT

UPDATED_AT"#
OrderDirection
ASC 
DESCJ"�
ListFlagTriggersResponsec
flag_triggers (2>.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrlRflagTriggers
cursor (	Rcursor
total_count (R
totalCounti
FlagTriggerWithUrlA
flag_trigger (2.bucketeer.feature.FlagTriggerRflagTrigger
url (	Rurl"1
FlagTriggerWebhookRequest
token (	Rtoken"
FlagTriggerWebhookResponse2� 
FeatureService[

GetFeature$.bucketeer.feature.GetFeatureRequest%.bucketeer.feature.GetFeatureResponse" ^
GetFeatures%.bucketeer.feature.GetFeaturesRequest&.bucketeer.feature.GetFeaturesResponse" a
ListFeatures&.bucketeer.feature.ListFeaturesRequest'.bucketeer.feature.ListFeaturesResponse" v
ListEnabledFeatures-.bucketeer.feature.ListEnabledFeaturesRequest..bucketeer.feature.ListEnabledFeaturesResponse" d
CreateFeature'.bucketeer.feature.CreateFeatureRequest(.bucketeer.feature.CreateFeatureResponse" d
UpdateFeature'.bucketeer.feature.UpdateFeatureRequest(.bucketeer.feature.UpdateFeatureResponse" g
EnableFeature'.bucketeer.feature.EnableFeatureRequest(.bucketeer.feature.EnableFeatureResponse"�j
DisableFeature(.bucketeer.feature.DisableFeatureRequest).bucketeer.feature.DisableFeatureResponse"�g
ArchiveFeature(.bucketeer.feature.ArchiveFeatureRequest).bucketeer.feature.ArchiveFeatureResponse" m
UnarchiveFeature*.bucketeer.feature.UnarchiveFeatureRequest+.bucketeer.feature.UnarchiveFeatureResponse" d
DeleteFeature'.bucketeer.feature.DeleteFeatureRequest(.bucketeer.feature.DeleteFeatureResponse" y
UpdateFeatureDetails..bucketeer.feature.UpdateFeatureDetailsRequest/.bucketeer.feature.UpdateFeatureDetailsResponse" �
UpdateFeatureVariations1.bucketeer.feature.UpdateFeatureVariationsRequest2.bucketeer.feature.UpdateFeatureVariationsResponse" 
UpdateFeatureTargeting0.bucketeer.feature.UpdateFeatureTargetingRequest1.bucketeer.feature.UpdateFeatureTargetingResponse" a
CloneFeature&.bucketeer.feature.CloneFeatureRequest'.bucketeer.feature.CloneFeatureResponse" d
CreateSegment'.bucketeer.feature.CreateSegmentRequest(.bucketeer.feature.CreateSegmentResponse" [

GetSegment$.bucketeer.feature.GetSegmentRequest%.bucketeer.feature.GetSegmentResponse" a
ListSegments&.bucketeer.feature.ListSegmentsRequest'.bucketeer.feature.ListSegmentsResponse" d
DeleteSegment'.bucketeer.feature.DeleteSegmentRequest(.bucketeer.feature.DeleteSegmentResponse" d
UpdateSegment'.bucketeer.feature.UpdateSegmentRequest(.bucketeer.feature.UpdateSegmentResponse" j
AddSegmentUser(.bucketeer.feature.AddSegmentUserRequest).bucketeer.feature.AddSegmentUserResponse"�s
DeleteSegmentUser+.bucketeer.feature.DeleteSegmentUserRequest,.bucketeer.feature.DeleteSegmentUserResponse"�j
GetSegmentUser(.bucketeer.feature.GetSegmentUserRequest).bucketeer.feature.GetSegmentUserResponse"�m
ListSegmentUsers*.bucketeer.feature.ListSegmentUsersRequest+.bucketeer.feature.ListSegmentUsersResponse" 
BulkUploadSegmentUsers0.bucketeer.feature.BulkUploadSegmentUsersRequest1.bucketeer.feature.BulkUploadSegmentUsersResponse" �
BulkDownloadSegmentUsers2.bucketeer.feature.BulkDownloadSegmentUsersRequest3.bucketeer.feature.BulkDownloadSegmentUsersResponse" m
EvaluateFeatures*.bucketeer.feature.EvaluateFeaturesRequest+.bucketeer.feature.EvaluateFeaturesResponse" U
ListTags".bucketeer.feature.ListTagsRequest#.bucketeer.feature.ListTagsResponse" p
CreateFlagTrigger+.bucketeer.feature.CreateFlagTriggerRequest,.bucketeer.feature.CreateFlagTriggerResponse" p
UpdateFlagTrigger+.bucketeer.feature.UpdateFlagTriggerRequest,.bucketeer.feature.UpdateFlagTriggerResponse" p
EnableFlagTrigger+.bucketeer.feature.EnableFlagTriggerRequest,.bucketeer.feature.EnableFlagTriggerResponse" s
DisableFlagTrigger,.bucketeer.feature.DisableFlagTriggerRequest-.bucketeer.feature.DisableFlagTriggerResponse" m
ResetFlagTrigger*.bucketeer.feature.ResetFlagTriggerRequest+.bucketeer.feature.ResetFlagTriggerResponse" p
DeleteFlagTrigger+.bucketeer.feature.DeleteFlagTriggerRequest,.bucketeer.feature.DeleteFlagTriggerResponse" g
GetFlagTrigger(.bucketeer.feature.GetFlagTriggerRequest).bucketeer.feature.GetFlagTriggerResponse" m
ListFlagTriggers*.bucketeer.feature.ListFlagTriggersRequest+.bucketeer.feature.ListFlagTriggersResponse" �
FlagTriggerWebhook,.bucketeer.feature.FlagTriggerWebhookRequest-.bucketeer.feature.FlagTriggerWebhookResponse"!���"/webhook/triggers/{token}B1Z/github.com/bucketeer-io/bucketeer/proto/featurebproto3