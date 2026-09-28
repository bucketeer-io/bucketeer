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
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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

func TestWithDB(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"default db", 0, 0},
		{"positive db index", 5, 5},
		{"negative db index is preserved for validation", -1, -1},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultOptions()
			WithDB(tt.input)(opts)
			assert.Equal(t, tt.expected, opts.db)
		})
	}
}

func TestNewClientWithDB(t *testing.T) {
	t.Parallel()

	t.Run("standalone mode applies db to underlying client", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithRedisMode(RedisModeStandalone),
			WithDB(3),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			stdClient, ok := rc.rc.(*goredis.Client)
			assert.True(t, ok)
			assert.Equal(t, 3, stdClient.Options().DB)
			c.Close()
		}
	})

	t.Run("cluster mode ignores db without error", func(t *testing.T) {
		logger := zap.NewNop()
		c, err := NewClient(
			"localhost:9999",
			WithLogger(logger),
			WithRedisMode(RedisModeCluster),
			WithDB(3),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c)

		if c != nil {
			rc := c.(*client)
			assert.Equal(t, ClientTypeCluster, rc.clientType)
			c.Close()
		}
	})

	t.Run("negative db fails instead of falling back to db 0", func(t *testing.T) {
		c, err := NewClient(
			"localhost:9999",
			WithLogger(zap.NewNop()),
			WithRedisMode(RedisModeStandalone),
			WithDB(-1),
		)
		assert.ErrorIs(t, err, ErrInvalidDB)
		assert.Nil(t, c)
	})
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

// testCA is an in-memory certificate authority used to issue server and
// client certificates for the TLS handshake tests.
type testCA struct {
	cert     *x509.Certificate
	key      *ecdsa.PrivateKey
	certPath string
}

func newTestCA(t *testing.T, dir string) *testCA {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "redis-test-ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)
	cert, err := x509.ParseCertificate(der)
	require.NoError(t, err)

	certPath := filepath.Join(dir, "ca.crt")
	require.NoError(t, os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600))
	return &testCA{cert: cert, key: key, certPath: certPath}
}

// issue signs a leaf certificate valid for localhost/127.0.0.1, writes the
// PEM cert/key under dir, and returns their paths.
func (ca *testCA) issue(
	t *testing.T,
	dir, prefix string,
	serial int64,
	usage x509.ExtKeyUsage,
) (certPath, keyPath string) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	template := &x509.Certificate{
		SerialNumber: big.NewInt(serial),
		Subject:      pkix.Name{CommonName: prefix},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{usage},
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, ca.cert, &key.PublicKey, ca.key)
	require.NoError(t, err)
	keyDER, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)

	certPath = filepath.Join(dir, prefix+".crt")
	keyPath = filepath.Join(dir, prefix+".key")
	require.NoError(t, os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600))
	require.NoError(t, os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}), 0o600))
	return certPath, keyPath
}

// startTLSRESPServer starts a minimal RESP server behind TLS that supports
// PING, SET and GET, and answers every other command (HELLO, CLIENT SETINFO,
// ...) with an error, which go-redis tolerates during its handshake. It
// returns the listening address.
func startTLSRESPServer(t *testing.T, tlsConf *tls.Config) string {
	t.Helper()

	ln, err := tls.Listen("tcp", "127.0.0.1:0", tlsConf)
	require.NoError(t, err)
	t.Cleanup(func() { ln.Close() })

	var mu sync.Mutex
	store := map[string]string{}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go serveRESP(conn, &mu, store)
		}
	}()
	return ln.Addr().String()
}

func serveRESP(conn net.Conn, mu *sync.Mutex, store map[string]string) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		args, err := readRESPCommand(r)
		if err != nil {
			// Includes TLS handshake failures, which surface on the first read.
			return
		}
		var reply string
		switch strings.ToUpper(args[0]) {
		case "PING":
			reply = "+PONG\r\n"
		case "SET":
			mu.Lock()
			store[args[1]] = args[2]
			mu.Unlock()
			reply = "+OK\r\n"
		case "GET":
			mu.Lock()
			v, ok := store[args[1]]
			mu.Unlock()
			if ok {
				reply = fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)
			} else {
				reply = "$-1\r\n"
			}
		default:
			reply = fmt.Sprintf("-ERR unknown command '%s'\r\n", args[0])
		}
		if _, err := io.WriteString(conn, reply); err != nil {
			return
		}
	}
}

func readRESPCommand(r *bufio.Reader) ([]string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(line, "*") {
		return nil, fmt.Errorf("unexpected RESP line: %q", line)
	}
	n, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil || n < 1 {
		return nil, fmt.Errorf("invalid RESP array length: %q", line)
	}
	args := make([]string, n)
	for i := range args {
		header, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		size, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(header, "$")))
		if err != nil {
			return nil, fmt.Errorf("invalid RESP bulk header: %q", header)
		}
		buf := make([]byte, size+2) // payload + CRLF
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		args[i] = string(buf[:size])
	}
	return args, nil
}

func TestNewClientTLSHandshake(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ca := newTestCA(t, dir)
	serverCertPath, serverKeyPath := ca.issue(t, dir, "server", 2, x509.ExtKeyUsageServerAuth)
	clientCertPath, clientKeyPath := ca.issue(t, dir, "client", 3, x509.ExtKeyUsageClientAuth)

	serverCert, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	require.NoError(t, err)
	clientCAs := x509.NewCertPool()
	clientCAs.AddCert(ca.cert)

	tlsAddr := startTLSRESPServer(t, &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{serverCert},
	})
	mtlsAddr := startTLSRESPServer(t, &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
	})

	patterns := []struct {
		desc      string
		addr      string
		tls       TLSConfig
		expectErr bool
	}{
		{
			desc: "succeeds when the server cert is verified with the configured CA",
			addr: tlsAddr,
			tls:  TLSConfig{Enabled: true, CACert: ca.certPath},
		},
		{
			desc:      "fails when the server cert is not signed by a trusted CA",
			addr:      tlsAddr,
			tls:       TLSConfig{Enabled: true},
			expectErr: true,
		},
		{
			desc: "succeeds without a trusted CA when verification is skipped",
			addr: tlsAddr,
			tls:  TLSConfig{Enabled: true, InsecureSkipVerify: true},
		},
		{
			desc:      "fails when TLS is disabled against a TLS-only server",
			addr:      tlsAddr,
			tls:       TLSConfig{},
			expectErr: true,
		},
		{
			desc: "mTLS succeeds with a client cert signed by the server's trusted CA",
			addr: mtlsAddr,
			tls: TLSConfig{
				Enabled: true,
				CACert:  ca.certPath,
				Cert:    clientCertPath,
				Key:     clientKeyPath,
			},
		},
		{
			desc:      "mTLS fails without a client cert",
			addr:      mtlsAddr,
			tls:       TLSConfig{Enabled: true, CACert: ca.certPath},
			expectErr: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			c, err := NewClient(
				p.addr,
				WithLogger(zap.NewNop()),
				WithRedisMode(RedisModeStandalone),
				WithDialTimeout(time.Second),
				WithTLS(p.tls),
			)
			require.NoError(t, err)
			defer c.Close()

			err = c.Set("key", "value", 0)
			if p.expectErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			val, err := c.Get("key")
			require.NoError(t, err)
			assert.Equal(t, []byte("value"), val)
		})
	}
}
