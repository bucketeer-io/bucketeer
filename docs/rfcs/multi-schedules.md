# Summary

The Multi Schedule function allows you to register multiple schedules at once.
It is managed together like a staged "progressive rollout".

## Function
### Create
It is possible to specify ON or OFF for multiple dates at the same time.
* Validation
  * Dates must be set in ascending order
  	* Error ex： 1/1[ON] 1/3[OFF] 1/2[ON]
  * Validation of ON/OFF order is not performed.
  	* OK ex： 1/1[ON] 1/2[ON] 1/3[OFF]

### Update
Settings for future dates can be updated or deleted.
* Validation
  * Schedules that have already been completed cannot be updated
  	* ex: Setting was 1/1 0:00 [ON]  1/2 0:00 [OFF] 1/3 0:00 [ON] .\
     As of 1/2 1:00, updates are possible only on 1/3 0:00 ON

### Delete
It is possible to delete uncompleted schedules.
* Validation
  * Schedules that have already been completed cannot be deleted.

### Stop
It is possible to stop a schedule that has already started.
* Validation
  * Schedules that have already been completed cannot be stopped.

## Implementation

### Storage

* Change `auto_ops_rule` table definition.

```sql
      CREATE TABLE `auto_ops_rule` (
        `id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
        `feature_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
        `ops_type` int NOT NULL,
        `clauses` json NOT NULL,
        `status` int NOT NULL,                       # <-- add
#       `triggered_at` bigint NOT NULL,                <-- delete
        `created_at` bigint NOT NULL,
        `updated_at` bigint NOT NULL,
#       `deleted` tinyint(1) NOT NULL DEFAULT '0',     <-- delete
        `stopped_at` bigint NOT NULL,                # <-- add
        `environment_namespace` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
        PRIMARY KEY (`id`,`environment_namespace`),
        KEY `foreign_auto_ops_rule_feature_id_environment_namespace` (`feature_id`,`environment_namespace`),
        CONSTRAINT `foreign_auto_ops_rule_feature_id_environment_namespace` FOREIGN KEY (`feature_id`, `environment_namespace`) REFERENCES `feature` (`id`, `environment_namespace`)
      ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
```

Table explanation:\
`ops_type`: Currently it held `Enable` or `Disable`, but it has been changed to a column that manages the type of `AutoOperation` such as `Schedule` and `Event Rate`. Change the original information to be managed by `ActionType` in each `clauses`.\
`status`: Added management of schedule progress (waiting, running, completed, stopped, deleted).\
`triggered_at`: Delete as progress can be managed with `status`.\
`deleted`: Deleted to manage by `status`.\
`stopped_at`: Added to manage stopped date and time.

Migrations explanation:
1. Add "ActionType" to `clauses` Json data based on `ops_type`.
2. Change `ops_type` to `Schedule` or `Event Rate` based on the Json data of `clauses`.
3. Update `status` based on `triggered_at` and `deleted`.\
   `completed` if `triggered_at` is not empty, `deleted` if `deleted` is On, and `waiting` otherwise.
4. Delete `triggered_at` and `triggered_at`.

## Release Steps
1. Run the migration and modify any existing code related to it so that the functionality works correctly.
2. Add support related to Multi Schedule to the code.
