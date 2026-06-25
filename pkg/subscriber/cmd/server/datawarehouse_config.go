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

package server

import (
	"encoding/json"
	"fmt"
)

// DataWarehouseConfig is the data-warehouse section of the events DWH persister config
// (onDemandProcessors.json). The server reads it to initialize the data-warehouse storages,
// which are separate from the operational database.
type DataWarehouseConfig struct {
	Type      string `json:"type"` // bigquery, mysql, postgres
	BatchSize int    `json:"batchSize"`
	Timezone  string `json:"timezone"`

	BigQuery BigQueryConfig `json:"bigquery"`
	MySQL    MySQLConfig    `json:"mysql"`
	Postgres PostgresConfig `json:"postgres"`
}

// BigQueryConfig is the BigQuery-specific data-warehouse configuration.
type BigQueryConfig struct {
	Project  string `json:"project"`
	Dataset  string `json:"dataset"`
	Location string `json:"location"`
}

// MySQLConfig is the MySQL-specific data-warehouse configuration.
type MySQLConfig struct {
	UseMainConnection bool   `json:"useMainConnection"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	User              string `json:"user"`
	Password          string `json:"password"`
	Database          string `json:"database"`
}

// PostgresConfig is the Postgres-specific data-warehouse configuration.
type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// parseDWHConfig extracts and validates the data-warehouse configuration from a raw
// persister config (as loaded from onDemandProcessors.json).
func parseDWHConfig(config interface{}) (DataWarehouseConfig, error) {
	jsonConfig, ok := config.(map[string]interface{})
	if !ok {
		return DataWarehouseConfig{}, fmt.Errorf("data warehouse: invalid persister config")
	}
	configBytes, err := json.Marshal(jsonConfig)
	if err != nil {
		return DataWarehouseConfig{}, err
	}
	var wrapper struct {
		DataWarehouse DataWarehouseConfig `json:"dataWarehouse"`
	}
	if err := json.Unmarshal(configBytes, &wrapper); err != nil {
		return DataWarehouseConfig{}, err
	}
	dwh := wrapper.DataWarehouse
	if err := dwh.validateAndSetDefaults(); err != nil {
		return DataWarehouseConfig{}, err
	}
	return dwh, nil
}

// validateAndSetDefaults validates the data-warehouse configuration and sets default values.
func (c *DataWarehouseConfig) validateAndSetDefaults() error {
	if c.Type == "" {
		return fmt.Errorf("dataWarehouse.type is required")
	}
	if c.BatchSize == 0 {
		c.BatchSize = 1000 // default
	}
	if c.Timezone == "" {
		c.Timezone = "UTC" // default
	}
	switch c.Type {
	case "bigquery":
		if c.BigQuery.Project == "" || c.BigQuery.Dataset == "" {
			return fmt.Errorf("bigquery project and dataset are required")
		}
	case "mysql":
		// MySQL configuration is always present in the struct.
	case "postgres":
		// Postgres configuration is always present in the struct.
	default:
		return fmt.Errorf("unsupported data warehouse type: %s", c.Type)
	}
	return nil
}
