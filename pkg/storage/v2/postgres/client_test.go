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

package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDSN(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		ssl      SSLConfig
		expected string
		isErr    bool
	}{
		{
			desc:     "success: ssl disabled",
			ssl:      SSLConfig{Mode: SSLModeDisable},
			expected: "postgres://bucketeer:password@localhost:5432/bucketeer?sslmode=disable",
		},
		{
			desc:     "success: ssl required",
			ssl:      SSLConfig{Mode: "require"},
			expected: "postgres://bucketeer:password@localhost:5432/bucketeer?sslmode=require",
		},
		{
			desc: "success: verify-full with a root certificate",
			ssl: SSLConfig{
				Mode:     "verify-full",
				RootCert: "/usr/local/certs/postgres/ca.crt",
			},
			expected: "postgres://bucketeer:password@localhost:5432/bucketeer" +
				"?sslmode=verify-full&sslrootcert=%2Fusr%2Flocal%2Fcerts%2Fpostgres%2Fca.crt",
		},
		{
			desc: "success: verify-full with a client certificate",
			ssl: SSLConfig{
				Mode:     "verify-full",
				RootCert: "/certs/ca.crt",
				Cert:     "/certs/tls.crt",
				Key:      "/certs/tls.key",
			},
			expected: "postgres://bucketeer:password@localhost:5432/bucketeer" +
				"?sslcert=%2Fcerts%2Ftls.crt&sslkey=%2Fcerts%2Ftls.key&sslmode=verify-full" +
				"&sslrootcert=%2Fcerts%2Fca.crt",
		},
		{
			desc:  "err: unknown ssl mode",
			ssl:   SSLConfig{Mode: "enable"},
			isErr: true,
		},
		{
			desc:  "err: empty ssl mode",
			ssl:   SSLConfig{},
			isErr: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			dsn, err := buildDSN("bucketeer", "password", "localhost", 5432, "bucketeer", p.ssl)
			if p.isErr {
				assert.Error(t, err)
				assert.Empty(t, dsn)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, p.expected, dsn)
		})
	}
}

func TestBuildDSNEscapesCredentials(t *testing.T) {
	t.Parallel()
	dsn, err := buildDSN("bucketeer@corp", "p@ss:w/rd?", "db.example.com", 5432, "bucketeer", SSLConfig{
		Mode: "require",
	})
	assert.NoError(t, err)
	assert.Equal(
		t,
		"postgres://bucketeer%40corp:p%40ss%3Aw%2Frd%3F@db.example.com:5432/bucketeer?sslmode=require",
		dsn,
	)
}

func TestSSLDefaultsToRequire(t *testing.T) {
	t.Parallel()
	assert.Equal(t, SSLModeRequire, defaultOptions().ssl.Mode)

	opts := defaultOptions()
	WithSSL(SSLConfig{RootCert: "/certs/ca.crt"})(opts)
	assert.Equal(t, SSLModeRequire, opts.ssl.Mode)
	assert.Equal(t, "/certs/ca.crt", opts.ssl.RootCert)

	opts = defaultOptions()
	WithSSL(SSLConfig{Mode: SSLModeDisable})(opts)
	assert.Equal(t, SSLModeDisable, opts.ssl.Mode)
}
