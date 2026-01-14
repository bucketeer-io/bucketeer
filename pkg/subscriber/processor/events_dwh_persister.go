// Copyright 2026 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller/codes"
	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber"
	storage "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/v2"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
)

// BigQuery specific configuration
type BigQueryConfig struct {
	Project  string `json:"project"`
	Dataset  string `json:"dataset"`
	Location string `json:"location"`
}

// MySQL specific configuration
type MySQLConfig struct {
	UseMainConnection bool   `json:"useMainConnection"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	User              string `json:"user"`
	Password          string `json:"password"`
	Database          string `json:"database"`
}

// Postgres specific configuration
type PostgresConfig struct {
	UseMainConnection bool   `json:"useMainConnection"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	User              string `json:"user"`
	Password          string `json:"password"`
	Database          string `json:"database"`
}

type eventsDWHPersisterConfig struct {
	// Persister-specific settings
	FlushInterval           int `json:"flushInterval"`
	FlushTimeout            int `json:"flushTimeout"`
	FlushSize               int `json:"flushSize"`
	MaxRetryGoalEventPeriod int `json:"maxRetryGoalEventPeriod,omitempty"`
	RetryGoalEventInterval  int `json:"retryGoalEventInterval,omitempty"`

	// Data warehouse configuration
	DataWarehouse DataWarehouseConfig `json:"dataWarehouse"`
}

type DataWarehouseConfig struct {
	Type      string `json:"type"` // bigquery, mysql
	BatchSize int    `json:"batchSize"`
	Timezone  string `json:"timezone"`

	// BigQuery specific configuration
	BigQuery BigQueryConfig `json:"bigquery"`
	// MySQL specific configuration
	MySQL MySQLConfig `json:"mysql"`
	// Postgres specific configuration
	Postgres PostgresConfig `json:"postgres"`
}

type eventsDWHPersister struct {
	eventsDWHPersisterConfig eventsDWHPersisterConfig
	mysqlClient              mysql.Client
	postgresClient           postgres.Client
	writer                   Writer
	subscriberType           string
	logger                   *zap.Logger
}

func NewEventsDWHPersister(
	ctx context.Context,
	config interface{},
	mysqlClient mysql.Client,
	postgresClient postgres.Client,
	redisClient redisv3.Client,
	persistentRedisClient redisv3.Client,
	exClient experimentclient.Client,
	ftClient featureclient.Client,
	persisterName string,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (subscriber.PubSubProcessor, error) {
	jsonConfig, ok := config.(map[string]interface{})
	if !ok {
		logger.Error("eventsDWHPersister: invalid config")
		return nil, ErrEventsDWHPersisterInvalidConfig
	}
	configBytes, err := json.Marshal(jsonConfig)
	if err != nil {
		logger.Error("eventsDWHPersister: failed to marshal config", zap.Error(err))
		return nil, err
	}
	var persisterConfig eventsDWHPersisterConfig
	err = json.Unmarshal(configBytes, &persisterConfig)
	if err != nil {
		logger.Error("eventsDWHPersister: failed to unmarshal config", zap.Error(err))
		return nil, err
	}

	// Validate configuration and set defaults
	if err := persisterConfig.validateAndSetDefaults(); err != nil {
		logger.Error("eventsDWHPersister: invalid configuration", zap.Error(err))
		return nil, err
	}

	// Determine the MySQL client to use for data warehouse operations
	var dwhMySQLClient mysql.Client
	if persisterConfig.DataWarehouse.Type == "mysql" {
		if persisterConfig.DataWarehouse.MySQL.UseMainConnection {
			// Use the existing MySQL client from main application
			dwhMySQLClient = mysqlClient
			logger.Info("Using main MySQL connection for data warehouse")
		} else {
			// Create a new MySQL client with separate connection
			dwhMySQLClient, err = createDedicatedMySQLClient(ctx, &persisterConfig.DataWarehouse.MySQL, logger)
			if err != nil {
				return nil, fmt.Errorf("failed to create dedicated MySQL client: %w", err)
			}
			logger.Info("Using dedicated MySQL connection for data warehouse",
				zap.String("host", persisterConfig.DataWarehouse.MySQL.Host),
				zap.String("database", persisterConfig.DataWarehouse.MySQL.Database),
			)
		}
	} else {
		// Default to main connection for non-MySQL types
		dwhMySQLClient = mysqlClient
	}

	// Determine the Postgres client to use for data warehouse operations
	var dwhPostgresClient postgres.Client
	if persisterConfig.DataWarehouse.Type == "postgres" {
		if persisterConfig.DataWarehouse.Postgres.UseMainConnection {
			// Use the existing Postgres client from main application
			dwhPostgresClient = postgresClient
			logger.Info("Using main Postgres connection for data warehouse")
		} else {
			// Create a new Postgres client with separate connection
			dwhPostgresClient, err = createDedicatedPostgresClient(ctx, &persisterConfig.DataWarehouse.Postgres, logger)
			if err != nil {
				return nil, fmt.Errorf("failed to create dedicated Postgres client: %w", err)
			}
			logger.Info("Using dedicated Postgres connection for data warehouse",
				zap.String("host", persisterConfig.DataWarehouse.Postgres.Host),
				zap.String("database", persisterConfig.DataWarehouse.Postgres.Database),
			)
		}
	} else {
		// Default to main connection for non-Postgres types
		dwhPostgresClient = postgresClient
	}

	e := &eventsDWHPersister{
		eventsDWHPersisterConfig: persisterConfig,
		mysqlClient:              dwhMySQLClient,
		postgresClient:           dwhPostgresClient,
		logger:                   logger,
	}
	experimentsCache := cachev3.NewExperimentsCache(cachev3.NewRedisCache(redisClient))
	location, err := locale.GetLocation(e.eventsDWHPersisterConfig.DataWarehouse.Timezone)
	if err != nil {
		return nil, err
	}

	switch persisterName {
	case EvaluationCountEventDWHPersisterName:
		e.subscriberType = subscriberEvaluationEventDWH

		// Create evaluation event writer based on data warehouse type
		var evalOptions []EvalEventWriterOption

		switch persisterConfig.DataWarehouse.Type {
		case "mysql":
			evalOptions = append(evalOptions, EvalEventWriterOption{
				DataWarehouseType: "mysql",
				MySQLClient:       dwhMySQLClient,
				BatchSize:         persisterConfig.DataWarehouse.BatchSize,
			})
		case "postgres":
			evalOptions = append(evalOptions, EvalEventWriterOption{
				DataWarehouseType: "postgres",
				PostgresClient:    dwhPostgresClient,
				BatchSize:         persisterConfig.DataWarehouse.BatchSize,
			})
		case "bigquery":
			// BigQuery is handled in the NewEvalEventWriter call below
		default:
			return nil, fmt.Errorf(
				"unsupported data warehouse type for evaluation events: %s",
				persisterConfig.DataWarehouse.Type,
			)
		}

		// Get BigQuery configuration
		project := persisterConfig.DataWarehouse.BigQuery.Project
		dataset := persisterConfig.DataWarehouse.BigQuery.Dataset
		bigQueryBatchSize := persisterConfig.DataWarehouse.BatchSize

		evalEventWriter, err := NewEvalEventWriter(
			ctx,
			logger,
			exClient,
			experimentsCache,
			project,
			dataset,
			bigQueryBatchSize,
			location,
			registerer,
			evalOptions...,
		)
		if err != nil {
			return nil, err
		}
		e.writer = evalEventWriter

	case GoalCountEventDWHPersisterName:
		e.subscriberType = subscriberGoalEventDWH

		// Create goal event writer based on data warehouse type
		var goalOptions []GoalEventWriterOption

		switch persisterConfig.DataWarehouse.Type {
		case "mysql":
			goalOptions = append(goalOptions, GoalEventWriterOption{
				DataWarehouseType: "mysql",
				MySQLClient:       dwhMySQLClient,
				BatchSize:         persisterConfig.DataWarehouse.BatchSize,
			})
		case "postgres":
			goalOptions = append(goalOptions, GoalEventWriterOption{
				DataWarehouseType: "postgres",
				PostgresClient:    dwhPostgresClient,
				BatchSize:         persisterConfig.DataWarehouse.BatchSize,
			})
		case "bigquery":
			// BigQuery is handled in the NewGoalEventWriter call below
		default:
			return nil, fmt.Errorf("unsupported data warehouse type for goal events: %s", persisterConfig.DataWarehouse.Type)
		}

		// Get BigQuery configuration
		project := persisterConfig.DataWarehouse.BigQuery.Project
		dataset := persisterConfig.DataWarehouse.BigQuery.Dataset
		location_str := persisterConfig.DataWarehouse.BigQuery.Location
		bigQueryBatchSize := persisterConfig.DataWarehouse.BatchSize

		// Get the max retry period and retry interval
		maxRetryPeriod := time.Duration(e.eventsDWHPersisterConfig.MaxRetryGoalEventPeriod) * time.Second
		retryInterval := time.Duration(e.eventsDWHPersisterConfig.RetryGoalEventInterval) * time.Second

		goalEventWriter, err := NewGoalEventWriter(
			ctx,
			logger,
			exClient,
			ftClient,
			experimentsCache,
			project,
			dataset,
			location_str,
			bigQueryBatchSize,
			location,
			persistentRedisClient,
			maxRetryPeriod,
			retryInterval,
			registerer,
			goalOptions...,
		)
		if err != nil {
			return nil, err
		}
		e.writer = goalEventWriter
	}
	return e, nil
}

func (e *eventsDWHPersister) Process(
	ctx context.Context,
	msgChan <-chan *puller.Message,
) error {
	batch := make(map[string]*puller.Message)
	ticker := time.NewTicker(time.Duration(e.eventsDWHPersisterConfig.FlushInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(e.subscriberType).Inc()
			id := msg.Attributes["id"]
			if id == "" {
				msg.Ack()
				// TODO: better log format for msg data
				subscriberHandledCounter.WithLabelValues(e.subscriberType, codes.MissingID.String()).Inc()
				continue
			}
			if previous, ok := batch[id]; ok {
				previous.Ack()
				subscriberHandledCounter.WithLabelValues(e.subscriberType, codes.DuplicateID.String()).Inc()
			}
			batch[id] = msg
			if len(batch) < e.eventsDWHPersisterConfig.FlushSize {
				continue
			}
			e.send(batch)
			batch = make(map[string]*puller.Message)
		case <-ticker.C:
			if len(batch) > 0 {
				e.send(batch)
				batch = make(map[string]*puller.Message)
			}
		case <-ctx.Done():
			batchSize := len(batch)
			e.logger.Debug("Context is done", zap.Int("batchSize", batchSize))
			if len(batch) > 0 {
				e.send(batch)
				e.logger.Debug(
					"All the left messages are processed successfully",
					zap.Int("batchSize", batchSize),
				)
			}
			return nil
		}
	}
}

func (e *eventsDWHPersister) Switch(ctx context.Context) (bool, error) {
	experimentStorage := storage.NewExperimentStorage(e.mysqlClient)
	count, err := experimentStorage.CountRunningExperiments(ctx)
	if err != nil {
		e.logger.Error("Failed to count running experiments", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

func (e *eventsDWHPersister) send(messages map[string]*puller.Message) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(e.eventsDWHPersisterConfig.FlushTimeout)*time.Second,
	)
	defer cancel()
	envEvents := e.extractEvents(messages)
	if len(envEvents) == 0 {
		e.logger.Error("all messages were bad")
		return
	}
	fails := e.writer.Write(ctx, envEvents)
	for id, m := range messages {
		if repeatable, ok := fails[id]; ok {
			if repeatable {
				m.Nack()
				subscriberHandledCounter.WithLabelValues(
					e.subscriberType,
					codes.RepeatableError.String(),
				).Inc()
			} else {
				m.Ack()
				subscriberHandledCounter.WithLabelValues(
					e.subscriberType,
					codes.NonRepeatableError.String(),
				).Inc()
			}
			continue
		}
		m.Ack()
		subscriberHandledCounter.WithLabelValues(e.subscriberType, codes.OK.String()).Inc()
	}
}
func (e *eventsDWHPersister) extractEvents(messages map[string]*puller.Message) environmentEventDWHMap {
	envEvents := environmentEventDWHMap{}
	handleBadMessage := func(m *puller.Message, err error) {
		m.Ack()
		e.logger.Error("Bad proto message",
			zap.Error(err),
			zap.String("messageID", m.ID),
			zap.ByteString("data", m.Data),
			zap.Any("attributes", m.Attributes),
		)
		subscriberHandledCounter.WithLabelValues(e.subscriberType, codes.BadMessage.String()).Inc()
	}
	for _, m := range messages {
		// Check if message data is empty
		if len(m.Data) == 0 {
			handleBadMessage(m, fmt.Errorf("message data is empty"))
			continue
		}
		event := &eventproto.Event{}
		if err := proto.Unmarshal(m.Data, event); err != nil {
			handleBadMessage(m, err)
			continue
		}
		var innerEvent ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(event.Event, &innerEvent); err != nil {
			handleBadMessage(m, err)
			continue
		}
		if innerEvents, ok := envEvents[event.EnvironmentId]; ok {
			innerEvents[event.Id] = innerEvent.Message
			continue
		}
		envEvents[event.EnvironmentId] = eventDWHMap{event.Id: innerEvent.Message}
	}
	return envEvents
}

// createDedicatedMySQLClient creates a new MySQL client with dedicated connection for data warehouse operations
func createDedicatedMySQLClient(ctx context.Context, config *MySQLConfig, logger *zap.Logger) (mysql.Client, error) {
	if config == nil {
		return nil, fmt.Errorf("mysql config is nil")
	}

	// Validate required fields
	if config.Host == "" || config.Database == "" || config.User == "" {
		return nil, fmt.Errorf("mysql host, database, and user are required for dedicated connection")
	}

	// Set default port if not specified
	port := config.Port
	if port == 0 {
		port = 3306 // Default MySQL port
	}

	// Create context with timeout for connection
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Create MySQL client with dedicated connection
	client, err := mysql.NewClient(
		ctx,
		config.User,
		config.Password,
		config.Host,
		port,
		config.Database,
		mysql.WithLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MySQL client: %w", err)
	}

	logger.Info("Created dedicated MySQL client for data warehouse",
		zap.String("host", config.Host),
		zap.Int("port", port),
		zap.String("database", config.Database),
		zap.String("user", config.User),
	)

	return client, nil
}

// createDedicatedPostgresClient creates a new Postgres client with dedicated connection for data warehouse operations
func createDedicatedPostgresClient(
	ctx context.Context,
	config *PostgresConfig,
	logger *zap.Logger,
) (postgres.Client, error) {
	if config == nil {
		return nil, fmt.Errorf("postgres config is nil")
	}

	// Validate required fields
	if config.Host == "" || config.Database == "" || config.User == "" {
		return nil, fmt.Errorf("postgres host, database, and user are required for dedicated connection")
	}

	// Set default port if not specified
	port := config.Port
	if port == 0 {
		port = 5432 // Default Postgres port
	}

	// Create context with timeout for connection
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Create Postgres client with dedicated connection
	client, err := postgres.NewClient(
		ctx,
		config.User,
		config.Password,
		config.Host,
		port,
		config.Database,
		postgres.WithLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Postgres client: %w", err)
	}

	logger.Info("Created dedicated Postgres client for data warehouse",
		zap.String("host", config.Host),
		zap.Int("port", port),
		zap.String("database", config.Database),
		zap.String("user", config.User),
	)

	return client, nil
}

// validateAndSetDefaults validates configuration and sets default values
func (config *eventsDWHPersisterConfig) validateAndSetDefaults() error {
	// Validate data warehouse type
	if config.DataWarehouse.Type == "" {
		return fmt.Errorf("dataWarehouse.type is required")
	}

	// Set default batch size if not specified
	if config.DataWarehouse.BatchSize == 0 {
		config.DataWarehouse.BatchSize = 1000 // default
	}

	// Set default timezone if not specified
	if config.DataWarehouse.Timezone == "" {
		config.DataWarehouse.Timezone = "UTC" // default
	}

	// Validate type-specific configuration
	switch config.DataWarehouse.Type {
	case "bigquery":
		if config.DataWarehouse.BigQuery.Project == "" || config.DataWarehouse.BigQuery.Dataset == "" {
			return fmt.Errorf("bigquery project and dataset are required")
		}
	case "mysql":
		// MySQL configuration is always present in the struct, no need to check for nil
	case "postgres":
		// Postgres configuration is always present in the struct, no need to check for nil
	default:
		return fmt.Errorf("unsupported data warehouse type: %s", config.DataWarehouse.Type)
	}

	return nil
}
