## Overview

This rfc summarizes the result of comparision for each cloud platform.

## Current middlewares we use

| Currently used middlewares | Usage                                                 |
| -------------------------- | ----------------------------------------------------- |
| Cloud Pub/Sub              | Event-driven-architecture                             |
| Cloud SQL                  | Almost Util data such as feature, segmentUser, apiKey |
| BigQuery                   | A/B test                                              |
| Cloud KMS                  | Webhook                                               |
| MemoryStore                | Feature Flag such as auto ops and evaluation count    |

### Cloud Pub/Sub

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

**Comparison table about order guarantee**

|                              | Cloud Pub/Sub | Amazon Simple Notification Service（SNS）、Amazon Simple Queueing Service（SQS） | Azure Service Bus Messaging |
| ---------------------------- | ------------- | ------------------------------------------------------------ | --------------------------- |
| Order guarantee(optional)    | Yes           | Yes                                                          | Yes                         |
| Topic                        |               | 10,0000 -> 1,000                                             |                             |
| Subscription                 |               | 12,500,000 -> 100                                            |                             |
| API throttling(Tokyo region) |               | 1500 -> 300 transactions per second (3000 if high throughput) |                             |
| messages only once           |               | ◯                                                            |                             |

#### How to use SNS

```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
// snippet-start:[sns.go-v2.Publish]
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSPublishAPI defines the interface for the Publish function.
// We use this interface to test the function using a mocked service.
type SNSPublishAPI interface {
	Publish(ctx context.Context,
		params *sns.PublishInput,
		optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// PublishMessage publishes a message to an Amazon Simple Notification Service (Amazon SNS) topic
// Inputs:
//     c is the context of the method call, which includes the Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a PublishOutput object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to Publish
func PublishMessage(c context.Context, api SNSPublishAPI, input *sns.PublishInput) (*sns.PublishOutput, error) {
	return api.Publish(c, input)
}

func main() {
	msg := flag.String("m", "", "The message to send to the subscribed users of the topic")
	topicARN := flag.String("t", "", "The ARN of the topic to which the user subscribes")

	flag.Parse()

	if *msg == "" || *topicARN == "" {
		fmt.Println("You must supply a message and topic ARN")
		fmt.Println("-m MESSAGE -t TOPIC-ARN")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sns.NewFromConfig(cfg)

	msgs := []string{"1", "2", "3", "4", "5"}

	for _, m := range msgs {
		input := &sns.PublishInput{
			Message:  &m,
			TopicArn: topicARN,
		}

		_, err := PublishMessage(context.TODO(), client, input)
		if err != nil {
			fmt.Println("Got an error publishing the message:")
			fmt.Println(err)
			return
		}
	}

	// fmt.Println("Message ID: " + *result.MessageId)
}
```

#### How to use SQS

```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
// snippet-start:[sqs.go-v2.ReceiveMessage]
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// SQSReceiveMessageAPI defines the interface for the GetQueueUrl function.
// We use this interface to test the function using a mocked service.
type SQSReceiveMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

type SQSDeleteMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	DeleteMessage(ctx context.Context,
		params *sqs.DeleteMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

// GetQueueURL gets the URL of an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a GetQueueUrlOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to GetQueueUrl.
func GetQueueURL(c context.Context, api SQSReceiveMessageAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

// RemoveMessage deletes a message from an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a DeleteMessageOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to DeleteMessage.
func RemoveMessage(c context.Context, api SQSDeleteMessageAPI, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return api.DeleteMessage(c, input)
}

// GetMessages gets the most recent message from an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a ReceiveMessageOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ReceiveMessage.
func GetMessages(c context.Context, api SQSReceiveMessageAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}

func main() {
	queue := flag.String("q", "", "The name of the queue")
	timeout := flag.Int("t", 5, "How long, in seconds, that the message is hidden from others")
	flag.Parse()

	if *queue == "" {
		fmt.Println("You must supply the name of a queue (-q QUEUE)")
		return
	}

	if *timeout < 0 {
		*timeout = 0
	}

	if *timeout > 12*60*60 {
		*timeout = 12 * 60 * 60
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sqs.NewFromConfig(cfg)

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: queue,
	}

	// Get URL of queue
	urlResult, err := GetQueueURL(context.TODO(), client, gQInput)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	queueURL := urlResult.QueueUrl

	gMInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   int32(*timeout),
	}

	msgResult, err := GetMessages(context.TODO(), client, gMInput)
	if err != nil {
		fmt.Println("Got an error receiving messages:")
		fmt.Println(err)
		return
	}

	for _, m := range msgResult.Messages {
		fmt.Println("Message ID:     " + *m.MessageId)
		fmt.Println("Message Handle: " + *m.ReceiptHandle)
		fmt.Println("Message Handle: " + *m.Body)

		dMInput := &sqs.DeleteMessageInput{
			QueueUrl:      queueURL,
			ReceiptHandle: m.ReceiptHandle,
		}

		_, err = RemoveMessage(context.TODO(), client, dMInput)
		if err != nil {
			fmt.Println("Got an error deleting the message:")
			fmt.Println(err)
			return
		}
	}
}

// snippet-end:[sqs.go-v2.ReceiveMessage]

```

### Cloud SQL

#### Comparison Table

|      | GCP        | AWS                                                     | Azure                                                      |
| ---- | ---------- | ------------------------------------------------------- | ---------------------------------------------------------- |
|      | Cloud  SQL | Amazon Relational Database Service (RDS), Amazon Aurora | Azure Database for MySQL and Azure Database for PostgreSQL |

**Comparison table about AlloyDB**

|        | GCP     | AWS           | Azure              |
| ------ | ------- | ------------- | ------------------ |
|        | AlloyDB | Amazon Aurora | Azure SQL Database |
| SLA(%) | >=99.99 | >=99.99       | >=99.995           |


#### How to use RDS

##### Normal Way

```go
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var dbName string = "mydb"
	var dbUser string = "admin"
	var dbHost string = "xxxxxxxxxxxxxxxxxxxxxxxx.rds.amazonaws.com"
	var dbPort int = 3306
	var dbEndpoint string = fmt.Sprintf("%s:%d", dbHost, dbPort)
	var dbPass = "PASS"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?collation=utf8mb4_bin",
		dbUser, dbPass, dbEndpoint, dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	fmt.Println("hoge")

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	var sample Sample

	err = db.QueryRow("select * from sample where id = 'hogehoge';").Scan(&sample.id, &sample.num)

	if err != nil {
		panic(err)
	}

	fmt.Println(sample.id, sample.num)
}

type Sample struct {
	id  string
	num int
}
```

##### IAM authentication

```go
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Sample struct {
	id  string
	num int
}

func main() {

	var dbName string = "mydb"
	var dbUser string = "iam_user"
	var dbHost string = "xxxxxxxxxxxxxxxxxxxx.rds.amazonaws.com"
	var dbPort int = 3306
	var dbEndpoint string = fmt.Sprintf("%s:%d", dbHost, dbPort)
	var region string = "ap-northeast-1"
	var pemFile string = "xxxxxxxxxx.pem"

	caCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(pemFile)
	if err != nil {
		panic(err)
	}
	if ok := caCertPool.AppendCertsFromPEM(pem); !ok {
		panic("fail")
	}
	mysql.RegisterTLSConfig("rds", &tls.Config{
		ClientCAs:          caCertPool,
		InsecureSkipVerify: true,
	})

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error: " + err.Error())
	}

	authenticationToken, err := auth.BuildAuthToken(
		context.TODO(), dbEndpoint, region, dbUser, cfg.Credentials)
	if err != nil {
		panic("failed to create authentication token: " + err.Error())
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?allowCleartextPasswords=true&tls=rds",
		dbUser, authenticationToken, dbEndpoint, dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	fmt.Println("hoge")

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	var sample Sample

	err = db.QueryRow("select * from sample where id = 'hogehoge';").Scan(&sample.id, &sample.num)

	if err != nil {
		panic(err)
	}

	fmt.Println(sample.id, sample.num)
}

```

### BigQuery

#### Comparison Table

|      | GCP        | AWS1          | AWS2                | AWS3                       | Azure                   |
| ---- | ---------- | ------------- | ------------------- | -------------------------- | ----------------------- |
|      | BigQuery   | Amazon Athena | Amazon Redshift     | Amazon Redshift Spectrum   | Azure Synapse Analytics |
| SLA  | >= 99.99%  |  >= 99.9%     | >=99.9%(Multi Node) | ?                          | \>= 99.9%               |

### MemoryStore

#### Comparison Table

|        | GCP         | AWS                | Azure                  |
| ------ | ----------- | ------------------ | ---------------------- |
|        | Memorystore | Amazon ElastiCache | Azure Cache            |
| SLA(%) | >=99.9      | >=99.9             | >=99.9 (from Standard) |

#### How to use Amazon ElastiCache

```go
// Based on https://github.com/aws-samples/amazon-elasticache-samples
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

var ctx = context.Background()

type redisClient struct {
	redis *redis.Client
	mysql *sql.DB
}

func (c *redisClient) fetch(ctx context.Context, query string) (interface{}, error) {
	val, err := c.redis.Get(ctx, query).Bytes()
	if err != nil {
		return "", err
	}
	samples := []Sample{}
	err = json.Unmarshal(val, &samples)
	return samples, err
}

var dbName string = "tutorial"
var dbUser string = "admin"
var dbHost string = "xxxxxxxxxxxxxxxxxxx.rds.amazonaws.com"
var dbPort int = 3306
var dbEndpoint string = fmt.Sprintf("%s:%d", dbHost, dbPort)
var dbPass = "PASS"
var query = "SELECT * FROM planet"

func main() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?collation=utf8mb4_bin",
		dbUser, dbPass, dbEndpoint, dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "xxxxxxxxxxxxxxxxxxx.cache.amazonaws.com:6379",
	})

	client := redisClient{redis: rdb}
	a, err := client.redis.Ping(ctx).Result()
	fmt.Println(a)

	val, err := client.fetch(ctx, query)
	if err != nil {
		if err == redis.Nil {
			samples := []Sample{}
			var sample Sample

			rows, err := db.Query(query)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				if err := rows.Scan(&sample.Id, &sample.Name); err != nil {
					log.Fatal(err)
				}
				samples = append(samples, sample)
			}

			if err != nil {
				panic(err)
			}

			fmt.Println(samples)

			decoded, err := json.Marshal(samples)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", decoded)

			rdb.SetEX(ctx, query, string(decoded), 5*time.Second).Err()
		}
	}
	fmt.Println("key", val)
}

type Sample struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

```

#### Controversial topic
##### 1. Can we configure Memorystore as optional?

Yes. We can configure memory store as a optional in YAML file.

### Cloud KMS

#### Comparison Table

|        | GCP       | AWS                              | Azure           |
| ------ | --------- | -------------------------------- | --------------- |
|        | Cloud KMS | AWS Key Management Service (KMS) | Azure Key Vault |
| SLA(%) | >=99.99   | >=99.999                         | >=99.99         |


#### How to use AWS KMS

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

const keyID = "<KEY ID>"

func main() {
	ctx := context.TODO()
	str := `{"body":{"Alert id": 123}}`

	a, _ := NewAwsKMSCrypto(ctx, keyID, "ap-northeast-1")
	json := []byte(str)
	result, err := a.Encrypt(ctx, json)
	if err != nil {
		log.Fatal(err)
	}
	decripted, err := a.Decrypt(ctx, result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(decripted)) // => {"body":{"Alert id": 123}}
}

type EncrypterDecrypter interface {
	Encrypt(ctx context.Context, data []byte) ([]byte, error)
	Decrypt(ctx context.Context, data []byte) ([]byte, error)
}

type awsKMSCrypto struct {
	client *kms.Client
	keyID  string
}

func NewAwsKMSCrypto(
	ctx context.Context,
	keyID, region string,
) (EncrypterDecrypter, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	client := kms.NewFromConfig(cfg)
	return awsKMSCrypto{
		client: client,
		keyID:  keyID,
	}, nil
}

func (c awsKMSCrypto) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	resp, err := c.client.Encrypt(ctx, &kms.EncryptInput{
		Plaintext: data,
		KeyId:     &c.keyID,
	})
	if err != nil {
		return nil, err
	}
	return resp.CiphertextBlob, nil
}

func (c awsKMSCrypto) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	resp, err := c.client.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: data,
		KeyId:          &c.keyID,
	})
	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil
}

```

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



https://cloud.google.com/bigquery/sla

https://aws.amazon.com/athena/sla/

https://aws.amazon.com/redshift/sla/
