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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateDWHMySQLClientMissingConfig(t *testing.T) {
	t.Parallel()

	empty := ""
	zero := 0
	s := &server{
		mysqlUser:   &empty,
		mysqlPass:   &empty,
		mysqlHost:   &empty,
		mysqlPort:   &zero,
		mysqlDBName: &empty,
	}
	patterns := []struct {
		desc        string
		config      DataWarehouseMySQLConfig
		expectedErr string
	}{
		{
			desc:        "main connection without mysql flags",
			config:      DataWarehouseMySQLConfig{UseMainConnection: true},
			expectedErr: "data warehouse useMainConnection needs the mysql-user, mysql-host, and mysql-db-name flags",
		},
		{
			desc:        "dedicated connection without host",
			config:      DataWarehouseMySQLConfig{User: "user", Database: "db"},
			expectedErr: "mysql host, database, and user are required for data warehouse connection",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			client, err := s.createDWHMySQLClient(context.Background(), p.config, nil, zap.NewNop())
			require.Error(t, err)
			assert.Equal(t, p.expectedErr, err.Error())
			assert.Nil(t, client)
		})
	}
}
