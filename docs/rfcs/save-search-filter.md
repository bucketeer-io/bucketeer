# Save Filters Summary

"Save Filters" saves the filters set for each feature list in the console.

This feature allows to give default settings to saved filters.

This feature should be designed to store filter information for the various information displayed in the console.

Filter and search information is based on the string in the query part of the URL.

ex) https://dev.bucketeer.jp/default/features?page=1&q=test&sort=-createdAt&tagIds=iOS

Filter and Search Info: ```q=test&sort=-createdAt&tagIds=iOS```

**Concerns**
- Since the string we save is a query parameter, the saved query may not work if the filter is changed.
   - ex: Save Query String [&sort=-createdAt&tagIds=iOS]
     If there is no longer a filter for Tag, saved queries will not work.

I have considered the following two methods, but considering the advantages and disadvantages, we will propose ***Part2***.

## Part1: Create a new search_filter table and manage one record for each filter.

### Package
packageName: savefilter

### API

In common with all functions, instead of receiving an email address from the user to identify the user, we use the [AccessToken.Email](https://github.com/bucketeer-io/bucketeer/blob/main/pkg/token/token.go#L26) obtained with [GetAccessToken](https://github.com/bucketeer-io/bucketeer/blob/41da1916100bf29cc925010b629f74c17e014d5f/pkg/rpc/auth.go#L70).

#### Create
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

#### Update
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

#### Delete
Receives a request for the FilterID to be deleted and deletes the corresponding filter.
* Request
  * ID: Filter ID
* Request Validation
  * ID: Reject ID that does not exist
  * Compare the saved email and the email obtained from the AccessToken and reject if they are different.

#### List
A list is retrieved for each filter target(FeatureFlag, Goal, etc.).
* Request
  * TargetType: Target feature type to filter (FeatureFlag, Goal, etc.)
  * EnvironmentId: Target Environment Id
* Request Validation
  * TargetType: Reject Unknown Type

#### Get
Filter information is obtained based on the filter ID.
  * ID: Filter ID
* Request Validation
  * ID: Reject ID that does not exist
  * Compare the saved email and the email obtained from the AccessToken and reject if they are different.

### Implementation

#### Storage

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

### Merit
- Since the `target_type` is managed by column, there is no need to parse compared to managing by JSON, so the cost of narrowing down records according to conditions is lower.

### Demerit
- Records are maintained for each user, so as the number of users increases, the number of records increases, leading to a decrease in performance.
- Since information is linked to accounts, when deleting an account, records related to the account must also be deleted, which increases management costs.

## Part2: Add and manage a new column that stores filter information in an array in the account_v2 table

### Package
Add functionality within the `account` package.

### Proto
```
message AccountV2 {
  string email = 1;
  string name = 2;
  string avatar_image_url = 3;
  string organization_id = 4;
  Role.Organization organization_role = 5;
  repeated EnvironmentRole environment_roles = 6;
  bool disabled = 7;
  int64 created_at = 8;
  int64 updated_at = 9;
  repeated SearchFilter search_filters = 10;  // ← Add
}
```

```
message SearchFilter {
  string id = 1;
  string name = 2;
  string query = 3;
  FilterTargetType filter_target_type = 4;
  string environment_id = 5;
  bool default_filter = 6;
  int64 created_at = 7;
  int64 updated_at = 8;
}
```

```
enum FilterTargetType {
  UNNKOWN = 0;
  FEATURE_FLAG = 1;
  GOAL = 2;
}
```

### Storage
Add `search_filters` column to `account_v2` table.

```
CREATE TABLE `account_v2` (
  `email` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `avatar_image_url` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `organization_id` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `organization_role` int NOT NULL,
  `environment_roles` json NOT NULL,
  `search_filters` json,        -- ← Add column
  `disabled` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` bigint NOT NULL,
  `updated_at` bigint NOT NULL,
  PRIMARY KEY (`email`,`organization_id`),
  KEY `account_v2_foreign_organization_id` (`organization_id`),
  CONSTRAINT `account_v2_foreign_organization_id` FOREIGN KEY (`organization_id`) REFERENCES `organization` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
```

### API
It is assumed that GetAccount is used to retrieve saved filters.

#### CreateSearchFilter
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

#### UpdateSearchFilter
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

#### DeleteSearchFilter
Receives a request for the FilterID to be deleted and deletes the corresponding filter.
* Request
  * ID: Filter ID
* Request Validation
  * ID: Reject ID that does not exist

### Merit
- Management costs are low because information can be managed by linking it to an account.
- It requires less man-hours because it is additionally implemented within the existing Account package processing.

### Demerit
- Since it is handled as an array-type Json, it is not possible to specify conditions in a query, and after retrieving all results, it is necessary to extract a filter that matches the conditions, resulting in high search costs.
  - Since users are not expected to save dozens of filters, we do not think it will have much of an impact.
