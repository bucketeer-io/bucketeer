# Summary

We'll migrate the evaluation count data from Druid to another place to achieve the following goals.

- The time series data retrieval is done in milliseconds instead of 5-20 seconds
- Split the data pipeline for Feature Flag and A/B Test, making the self-host process simple and cheaper
- Be able to see the default value count

## Proposal

### Redis

* Pros
  * Implementation cost (We already use Redis)
  * No need for schema management
  * Data retrieval can be done directly from Redis in milliseconds, and no need to summarize the count
  * Doesn't need much space to store the counts
  * No need to create a service and a table to summarize the daily count
  * No need to create a scheduled workflow to delete old data
  * Support up to 60k connections
  * Cluster cost is cheap

* Cons
  * Need to create a scheduled workflow to backup the current data
  * There is an error rate of up to 0.81% for the unique count

See also: [pricing](https://cloud.google.com/memorystore/docs/redis/pricing)

### PostgreSQL

* Pros
  * The unique count will be accurate
  * Backup can be done automatically

* Cons
  * It requires much more space to store the data (Terabytes)
  * Need schema management
  * Need to create a service to summarize the daily count and save it in another table to speed up the data retrieval
  * Need to create a scheduled workflow to delete old data
  * Need to create a scheduled workflow to create new partitions daily to make it easier when we delete old data
  * Implementation cost
  * Max connection is limited and requires adjustments
  * Instance cost is not cheap due to the high volume of requests

See also: [pricing](https://cloud.google.com/sql/pricing)

## Implementation

### Infra

We may need to increase the current Redis storage via Terraform. Currently, we use 1GB.

### Server

#### Event Persister

We will implement the Redis using `INCR` interface to increment the event counter
For the user count, we will use the `PFADD` (HyperLogLog) interface to increment the unique counter.

Also, we will use the `EXPIRE` to set a TTL of 31 days so that the keys will delete automatically.

##### Key format

Event count: `ec:feature_flag_id:variation_id:daily_timestamp`
User count: `uc:feature_flag_id:variation_id:daily_timestamp`

**Note:** For default evaluation events, we set the variation id as `default`.

#### Event Counter Storage

We will change the event counter API's storage interface to retrieve the data from Redis instead of Druid and convert the data to the current Timeseries format.
No changes are needed in the console UI.

##### Event count

We will get multiple counters using the `MGET` interface.

##### User count

We will get the unique count using the `PFCOUNT` interface.

## Migration

Because there is no need to rush, I'm going to implement it to double-write the data for 30 days and then delete the old implementation after we confirm everything is okay.

No need to stop event persister services during this period.
