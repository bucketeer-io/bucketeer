# Summary

Audit logs is a feature to store the history of changes that we made to objects like Feature Flags, Push, API keys, … 
We can see the user who did the action in either the console or public APIs.

<div style="text-align: center;">
  <img width="50%" src="/docs/rfcs/images/audit-log-dashboard.png" alt="audit log dashboard">
</div>

When this document was drafting, the editor who did the action is got from access token: 
https://github.com/bucketeer-io/bucketeer/blob/main/pkg/role/role.go#L48-L74

For public APIs, the authorization method is API keys, so after our API gateway request to web gRPC service, 
the access token is got from internal environment, which is the service token 
(https://github.com/bucketeer-io/bucketeer/blob/main/pkg/api/cmd/server.go#L114), so the email being shown 
in audit log for public API write action is always internal Bucketeer email, which is not really correct.

<div style="text-align: center;">
  <img width="50%" src="/docs/rfcs/images/public-api-default-flow.png" alt="public api default flow">
</div>

We need to save the signature (in this case, email) of the API creator and decide who should be the editor 
then apply a mechanism to get the correct editor in api layer when making calls to the Web gRPC service.

# Solutions

## Problem 1: save the creator of API key when we create API key in Console

We can add a column `created_by` in api_key table so the creator can be saved. For old api keys, we can find creator
in audit logs by filter type = 400 (create API key event) in audit_log table, select editor and `entity_id` 
⇒ migration data for old api_key records and migration data for old `audig_log` records.

<div style="text-align: center">
  <img width="15%" src="/docs/rfcs/images/api-key-table-add-created-by.png" alt="api_key table ERD" style="padding-right: 10px">
  <img width="35%" src="/docs/rfcs/images/audit-log-snapshot.png" alt="audit log snapshot example">
</div>

We also need to add API key owner/creator in the UI and in the CreatorAPIkey web API request message:

```protobuf
message CreateAPIKeyRequest {
  string environment_namespace = 1;
  string name = 2;
  account.APIKey.Role role = 3;
  string created_by = 4;
}
```

As we add new property to `api_key` object, the response of the `GetAPIKeyBySearchingAllEnvironments` API 
and `ListAPIKeys` API also return an extra `created_by` field in `api_key` object:

```protobuf
message EnvironmentAPIKey {
  string environment_namespace = 1 [deprecated = true];
  APIKey api_key = 2;
  bool environment_disabled = 3;
  string project_id = 4 [deprecated = true];
  environment.EnvironmentV2 environment = 5;
  string project_url_code = 6;
  string created_by = 7;
}
```

## Problem 2: Get the correct editor in API layer when making call to Web service
I suggest 2 options:
  - Let the editor be the creator of the API key.
  - We let the client decide the editor of the action by adding editor email in the request body of public API.

### Option 1: Let the editor be the creator of the API key <span style="color:green">(preferred)</span>

We can extract the APIKey creator in function `getEnvironmentAPIKey` as we will update the response of 
`GetAPIKeyBySearchingAllEnvironments` (https://github.com/bucketeer-io/bucketeer/blob/main/pkg/account/api/api_key.go#L455)

Now to save the API key creator to audit log, we need to overwrite the editor, 
the idea is to use context to share information between services.

The `created_by` value will be added into context before forward to Web gRPC service:

```go
const APIKeyCreatorMDKey string = "apikey-creator"

headerMetadata := metadata.New(map[string]string{
	APIKeyCreatorMDKey: envAPIKey.ApiKey.CreatedBy,
})
ctx = metadata.NewOutgoingContext(ctx, headerMetadata)
```

When receive request, we can get it from context metadata again 
(let’s add this in https://github.com/bucketeer-io/bucketeer/blob/main/pkg/role/role.go#L48-L74):

```go
md, ok := metadata.FromIncomingContext(ctx)
if ok {
    apiKeyCreator := md.Get(APIKeyCreatorMDKey)
    // verify creator email then form *eventproto.Editor ...
}
```

Below is the overall updated flow to save the API key creator to audit log, 
red color means those are new flows need to be implemented in this option:

<div style="text-align: center;">
  <img width="50%" src="/docs/rfcs/images/option-1-get-audit-log-editor.png" alt="option 1 flow">
</div>

**Pros and cons**

Pros:

- We don’t have to change the signature of any public or web API, no breaking change.
- Other than saving API keys and create migration for old data in the database, we only need to add 
API creator to context before forward to web gRPC service (for every APIs) and overwrite the editor.
- The implementation can be fast and compact.

Cons:

- The editor is not fully specific, we only know the one that responsible for creating the API, 
not the one that actually do the action.

### Option 2: let the client decide the editor of the action

In this option, we still save the API creator in the audit log but also, we change the request body 
of update and create APIs by adding `creator_email` or `updater_email`:

```protobuf
message UpdateFeatureRequest {
  string comment = 1;
  string id = 2;
  google.protobuf.StringValue name = 3;
  google.protobuf.StringValue description = 4;
  repeated string tags = 5;
  google.protobuf.BoolValue enabled = 6;
  google.protobuf.BoolValue archived = 7;
  repeated feature.Variation variations = 8;
  repeated feature.Prerequisite prerequisites = 9;
  repeated feature.Target targets = 10;
  repeated feature.Rule rules = 11;
  feature.Strategy default_strategy = 12;
  google.protobuf.StringValue off_variation = 13;
  google.protobuf.StringValue updater_email = 14;
}
```

This solution is an extension of option 1, and we can reconstruct the audit log message for public API like this:
```
API creator: api-creator@bucketeer.io, Editor: editor@gmail.com
Push name has been updated
```

The code should change like this:
```go
// get editor and pass to context
const (
    APIKeyCreatorMDKey = "apikey-creator"
    APIEditorMDKey     = "api-editor"
)

headerMetadata := metadata.New(map[string]string{
    APIKeyCreatorMDKey: envAPIKey.ApiKey.CreatedBy,
    APIEditorMDKey:     req.updater_email.value,
})
ctx = metadata.NewOutgoingContext(ctx, headerMetadata)
```
```go
// get editor from context
md, ok := metadata.FromIncomingContext(ctx)
if ok {
    apiKeyCreator := md.Get(APIKeyCreatorMDKey)
    apiEditor := md.Get(APIEditorMDKey)
    
    editors := fmt.Sprint("API creator: %s, Editor: %s",
        apiKeyCreator,
        apiEditor,
    )
    // verify creator email then form *eventproto.Editor ...
}
```
*The reason that we save both API triggerer and API editor but not one of them is because in case 
input `updater_email` is “anonymous@hacker.io” or something like that, we still do not know 
exactly who did the action. The target here should be what has happened and who did it.*

The new API flow should be below:
<div style="text-align: center;">
  <img width="50%" src="/docs/rfcs/images/option-2-get-audit-log-editor.png" alt="option 2 flow">
</div>

**Pros and cons:**

Pros:
- We can define exactly who is the editor of the object, the content is more specific.

Cons:
- We have to change the request of every current public APIs and their associated description.
- The audit log message format might be changed as well.
- For old audit log data, we can only know the editor is the API key creator (like option 1).
- The implementation is more time-consuming than option 1 as more changes.

# Conclusion

The tasks and time estimation to resolve this issue can be:

| Task                                            | Description                                                                                                                                                                                                                                                      | Time estimation |
|:------------------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:---------------:|
| Update APIKey APIs (Breaking change)            | - Add column `created_by` in `api_key` table <br/>- Update `CreateAPIKey`: add `created_by` in request and response; Update `ListAPIKeys` and `GetAPIKeyBySearchingAllEnvironments`: add `created_by` in `api_key` object response                               |   1 - 2 days    |
| Update APIKey UI (Breaking change)              | - Add field email input as the APIKey owner                                                                                                                                                                                                                      |                 |
| Create data migration for old audit log data    | - Create backup <br/>- Query `audit_log` table, get all records with type = 400 (create API key event) and get all editors as API creator of associated API <br/>- Migrate old event that use APIKey in `audit_log` to new editor (APIKey creator)               |  1 - 2.5 days   |
| (Option 1) Implement overwrite editor mechanism | - Get API creator of APIKey and save to context before call web gRPC API (every create/update public APIs) <br/>- Overwrite editor if API creator in context metadata is not nil                                                                                 |  1 - 2.5 days   |
| (Option 2) Implement overwrite editor mechanism | - Change request body message of every create/update public APIs, also change the API description document <br/>- Get API creator of APIKey, form editors then forward to gRPC web service <br/>- Overwrite editor if API creator in context metadata is not nil |   3 - 5 days    |
