
�
#proto/environment/environment.protobucketeer.environment"�
EnvironmentV2
id (	Rid
name (	Rname
url_code (	RurlCode 
description (	Rdescription

project_id (	R	projectId
archived (Rarchived

created_at (R	createdAt

updated_at (R	updatedAt'
organization_id	 (	RorganizationId'
require_comment
 (RrequireComment,
feature_flag_count (RfeatureFlagCountB5Z3github.com/bucketeer-io/bucketeer/proto/environmentbproto3
�
proto/environment/project.protobucketeer.environment"�
Project
id (	Rid 
description (	Rdescription
disabled (Rdisabled
trial (Rtrial#
creator_email (	RcreatorEmail

created_at (R	createdAt

updated_at (R	updatedAt
name (	Rname
url_code	 (	RurlCode'
organization_id
 (	RorganizationId+
environment_count (RenvironmentCount,
feature_flag_count (RfeatureFlagCountB5Z3github.com/bucketeer-io/bucketeer/proto/environmentbproto3
�
$proto/environment/organization.protobucketeer.environment"�
Organization
id (	Rid
name (	Rname
url_code (	RurlCode 
description (	Rdescription
disabled (Rdisabled
archived (Rarchived
trial (Rtrial

created_at (R	createdAt

updated_at	 (R	updatedAt!
system_admin
 (RsystemAdmin#
project_count (RprojectCount+
environment_count (RenvironmentCount

user_count (R	userCount
owner_email (	R
ownerEmailB5Z3github.com/bucketeer-io/bucketeer/proto/environmentbproto3
�
!proto/account/search_filter.protobucketeer.account"�
SearchFilter
id (	Rid
name (	Rname
query (	RqueryQ
filter_target_type (2#.bucketeer.account.FilterTargetTypeRfilterTargetType%
environment_id (	RenvironmentId%
default_filter (RdefaultFilter*1
FilterTargetType
UNKNOWN 
FEATURE_FLAGB1Z/github.com/bucketeer-io/bucketeer/proto/accountbproto3
�
proto/account/account.protobucketeer.account#proto/environment/environment.protoproto/environment/project.proto$proto/environment/organization.proto!proto/account/search_filter.proto"�
Account
id (	Rid
email (	Remail
name (	Rname3
role (2.bucketeer.account.Account.RoleRrole
disabled (Rdisabled

created_at (R	createdAt

updated_at (R	updatedAt
deleted (Rdeleted"9
Role

VIEWER 

EDITOR	
OWNER

UNASSIGNEDc:"�
	AccountV2
email (	Remail
name (	Rname(
avatar_image_url (	RavatarImageUrl'
organization_id (	RorganizationId[
organization_role (2..bucketeer.account.AccountV2.Role.OrganizationRorganizationRoleY
environment_roles (2,.bucketeer.account.AccountV2.EnvironmentRoleRenvironmentRoles
disabled (Rdisabled

created_at (R	createdAt

updated_at	 (R	updatedAtF
search_filters
 (2.bucketeer.account.SearchFilterRsearchFilters

first_name (	R	firstName
	last_name (	RlastName
language (	Rlanguage
	last_seen (RlastSeen(
avatar_file_type (	RavatarFileType!
avatar_image (RavatarImage+
environment_count (RenvironmentCount
tags (	Rtags�
Role"Y
Environment
Environment_UNASSIGNED 
Environment_VIEWER
Environment_EDITOR"t
Organization
Organization_UNASSIGNED 
Organization_MEMBER
Organization_ADMIN
Organization_OWNER{
EnvironmentRole%
environment_id (	RenvironmentIdA
role (2-.bucketeer.account.AccountV2.Role.EnvironmentRrole"�
ConsoleAccount
email (	Remail
name (	Rname

avatar_url (	R	avatarUrl&
is_system_admin (RisSystemAdminG
organization (2#.bucketeer.environment.OrganizationRorganization[
organization_role (2..bucketeer.account.AccountV2.Role.OrganizationRorganizationRole^
environment_roles (21.bucketeer.account.ConsoleAccount.EnvironmentRoleRenvironmentRolesF
search_filters (2.bucketeer.account.SearchFilterRsearchFilters

first_name	 (	R	firstName
	last_name
 (	RlastName
language (	Rlanguage(
avatar_file_type (	RavatarFileType!
avatar_image (RavatarImage
	last_seen (RlastSeen�
EnvironmentRoleF
environment (2$.bucketeer.environment.EnvironmentV2Renvironment8
project (2.bucketeer.environment.ProjectRprojectA
role (2-.bucketeer.account.AccountV2.Role.EnvironmentRroleB1Z/github.com/bucketeer-io/bucketeer/proto/accountbproto3
�
proto/account/api_key.protobucketeer.account#proto/environment/environment.proto"�
APIKey
id (	Rid
name (	Rname2
role (2.bucketeer.account.APIKey.RoleRrole
disabled (Rdisabled

created_at (R	createdAt

updated_at (R	updatedAt

maintainer (	R
maintainer
api_key (	RapiKey 
description	 (	Rdescription)
environment_name
 (	RenvironmentName"y
Role
UNKNOWN 

SDK_CLIENT

SDK_SERVER
PUBLIC_API_READ_ONLY
PUBLIC_API_WRITE
PUBLIC_API_ADMIN"�
EnvironmentAPIKey2
api_key (2.bucketeer.account.APIKeyRapiKey1
environment_disabled (RenvironmentDisabled!

project_id (	BR	projectIdF
environment (2$.bucketeer.environment.EnvironmentV2Renvironment(
project_url_code (	RprojectUrlCodeJB1Z/github.com/bucketeer-io/bucketeer/proto/accountbproto3
�
proto/account/command.protobucketeer.accountproto/account/account.protoproto/account/api_key.proto!proto/account/search_filter.proto"�
CreateAccountV2Command
email (	Remail
name (	Rname(
avatar_image_url (	RavatarImageUrl[
organization_role (2..bucketeer.account.AccountV2.Role.OrganizationRorganizationRoleY
environment_roles (2,.bucketeer.account.AccountV2.EnvironmentRoleRenvironmentRoles

first_name (	R	firstName
	last_name (	RlastName
language (	Rlanguage
tags	 (	Rtags"0
ChangeAccountV2NameCommand
name (	Rname"@
ChangeAccountV2FirstNameCommand

first_name (	R	firstName"=
ChangeAccountV2LastNameCommand
	last_name (	RlastName"<
ChangeAccountV2LanguageCommand
language (	Rlanguage"P
$ChangeAccountV2AvatarImageUrlCommand(
avatar_image_url (	RavatarImageUrl"k
ChangeAccountV2AvatarCommand!
avatar_image (RavatarImage(
avatar_file_type (	RavatarFileType"0
ChangeAccountV2TagsCommand
tags (	Rtags"=
ChangeAccountV2LastSeenCommand
	last_seen (RlastSeen"l
&ChangeAccountV2OrganizationRoleCommandB
role (2..bucketeer.account.AccountV2.Role.OrganizationRrole"�
&ChangeAccountV2EnvironmentRolesCommandB
roles (2,.bucketeer.account.AccountV2.EnvironmentRoleRrolesb

write_type (2C.bucketeer.account.ChangeAccountV2EnvironmentRolesCommand.WriteTypeR	writeType"S
	WriteType
WriteType_UNSPECIFIED 
WriteType_OVERRIDE
WriteType_PATCH"
EnableAccountV2Command"
DisableAccountV2Command"
DeleteAccountV2Command"]
CreateAPIKeyCommand
name (	Rname2
role (2.bucketeer.account.APIKey.RoleRrole"-
ChangeAPIKeyNameCommand
name (	Rname"
EnableAPIKeyCommand"
DisableAPIKeyCommand"�
CreateSearchFilterCommand
name (	Rname
query (	RqueryQ
filter_target_type (2#.bucketeer.account.FilterTargetTypeRfilterTargetType%
environment_id (	RenvironmentId%
default_filter (RdefaultFilter"C
ChangeSearchFilterNameCommand
id (	Rid
name (	Rname"F
ChangeSearchFilterQueryCommand
id (	Rid
query (	Rquery"Y
 ChangeDefaultSearchFilterCommand
id (	Rid%
default_filter (RdefaultFilter"+
DeleteSearchFilterCommand
id (	RidB1Z/github.com/bucketeer-io/bucketeer/proto/accountbproto3
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
�
proto/common/string.protobucketeer.common")
StringListValue
values (	RvaluesB0Z.github.com/bucketeer-io/bucketeer/proto/commonbproto3
��
proto/account/service.protobucketeer.accountgoogle/protobuf/wrappers.protogoogle/api/annotations.protogoogle/api/field_behavior.proto.protoc-gen-openapiv2/options/annotations.protoproto/account/account.protoproto/common/string.protoproto/account/api_key.protoproto/account/command.proto$proto/environment/organization.proto"7
GetMeRequest'
organization_id (	RorganizationId"L
GetMeResponse;
account (2!.bucketeer.account.ConsoleAccountRaccount"
GetMyOrganizationsRequest"8
 GetMyOrganizationsByEmailRequest
email (	Remail"g
GetMyOrganizationsResponseI
organizations (2#.bucketeer.environment.OrganizationRorganizations"�
CreateAccountV2Request'
organization_id (	RorganizationIdG
command (2).bucketeer.account.CreateAccountV2CommandBRcommand
email (	Remail
name (	Rname(
avatar_image_url (	RavatarImageUrl[
organization_role (2..bucketeer.account.AccountV2.Role.OrganizationRorganizationRoleY
environment_roles (2,.bucketeer.account.AccountV2.EnvironmentRoleRenvironmentRoles

first_name (	R	firstName
	last_name	 (	RlastName
language
 (	Rlanguage
tags (	Rtags"Q
CreateAccountV2Response6
account (2.bucketeer.account.AccountV2Raccount"�
EnableAccountV2Request
email (	Remail'
organization_id (	RorganizationIdG
command (2).bucketeer.account.EnableAccountV2CommandBRcommand"Q
EnableAccountV2Response6
account (2.bucketeer.account.AccountV2Raccount"�
DisableAccountV2Request
email (	Remail'
organization_id (	RorganizationIdH
command (2*.bucketeer.account.DisableAccountV2CommandBRcommand"R
DisableAccountV2Response6
account (2.bucketeer.account.AccountV2Raccount"�
DeleteAccountV2Request
email (	Remail'
organization_id (	RorganizationIdG
command (2).bucketeer.account.DeleteAccountV2CommandBRcommand"
DeleteAccountV2Response"�
UpdateAccountV2Request
email (	Remail'
organization_id (	RorganizationIda
change_name_command (2-.bucketeer.account.ChangeAccountV2NameCommandBRchangeNameCommandv
change_avatar_url_command (27.bucketeer.account.ChangeAccountV2AvatarImageUrlCommandBRchangeAvatarUrlCommand�
 change_organization_role_command (29.bucketeer.account.ChangeAccountV2OrganizationRoleCommandBRchangeOrganizationRoleCommand�
 change_environment_roles_command (29.bucketeer.account.ChangeAccountV2EnvironmentRolesCommandBRchangeEnvironmentRolesCommandq
change_first_name_command (22.bucketeer.account.ChangeAccountV2FirstNameCommandBRchangeFirstNameCommandn
change_last_name_command (21.bucketeer.account.ChangeAccountV2LastNameCommandBRchangeLastNameCommandm
change_language_command	 (21.bucketeer.account.ChangeAccountV2LanguageCommandBRchangeLanguageCommandn
change_last_seen_command
 (21.bucketeer.account.ChangeAccountV2LastSeenCommandBRchangeLastSeenCommandg
change_avatar_command (2/.bucketeer.account.ChangeAccountV2AvatarCommandBRchangeAvatarCommand0
name (2.google.protobuf.StringValueRnameF
avatar_image_url (2.google.protobuf.StringValueRavatarImageUrll
organization_role (2?.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValueRorganizationRoleY
environment_roles (2,.bucketeer.account.AccountV2.EnvironmentRoleRenvironmentRoles;

first_name (2.google.protobuf.StringValueR	firstName9
	last_name (2.google.protobuf.StringValueRlastName8
language (2.google.protobuf.StringValueRlanguage8
	last_seen (2.google.protobuf.Int64ValueRlastSeenQ
avatar (29.bucketeer.account.UpdateAccountV2Request.AccountV2AvatarRavatar6
disabled (2.google.protobuf.BoolValueRdisabled5
tags (2!.bucketeer.common.StringListValueRtagsa
change_tags_command (2-.bucketeer.account.ChangeAccountV2TagsCommandBRchangeTagsCommand^
AccountV2Avatar!
avatar_image (RavatarImage(
avatar_file_type (	RavatarFileType[
OrganizationRoleValueB
role (2..bucketeer.account.AccountV2.Role.OrganizationRrole"Q
UpdateAccountV2Response6
account (2.bucketeer.account.AccountV2Raccount"T
GetAccountV2Request
email (	Remail'
organization_id (	RorganizationId"N
GetAccountV2Response6
account (2.bucketeer.account.AccountV2Raccount"a
"GetAccountV2ByEnvironmentIDRequest
email (	Remail%
environment_id (	RenvironmentId"]
#GetAccountV2ByEnvironmentIDResponse6
account (2.bucketeer.account.AccountV2Raccount"�
ListAccountsV2Request
	page_size (RpageSize
cursor (	Rcursor'
organization_id (	RorganizationIdK
order_by (20.bucketeer.account.ListAccountsV2Request.OrderByRorderBy`
order_direction (27.bucketeer.account.ListAccountsV2Request.OrderDirectionRorderDirection%
search_keyword (	RsearchKeyword6
disabled (2.google.protobuf.BoolValueRdisabledH
organization_role (2.google.protobuf.Int32ValueRorganizationRoleC
environment_id	 (2.google.protobuf.StringValueRenvironmentIdF
environment_role
 (2.google.protobuf.Int32ValueRenvironmentRole
tags (	Rtags"�
OrderBy
DEFAULT 	
EMAIL

CREATED_AT

UPDATED_AT
ORGANIZATION_ROLE
ENVIRONMENT_COUNT
	LAST_SEEN	
STATE
TAGS"#
OrderDirection
ASC 
DESC"�
ListAccountsV2Response8
accounts (2.bucketeer.account.AccountV2Raccounts
cursor (	Rcursor
total_count (R
totalCount"�
CreateAPIKeyRequestD
command (2&.bucketeer.account.CreateAPIKeyCommandBRcommand%
environment_id (	RenvironmentId
name (	Rname2
role (2.bucketeer.account.APIKey.RoleRrole

maintainer (	R
maintainer 
description (	RdescriptionJ"J
CreateAPIKeyResponse2
api_key (2.bucketeer.account.APIKeyRapiKey"�
ChangeAPIKeyNameRequest
id (	RidD
command (2*.bucketeer.account.ChangeAPIKeyNameCommandRcommand%
environment_id (	RenvironmentIdJ"
ChangeAPIKeyNameResponse"�
EnableAPIKeyRequest
id (	Rid@
command (2&.bucketeer.account.EnableAPIKeyCommandRcommand%
environment_id (	RenvironmentIdJ"
EnableAPIKeyResponse"�
DisableAPIKeyRequest
id (	RidA
command (2'.bucketeer.account.DisableAPIKeyCommandRcommand%
environment_id (	RenvironmentIdJ"
DisableAPIKeyResponse"O
GetAPIKeyRequest
id (	Rid%
environment_id (	RenvironmentIdJ"G
GetAPIKeyResponse2
api_key (2.bucketeer.account.APIKeyRapiKey"�
ListAPIKeysRequest
	page_size (RpageSize
cursor (	RcursorH
order_by (2-.bucketeer.account.ListAPIKeysRequest.OrderByRorderBy]
order_direction (24.bucketeer.account.ListAPIKeysRequest.OrderDirectionRorderDirection%
search_keyword (	RsearchKeyword6
disabled (2.google.protobuf.BoolValueRdisabled)
environment_id (	BRenvironmentId'
environment_ids	 (	RenvironmentIds,
organization_id
 (	B�ARorganizationId"f
OrderBy
DEFAULT 
NAME

CREATED_AT

UPDATED_AT
ROLE
ENVIRONMENT	
STATE"#
OrderDirection
ASC 
DESCJ"�
ListAPIKeysResponse4
api_keys (2.bucketeer.account.APIKeyRapiKeys
cursor (	Rcursor
total_count (R
totalCount"6
GetEnvironmentAPIKeyRequest
api_key (	RapiKey"t
GetEnvironmentAPIKeyResponseT
environment_api_key (2$.bucketeer.account.EnvironmentAPIKeyRenvironmentApiKey"�
CreateSearchFilterRequest
email (	Remail'
organization_id (	RorganizationId%
environment_id (	RenvironmentIdF
command (2,.bucketeer.account.CreateSearchFilterCommandRcommand"
CreateSearchFilterResponse"�
UpdateSearchFilterRequest
email (	Remail'
organization_id (	RorganizationId%
environment_id (	RenvironmentId`
change_name_command (20.bucketeer.account.ChangeSearchFilterNameCommandRchangeNameCommandc
change_query_command (21.bucketeer.account.ChangeSearchFilterQueryCommandRchangeQueryCommandv
change_default_filter_command (23.bucketeer.account.ChangeDefaultSearchFilterCommandRchangeDefaultFilterCommand"
UpdateSearchFilterResponse"�
DeleteSearchFilterRequest
email (	Remail'
organization_id (	RorganizationId%
environment_id (	RenvironmentIdF
command (2,.bucketeer.account.DeleteSearchFilterCommandRcommand"
DeleteSearchFilterResponse"�
UpdateAPIKeyRequest
id (	Rid%
environment_id (	RenvironmentId0
name (2.google.protobuf.StringValueRname>
description (2.google.protobuf.StringValueRdescription2
role (2.bucketeer.account.APIKey.RoleRrole6
disabled (2.google.protobuf.BoolValueRdisabled<

maintainer (2.google.protobuf.StringValueR
maintainer"
UpdateAPIKeyResponse2�z
AccountService�
GetMe.bucketeer.account.GetMeRequest .bucketeer.account.GetMeResponse"��A�
AccountGet MeGet the user console account.*web.v1.account.get_meJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }���/v1/account/get_me�
GetMyOrganizations,.bucketeer.account.GetMyOrganizationsRequest-.bucketeer.account.GetMyOrganizationsResponse"��A�
AccountGet My Organizations.Get all the organizations for a specific user.*web.v1.account.my_organizationsJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }���/v1/account/my_organizations�
GetMyOrganizationsByEmail3.bucketeer.account.GetMyOrganizationsByEmailRequest-.bucketeer.account.GetMyOrganizationsResponse"��A�
AccountGet My Organizations By Email#Get all the organizations by email.*(web.v1.account.my_organizations_by_emailJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }���'%/v1/account/my_organizations_by_email�
CreateAccountV2).bucketeer.account.CreateAccountV2Request*.bucketeer.account.CreateAccountV2Response"��A�
AccountCreate\Create an account to have access to the console. To call this API, you need an `ADMIN` role.*web.v1.account.create_accountJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }���"/v1/account/create_account:*�
EnableAccountV2).bucketeer.account.EnableAccountV2Request*.bucketeer.account.EnableAccountV2Response"��A�
AccountEnable\Enable an account to have access to the console. To call this API, you need an `ADMIN` role.*web.v1.account.enable_accountJ�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }���"/v1/account/enable_account:*�
DisableAccountV2*.bucketeer.account.DisableAccountV2Request+.bucketeer.account.DisableAccountV2Response"��A�
AccountDisable^Disable an account to block access to the console. To call this API, you need an `ADMIN` role.*web.v1.account.disable_accountJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }��� "/v1/account/disable_account:*�
UpdateAccountV2).bucketeer.account.UpdateAccountV2Request*.bucketeer.account.UpdateAccountV2Response"��A�
AccountUpdate>Update an account. To call this API, you need an `ADMIN` role.*web.v1.account.update_accountJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }���"/v1/account/update_account:*�
DeleteAccountV2).bucketeer.account.DeleteAccountV2Request*.bucketeer.account.DeleteAccountV2Response"��A�
AccountDelete>Delete an account. To call this API, you need an `ADMIN` role.*web.v1.account.delete_accountJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }���"/v1/account/delete_account:*�
GetAccountV2&.bucketeer.account.GetAccountV2Request'.bucketeer.account.GetAccountV2Response"��A�
AccountGetGet an account.*web.v1.account.get_accountJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }���/v1/account/get_account�
GetAccountV2ByEnvironmentID5.bucketeer.account.GetAccountV2ByEnvironmentIDRequest6.bucketeer.account.GetAccountV2ByEnvironmentIDResponse"��A�
AccountGet Account By EnvironmentGet an account by environment.*)web.v1.account.get_account_by_environmentJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }���(&/v1/account/get_account_by_environment�
ListAccountsV2(.bucketeer.account.ListAccountsV2Request).bucketeer.account.ListAccountsV2Response"��A�
AccountListList accounts.*web.v1.account.list_accountsJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }���/v1/account/list_accounts�
CreateAPIKey&.bucketeer.account.CreateAPIKeyRequest'.bucketeer.account.CreateAPIKeyResponse"��A�
API KeyCreate[Create an API key to be used on the client SDK. To call this API, you need an `ADMIN` role.*web.v1.account.create_api_keyJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }���"/v1/account/create_api_key:*�
ChangeAPIKeyName*.bucketeer.account.ChangeAPIKeyNameRequest+.bucketeer.account.ChangeAPIKeyNameResponse"��A�
API KeyChange API Key NameDChange the API Key Name. To call this API, you need an `ADMIN` role.*"web.v1.account.change_api_key_nameJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }���$"/v1/account/change_api_key_name:*�
EnableAPIKey&.bucketeer.account.EnableAPIKeyRequest'.bucketeer.account.EnableAPIKeyResponse"��A�
API KeyEnable>Enable an API Key. To call this API, you need an `ADMIN` role.*web.v1.account.enable_api_keyJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }���"/v1/account/enable_api_key:*�
DisableAPIKey'.bucketeer.account.DisableAPIKeyRequest(.bucketeer.account.DisableAPIKeyResponse"��A�
API KeyDisable?Disable an API Key. To call this API, you need an `ADMIN` role.*web.v1.account.disable_api_keyJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }��� "/v1/account/disable_api_key:*�
	GetAPIKey#.bucketeer.account.GetAPIKeyRequest$.bucketeer.account.GetAPIKeyResponse"��A�
API KeyGetGet an API Key.*web.v1.account.get_api_keyJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }���/v1/account/get_api_key�
ListAPIKeys%.bucketeer.account.ListAPIKeysRequest&.bucketeer.account.ListAPIKeysResponse"��A�
API KeyListList API Keys.*web.v1.account.list_api_keysJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }���/v1/account/list_api_keys�
GetEnvironmentAPIKey..bucketeer.account.GetEnvironmentAPIKeyRequest/.bucketeer.account.GetEnvironmentAPIKeyResponse"��A�
API KeyGet Environment API KeyGet an environment API Key.*&web.v1.account.get_environment_api_keyJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }���%#/v1/account/get_environment_api_key�
CreateSearchFilter,.bucketeer.account.CreateSearchFilterRequest-.bucketeer.account.CreateSearchFilterResponse"��A�
AccountCreate Search FilterDCreate a search filter. To call this API, you need an `VIEWER` role.*#web.v1.account.create_search_filterJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }���%" /v1/account/create_search_filter:*�
UpdateSearchFilter,.bucketeer.account.UpdateSearchFilterRequest-.bucketeer.account.UpdateSearchFilterResponse"��A�
AccountUpdate Search FilterDUpdate a search filter. To call this API, you need an `VIEWER` role.*#web.v1.account.update_search_filterJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }���%" /v1/account/update_search_filter:*�
DeleteSearchFilter,.bucketeer.account.DeleteSearchFilterRequest-.bucketeer.account.DeleteSearchFilterResponse"��A�
AccountDelete Search FilterDDelete a search filter. To call this API, you need an `VIEWER` role.*#web.v1.account.delete_search_filterJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }���%" /v1/account/delete_search_filter:*�
UpdateAPIKey&.bucketeer.account.UpdateAPIKeyRequest'.bucketeer.account.UpdateAPIKeyResponse"��A�
API KeyUpdate API Key>Update an API Key. To call this API, you need an `ADMIN` role.*web.v1.account.update_api_keyJ�
400�
:Returned for bad requests that may have failed validation.
.google.rpc.Status"V
application/jsonB{ "code": 3, "message": "invalid arguments error", "details": [] }J�
401�
=Request could not be authenticated (authentication required).
.google.rpc.Status"Q
application/json={ "code": 16, "message": "not authenticated", "details": [] }J�
403�
8Request does not have permission to access the resource.
.google.rpc.Status"M
application/json9{ "code": 7, "message": "not authorized", "details": [] }���2/v1/account/update_api_key:*B1Z/github.com/bucketeer-io/bucketeer/proto/accountbproto3