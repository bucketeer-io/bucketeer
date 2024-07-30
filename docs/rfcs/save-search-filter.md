# Save Filters Summary

"Save Filters" saves the filters set for each feature list in the console.

This feature allows you to give default settings to saved filters.

This feature should be designed to store filter information for the various information displayed in the console.

## package
packageName: savefilter

## Function
Filter and search information is based on the string in the query part of the URL.

ex) https://dev.bucketeer.jp/default/features?page=1&q=test&sort=-createdAt&tagIds=iOS

Filter and Search Info: ```q=test&sort=-createdAt&tagIds=iOS```

In common with all functions, instead of receiving an email address from the user to identify the user, we use the [AccessToken.Email](https://github.com/bucketeer-io/bucketeer/blob/main/pkg/token/token.go#L26) obtained with [GetAccessToken](https://github.com/bucketeer-io/bucketeer/blob/41da1916100bf29cc925010b629f74c17e014d5f/pkg/rpc/auth.go#L70).

### Create
Receives and saves requests for filter titles, query strings, and default settings flags.
* Request
  * Name: Filter Name String
  * Query: Query Parameters String
  * DefaultFilter: Filter to set by default
  * TargetType: Target feature type to filter (FeatureFlag, Goal, etc.)
  * EnvironmentId: Target Environment Id
* Request Validation
  * Name: Reject empty string(Same name is allowed)
  * Query: Reject empty string(Same query is allowed)
  * DefaultFilter: ー
  * TargetType: Reject Unknown Type
  * EnvironmentId: check EnvironmentRole

### Update
Receives and updates requests for filter titles, query strings, and default settings flags.
* Request
  * ID: Filter ID
  * Name: Filter Name String
  * Query: Query Parameters String
  * DefaultFilter: Filter to set by default
* Request Validation
  * ID: ID that does not exist
  * Name: Reject empty string(Same name is allowed)
  * Query: Reject empty string(Same name is allowed)
  * DefaultFilter: ー
  * Compare the saved email and the email obtained from the AccessToken and reject if they are different.

### Delete
Receives a request for the FilterID to be deleted and deletes the corresponding filter.
* Request
  * ID: Filter ID
* Request Validation
  * ID: Reject ID that does not exist
  * Compare the saved email and the email obtained from the AccessToken and reject if they are different.

### List
A list is retrieved for each filter target(FeatureFlag, Goal, etc.).
* Request
  * TargetType: Target feature type to filter (FeatureFlag, Goal, etc.)
  * EnvironmentId: Target Environment Id
* Request Validation
  * TargetType: Reject Unknown Type

### Get
Filter information is obtained based on the filter ID.
  * ID: Filter ID
* Request Validation
  * ID: Reject ID that does not exist
  * Compare the saved email and the email obtained from the AccessToken and reject if they are different.

## Implementation

### Storage

* Create `search_filter` table definition.

```sql
      CREATE TABLE `search_filter` (
        `id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
        `name` sting NOT NULL,
        `query` sting NOT NULL,
        `target_type` int NOT NULL,
        `default` tinyint(1) NOT NULL DEFAULT '0',
        `environment_id` sting NOT NULL,
        `account_email` sting NOT NULL,
        `created_at` bigint NOT NULL,
        `updated_at` bigint NOT NULL,
        FOREIGN KEY (account_email) REFERENCES account_v2(email) ON DELETE CASCADE
      ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
```

## Concerns
1. Since the string you save is a query parameter, the saved query may not work if the filter is changed.
   - ex: Save Query String [&sort=-createdAt&tagIds=iOS]
     If there is no longer a filter for Tag, saved queries will not work.
