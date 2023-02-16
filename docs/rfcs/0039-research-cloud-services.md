## Overview

## Current middlewares we use

| Currently used middlewares | Usage                                                 |
| -------------------------- | ----------------------------------------------------- |
| Cloud Pub/Sub & Bigtable   | Event-driven-architecture                             |
| Cloud SQL & Memorystore    | Almost Util data such as feature, segmentUser, apiKey |
| Druid (GCS) & Kafka        | Calculator                                            |
| Cloud KMS                  | Webhook                                               |

### Cloud Pub/Sub & Bigtable

![event-pipeline](./images/0039-image1.png)

#### Comparison Table

|                                    | GCP                                                          | AWS                                                          | Azure                                                        |
| ---------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
|                                    | Cloud Pub/Sub                                                | Amazon Simple Notification Service（SNS）、Amazon Simple Queueing Service（SQS） | Azure Service Bus Messaging                                  |
| total size of Publish Request      | 10MB                                                         | 0.25MB                                                       | 0.25MB for [Standard tier](https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-premium-messaging)<br/>100 MB for [Premium tier](https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-premium-messaging). |
| number of messages per transaction | 1,000                                                        | 10                                                           | 100                                                          |
| Topic                              | 10,000                                                       | 10,0000                                                      | 10,000                                                       |
| Subscription per topic             | 10,000                                                       | 12,500,000                                                   | 2,000                                                        |
| Expiration period in topic         | 31 days                                                      | 4 weeks                                                      | 14 days                                                      |
| message retention period           | 7 days                                                       | 14 days (default is 4 days)                                  | 14 days                                                      |
| API throttling(Tokyo region)       |                                                              | 1500 transactions per second                                 |                                                              |
| At-Least-Once                      | ◯                                                            | ◯                                                            | ◯                                                            |
| Retry logic about subscription     | redelivers every subsequent message with the same ordering key, including acknowledged messages | redelivers only the messages                                 | redelivers only the messages                                 |
| SLA(%)                             | >=99.95                                                      | >=99.9                                                       | >=99.9                                                       |



|      | GCP      | AWS             | Azure           |
| ---- | -------- | --------------- | --------------- |
|      | Bigtable | Amazon DynamoDB | Azure Cosmos DB |
|      |          |                 |                 |

#### Controversial topic

##### 1. Can we stop using Bigtable?

We use Bigtable to store userEvaluation because PubSub doesn't guarantee order in the default. However, we can stop it if we use the order guarantee feature.

**Comparison table about order guarantee**

|                              | Cloud Pub/Sub | Amazon Simple Notification Service（SNS）、Amazon Simple Queueing Service（SQS） | Azure Service Bus Messaging |
| ---------------------------- | ------------- | ------------------------------------------------------------ | --------------------------- |
| Order guarantee(optional)    | Yes           | Yes                                                          | Yes                         |
| Topic                        |               | 10,0000 -> 1,000                                             |                             |
| Subscription                 |               | 12,500,000 -> 100                                            |                             |
| API throttling(Tokyo region) |               | 1500 -> 300 transactions per second (3000 if high throughput) |                             |
| messages only once           |               | ◯                                                            |                             |

**Important notice about enabling order guarantee**

* GCP
  * Possibility of increasing latency
  * Increase the number of duplicates
  * We have to design the ordering key so that granularity is not increased.
* AWS
  * By default, it limits up to 300 messages per second.
    * However, it can be configured to handle 6000 messages by enabling throughput mode.
* Azure
  * We have to design the session id so that granularity is not increased.

In conclusion, we can enable an order guarantee because we can set the key for each evaluation(userEvaluations + goalEvaluations). In short, we can design the key with fine granularity.

**Conversion plan from existing subscriber to ordering subscriber**

1. Create the new pubsub topic(bucketeer-xxx-evaluation-goal-events) in terraform.

2. Create the new subscription(bucketeer-xxx-evaluation-goal-events-event-persister) with turning on ordering feature and target topic.

3. Implement for enabling message ordering
  * Publisher
    * https://cloud.google.com/pubsub/docs/publisher#using-ordering-keys
  * Subscriber
    * https://cloud.google.com/pubsub/docs/ordering#enabling_message_ordering
  * Stores all events into dummy table such as dummy_evaluation_event, dummy_goal_event

4. Check if all messages are correctly stored into dummy tables correctly.

5. Move evaluations from BigTable into RDB.

6. Remove evaluation event persister and goal event persister.

7. Stop using Bigtable


### Cloud SQL & Memorystore

![event-pipeline](./images/0039-image3.png)

#### Comparison Table

|      | GCP        | AWS                                                     | Azure                                                      |
| ---- | ---------- | ------------------------------------------------------- | ---------------------------------------------------------- |
|      | Cloud  SQL | Amazon Relational Database Service (RDS), Amazon Aurora | Azure Database for MySQL and Azure Database for PostgreSQL |

|        | GCP         | AWS                | Azure                  |
| ------ | ----------- | ------------------ | ---------------------- |
|        | Memorystore | Amazon ElastiCache | Azure Cache            |
| SLA(%) | >=99.9      | >=99.9             | >=99.9 (from Standard) |

#### Controversial topic

##### 1. Can we configure Memorystore as optional?

Yes. We can configure memory store as a optional in YAML file.

### Druid (GCS) & Kafka

![data-pipeline](./images/0039-image4.png)

We use Druid as a relay DB and usual DB(fething data directory) and Kafka as an intereface for Druid. The amount of data is a huge and we have to handle them as a high performance when fething data directory.

#### Controversial topic

##### 1. Can we stop using Druid & Kafka?

We have the following problems when using Druid & Kafka.

* It's hard to maintainance self-hosted service.
* It's hard to solve the problem when something such as error occurs.

Therefore, We have the following options:

**1.  Replace with managed service**

|        | GCP     | AWS           | Azure              |
| ------ | ------- | ------------- | ------------------ |
|        | AlloyDB | Amazon Aurora | Azure SQL Database |
| SLA(%) | >=99.99 | >=99.99       | >=99.995           |

Pros

* If we can use AlloyDB, there is a possibility that we can use single DB.

Cons

* AlloyDB may not match with our use cases. For example, it takes time longer than Druid.

**2.  Preprocess evaluation data** 

![data-pipeline2](./images/0039-image5.png)

Pros

* This architecture can be consistent with Calculator.
* EventCounter service can access to only single DB.
* We don't need high performance DB.

Cons

* It takes time longer than 1 to finish this task.

##### Conclusion

We decided not to conclude this topic. Instead, we decided to store Evaluation Event to another DB, too. Therefore we can divide features into single feature flag service and A/B test service. 

Because Evaluation Event is a large data, we need to switch MySQL and PostgreSQL. Since PostgreSQL is ORDBMS, we need to define table as follows:

MySQL

```mysql
CREATE TABLE IF NOT EXISTS `feature` (
  `id` VARCHAR(255) NOT NULL,
  `tags` JSON NOT NULL,
  ...
  PRIMARY KEY (`id`, `environment_namespace`)
);
CREATE TABLE IF NOT EXISTS `tag` (
  `id` VARCHAR(255) NOT NULL,
  `created_at` BIGINT(20) NOT NULL,
  `updated_at` BIGINT(20) NOT NULL,
  `environment_namespace` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`, `environment_namespace`)
);
```

PostgreSQL

```postgresql
CREATE TABLE IF NOT EXISTS "feature" (
  "id" VARCHAR(255) NOT NULL,
  "tags" tags NOT NULL,
  ...
  PRIMARY KEY ("id", "environment_namespace")
);

INSERT INTO feature VALUES ('id', ROW('id', 1667370510, 1667370510, 'production'), ...);
...
```



### Cloud KMS

![webhook](./images/0039-image2.png)

#### Comparison Table

|        | GCP       | AWS                              | Azure           |
| ------ | --------- | -------------------------------- | --------------- |
|        | Cloud KMS | AWS Key Management Service (KMS) | Azure Key Vault |
| SLA(%) | >=99.99   | >=99.999                         | >=99.99         |

## Ref

Order guarantee

https://cloud.google.com/pubsub/docs/ordering#console

https://medium.com/google-cloud/google-cloud-pub-sub-ordered-delivery-1e4181f60bc8

https://aws.amazon.com/jp/about-aws/whats-new/2022/10/amazon-sqs-increased-throughput-quota-fifo-high-throughput-ht-mode-6000-transactions-per-second-tps/

https://learn.microsoft.com/en-us/azure/service-bus-messaging/message-sessions

https://devblogs.microsoft.com/premier-developer/ordering-messages-in-azure-service-bus/

https://christina04.hatenablog.com/entry/gcp-cloud-pubsub-ordering-key-concern



Quotas

https://cloud.google.com/pubsub/quotas

https://docs.aws.amazon.com/general/latest/gr/sns.html

https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/quotas-messages.html

https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-quotas

SLA

https://aws.amazon.com/messaging/sla/

https://cloud.google.com/pubsub/sla

https://azure.microsoft.com/ja-jp/support/legal/sla/service-bus/v1_1/



https://azure.microsoft.com/ja-jp/support/legal/sla/cache/v1_1/

https://aws.amazon.com/elasticache/sla/

https://cloud.google.com/memorystore/sla



https://aws.amazon.com/kms/sla/

https://cloud.google.com/kms/sla

https://azure.microsoft.com/ja-jp/updates/akv-sla-raised-to-9999/



https://azure.microsoft.com/ja-jp/support/legal/sla/azure-sql-database/v1_8/



https://aws.amazon.com/rds/aurora/sla/

https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues-understanding-logic.html
