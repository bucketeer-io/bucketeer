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

package v3

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// writeTestCertKeyPair generates a self-signed EC certificate/key pair and
// writes them as PEM files under dir, returning their paths.
func writeTestCertKeyPair(t *testing.T, dir, prefix string) (certPath, keyPath string) {
	t.Helper()

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "redis-test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:         true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	require.NoError(t, err)

	certPath = filepath.Join(dir, prefix+".crt")
	keyPath = filepath.Join(dir, prefix+".key")

	certOut, err := os.Create(certPath)
	require.NoError(t, err)
	defer certOut.Close()
	require.NoError(t, pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}))

	keyBytes, err := x509.MarshalECPrivateKey(priv)
	require.NoError(t, err)
	keyOut, err := os.Create(keyPath)
	require.NoError(t, err)
	defer keyOut.Close()
	require.NoError(t, pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes}))

	return certPath, keyPath
}

func TestNewClientIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name          string
		addr          string
		expectError   bool
		expectCluster bool
	}{
		{
			name:          "standalone redis on default port",
			addr:          "localhost:6379",
			expectError:   false,
			expectCluster: false,
		},
		{
			name:          "unreachable redis",
			addr:          "localhost:9999",
			expectError:   false,
			expectCluster: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			logger := zap.NewNop()
			client, err := NewClient(tt.addr, WithLogger(logger))

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)

				if client != nil {
					client.Close()
				}
			}
		})
	}
}

func TestNewClientBehavior(t *testing.T) {
	t.Parallel()

	t.Run("unreachable redis returns client with auto mode", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient("localhost:9999", WithLogger(logger))

		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			assert.Equal(t, ClientTypeStandard, rc.clientType)
			c.Close()
		}
	})

	t.Run("options are applied", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithPoolSize(20),
			WithMinIdleConns(5),
			WithPassword("test-password"),
		)

		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			c.Close()
		}
	})
}

func TestNewClientWithRedisMode(t *testing.T) {
	t.Parallel()

	t.Run("cluster mode creates ClusterClient", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithRedisMode(RedisModeCluster),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			assert.Equal(t, ClientTypeCluster, rc.clientType)
			_, ok := rc.rc.(*goredis.ClusterClient)
			assert.True(t, ok)
			c.Close()
		}
	})

	t.Run("standalone mode creates standard Client", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithRedisMode(RedisModeStandalone),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			assert.Equal(t, ClientTypeStandard, rc.clientType)
			_, ok := rc.rc.(*goredis.Client)
			assert.True(t, ok)
			c.Close()
		}
	})

	t.Run("auto mode defaults to standalone when unreachable", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithRedisMode(RedisModeAuto),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			assert.Equal(t, ClientTypeStandard, rc.clientType)
			c.Close()
		}
	})

	t.Run("invalid mode falls back to auto", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithRedisMode("invalid"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			// auto mode defaults to standalone when Redis is unreachable
			assert.Equal(t, ClientTypeStandard, rc.clientType)
			c.Close()
		}
	})

	t.Run("case-insensitive mode parsing", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithRedisMode("CLUSTER"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			assert.Equal(t, ClientTypeCluster, rc.clientType)
			c.Close()
		}
	})
}

func TestWithRedisMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    RedisMode
		expected RedisMode
	}{
		{"cluster", RedisModeCluster, RedisModeCluster},
		{"standalone", RedisModeStandalone, RedisModeStandalone},
		{"auto", RedisModeAuto, RedisModeAuto},
		{"uppercase CLUSTER", "CLUSTER", RedisModeCluster},
		{"mixed case Standalone", "Standalone", RedisModeStandalone},
		{"invalid falls back to auto", "invalid", RedisModeAuto},
		{"empty falls back to auto", "", RedisModeAuto},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultOptions()
			WithRedisMode(tt.input)(opts)
			assert.Equal(t, tt.expected, opts.redisMode)
		})
	}
}

func TestClientTypeString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "cluster", clientTypeString(ClientTypeCluster))
	assert.Equal(t, "standalone", clientTypeString(ClientTypeStandard))
}

func TestWithTLS(t *testing.T) {
	t.Parallel()

	cfg := TLSConfig{Enabled: true, CACert: "/path/to/ca.crt"}
	opts := defaultOptions()
	WithTLS(cfg)(opts)
	assert.Equal(t, cfg, opts.tls)
}

func TestBuildTLSConfig(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	caCertPath, _ := writeTestCertKeyPair(t, dir, "ca")
	certPath, keyPath := writeTestCertKeyPair(t, dir, "client")

	t.Run("disabled returns nil config", func(t *testing.T) {
		t.Parallel()
		tlsConfig, err := buildTLSConfig(TLSConfig{Enabled: false})
		require.NoError(t, err)
		assert.Nil(t, tlsConfig)
	})

	t.Run("enabled with no cert paths uses system pool", func(t *testing.T) {
		t.Parallel()
		tlsConfig, err := buildTLSConfig(TLSConfig{Enabled: true})
		require.NoError(t, err)
		require.NotNil(t, tlsConfig)
		assert.Nil(t, tlsConfig.RootCAs)
		assert.False(t, tlsConfig.InsecureSkipVerify)
	})

	t.Run("insecure skip verify is propagated", func(t *testing.T) {
		t.Parallel()
		tlsConfig, err := buildTLSConfig(TLSConfig{Enabled: true, InsecureSkipVerify: true})
		require.NoError(t, err)
		require.NotNil(t, tlsConfig)
		assert.True(t, tlsConfig.InsecureSkipVerify)
	})

	t.Run("valid CA cert is loaded", func(t *testing.T) {
		t.Parallel()
		tlsConfig, err := buildTLSConfig(TLSConfig{Enabled: true, CACert: caCertPath})
		require.NoError(t, err)
		require.NotNil(t, tlsConfig)
		assert.NotNil(t, tlsConfig.RootCAs)
	})

	t.Run("missing CA cert file errors", func(t *testing.T) {
		t.Parallel()
		_, err := buildTLSConfig(TLSConfig{Enabled: true, CACert: "/nonexistent/ca.crt"})
		assert.Error(t, err)
	})

	t.Run("invalid CA cert content errors", func(t *testing.T) {
		t.Parallel()
		badCACert := filepath.Join(dir, "bad-ca.crt")
		require.NoError(t, os.WriteFile(badCACert, []byte("not a pem cert"), 0o600))
		_, err := buildTLSConfig(TLSConfig{Enabled: true, CACert: badCACert})
		assert.Error(t, err)
	})

	t.Run("valid client cert and key are loaded", func(t *testing.T) {
		t.Parallel()
		tlsConfig, err := buildTLSConfig(TLSConfig{Enabled: true, Cert: certPath, Key: keyPath})
		require.NoError(t, err)
		require.NotNil(t, tlsConfig)
		assert.Len(t, tlsConfig.Certificates, 1)
	})

	t.Run("cert without key errors", func(t *testing.T) {
		t.Parallel()
		_, err := buildTLSConfig(TLSConfig{Enabled: true, Cert: certPath})
		assert.Error(t, err)
	})

	t.Run("key without cert errors", func(t *testing.T) {
		t.Parallel()
		_, err := buildTLSConfig(TLSConfig{Enabled: true, Key: keyPath})
		assert.Error(t, err)
	})

	t.Run("mismatched cert and key errors", func(t *testing.T) {
		t.Parallel()
		_, otherKeyPath := writeTestCertKeyPair(t, dir, "other")
		_, err := buildTLSConfig(TLSConfig{Enabled: true, Cert: certPath, Key: otherKeyPath})
		assert.Error(t, err)
	})
}

func TestNewClientWithTLS(t *testing.T) {
	t.Parallel()

	t.Run("TLS enabled against unreachable host does not fail startup", func(t *testing.T) {
		t.Parallel()
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithTLS(TLSConfig{Enabled: true}),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		c.Close()
	})

	t.Run("invalid TLS config returns error", func(t *testing.T) {
		t.Parallel()
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithTLS(TLSConfig{Enabled: true, CACert: "/nonexistent/ca.crt"}),
		)
		assert.Error(t, err)
		assert.Nil(t, c)
	})
}
