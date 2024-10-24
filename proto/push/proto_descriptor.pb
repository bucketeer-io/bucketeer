
�
proto/push/command.protobucketeer.push"q
CreatePushCommand
tags (	Rtags
name (	Rname.
fcm_service_account (RfcmServiceAccountJ"(
AddPushTagsCommand
tags (	Rtags"+
DeletePushTagsCommand
tags (	Rtags"
DeletePushCommand"'
RenamePushCommand
name (	RnameB.Z,github.com/bucketeer-io/bucketeer/proto/pushbproto3
�
proto/push/push.protobucketeer.push"�
Push
id (	Rid
tags (	Rtags
deleted (Rdeleted
name (	Rname

created_at (R	createdAt

updated_at (R	updatedAt.
fcm_service_account (	RfcmServiceAccountJB.Z,github.com/bucketeer-io/bucketeer/proto/pushbproto3
�
proto/push/service.protobucketeer.pushproto/push/push.protoproto/push/command.proto"�
CreatePushRequest3
environment_namespace (	RenvironmentNamespace?
command (2!.bucketeer.push.CreatePushCommandBRcommand
tags (	Rtags
name (	Rname.
fcm_service_account (RfcmServiceAccount">
CreatePushResponse(
push (2.bucketeer.push.PushRpush"�
ListPushesRequest3
environment_namespace (	RenvironmentNamespace
	page_size (RpageSize
cursor (	RcursorD
order_by (2).bucketeer.push.ListPushesRequest.OrderByRorderByY
order_direction (20.bucketeer.push.ListPushesRequest.OrderDirectionRorderDirection%
search_keyword (	RsearchKeyword"@
OrderBy
DEFAULT 
NAME

CREATED_AT

UPDATED_AT"#
OrderDirection
ASC 
DESC"{
ListPushesResponse,
pushes (2.bucketeer.push.PushRpushes
cursor (	Rcursor
total_count (R
totalCount"�
DeletePushRequest3
environment_namespace (	RenvironmentNamespace
id (	Rid;
command (2!.bucketeer.push.DeletePushCommandRcommand"
DeletePushResponse"�
UpdatePushRequest3
environment_namespace (	RenvironmentNamespace
id (	RidU
add_push_tags_command (2".bucketeer.push.AddPushTagsCommandRaddPushTagsCommand^
delete_push_tags_command (2%.bucketeer.push.DeletePushTagsCommandRdeletePushTagsCommandQ
rename_push_command (2!.bucketeer.push.RenamePushCommandRrenamePushCommand"
UpdatePushResponse2�
PushServiceU

ListPushes!.bucketeer.push.ListPushesRequest".bucketeer.push.ListPushesResponse" U

CreatePush!.bucketeer.push.CreatePushRequest".bucketeer.push.CreatePushResponse" U

DeletePush!.bucketeer.push.DeletePushRequest".bucketeer.push.DeletePushResponse" U

UpdatePush!.bucketeer.push.UpdatePushRequest".bucketeer.push.UpdatePushResponse" B.Z,github.com/bucketeer-io/bucketeer/proto/pushbproto3