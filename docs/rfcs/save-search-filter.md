# Save Filters Summary

"Save Filters" saves the filters set for each feature list in the console.
This feature allows you to give default settings to saved filters.
This feature should be designed to store filter information for the various information displayed in the console.

## package
packageName: savefilter

## Function
### Create
Receives and saves requests for filter titles, query strings, and default settings flags.
* Request
  * Name: Filter Name Stirng
  * Query: Query Parameters Stirng
  * DefaultFlag: Filter flags to set by default
  * TargetType: Target feature type to filter (FeatureFlag, Goal, etc.)
* Request Validation
  * Name: empty string(Same name is allowed)
  * Query: empty string(Same guary is allowed)
  * DefaultFlag: Noting
  * TargetType: Unknown Type

### Update
Receives and updates requests for filter titles, query strings, and default settings flags.
* Request
  * ID: Filter ID
  * Name: Filter Name Stirng
  * Query: Query Parameters Stirng
  * DefaultFlag: Filter flags to set by default
* Request Validation
  * ID: ID that does not exist
  * Name: empty string(Same name is allowed)
  * Query: empty string(Same name is allowed)
  * DefaultFlag: Noting

### Delete
Receives a request for the FilterID to be deleted and deletes the corresponding filter.
* Request
  * ID: Filter ID
* Request Validation
  * ID: ID that does not exist

### List
A list is retrieved for each filter target(FeatureFlag, Goal, etc.).
* Request
  * TargetType: Target feature type to filter (FeatureFlag, Goal, etc.)
* Request Validation
  * TargetType: Unknown Type

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
        `created_at` bigint NOT NULL,
        `updated_at` bigint NOT NULL,
      ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
```

## Concerns
1. Since the string you save is a query parameter, the saved query may not work if the filter is changed.
   - ex: Save Query String [&sort=-createdAt&tagIds=iOS]
     If there is no longer a filter for Tag, saved queries will not work.
