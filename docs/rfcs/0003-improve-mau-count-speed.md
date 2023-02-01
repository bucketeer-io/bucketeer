# Summary

We'll migrate the user event data from Druid to another DB to achieve the following goal.

- Retrieve the MAU count faster. Currently, it takes more than 1 hour depending on the number of events stored in the Druid, and sometimes it fails during the querying

## Proposal

I propose saving the user id in one row instead of all events and updating the count column to know how many events we received.
So we can retrieve the count much faster even if we use a relational database.<br />
Also, it reduces the Druid instance and storage costs.

### Table

```sql
CREATE TABLE IF NOT EXISTS `mau` (
  `user_id` VARCHAR(255) NOT NULL,
  `yearmonth` VARCHAR(6) NOT NULL,
  `source_id` VARCHAR(30) NOT NULL,
  `event_count` INT(11) UNSIGNED NOT NULL,
  `created_at` BIGINT(20) NOT NULL,
  `updated_at` BIGINT(20) NOT NULL,
  `environment_namespace` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`environment_namespace`, `yearmonth`, `source_id`, `user_id`)
)
PARTITION BY RANGE COLUMNS(`yearmonth`) (
  PARTITION p202211 VALUES LESS THAN ('202212'),
  PARTITION p202212 VALUES LESS THAN ('202301'),
  PARTITION p202301 VALUES LESS THAN ('202302')
);
```

**Note:** The partitions will be used to delete large data faster without locking the table during the process.

### Upsert

We will increment the `event_count` column if the user exists.

```sql
INSERT INTO mau (
  user_id,
  yearmonth,
  source_id,
  event_count,
  created_at,
  updated_at,
  environment_namespace
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
) ON DUPLICATE KEY UPDATE
  event_count = event_count + 1
  updated_at = VALUES(1668737826),
```

### Retrieve

```sql
SELECT
  COUNT(*) as user_count,
  IFNULL(SUM(event_count), 0) as event_count
FROM
  mau
where
  environment_namespace = 'namespace' AND
  yearmonth = '202212'
```

## Server changes

### Event Persister

We will add the upsert implementation in the `event-persister-user-event`.

**Note:** We will delete the Kafka implementation once we have gathered the entire month's events.

### Event Counter

We will add a new API called `GetMAUCount` in the `event-counter` service to retrieve the count from the new DB instead of Druid, using the same response format.<br />
We also need to make these changes in the [service.proto](https://github.com/bucketeer-io/bucketeer/blob/main/proto/eventcounter/service.proto#L155) file.

**Note:** We will delete the old `GetUserCountV2` API once we have gathered the entire month's events.

### Notification Sender

We will change the API name from `GetUserCountV2` to `GetMAUCount` in the `notification-sender` service.

### Send MAU script

We will change the API name from `GetUserCountV2` to `GetMAUCount`.

## Data Deletion

We will add a cronjob in Kubernetes which will send a request to the batch service API and delete the partition. It will also create a new partition if needed.

## Backup

CloudSQL automatically backs up the database once a day, so we don't need to do anything.<br />
But in case we need to back up more than once, we will need to add a cronjob. For this proposal once a day is enough.

## Migration

Because there is no need to rush, I'm going to implement it to double-write the data for 30 days and then delete the Kafka implementation after we confirm everything is okay.

No need to stop event persister services during this period.
