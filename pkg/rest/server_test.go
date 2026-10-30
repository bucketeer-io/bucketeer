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

package rest

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	certPath = "testdata/server.crt"
	keyPath  = "testdata/server.key"
	timeout  = 10 * time.Second
	port     = 9222
	// Ports of the servers started by individual tests with custom timeouts.
	slowHeadersPort = 9223
	streamingPort   = 9224
)

type dummyService struct{}

func (*dummyService) Register(mux *http.ServeMux) {}

// handlerService registers a single handler, so tests can serve custom responses.
type handlerService struct {
	path    string
	handler http.HandlerFunc
}

func (s *handlerService) Register(mux *http.ServeMux) {
	mux.HandleFunc(s.path, s.handler)
}

func newServer() *Server {
	logger := zap.NewExample()
	server := NewServer(
		certPath,
		keyPath,
		WithLogger(logger),
		WithPort(port),
		WithService(&dummyService{}),
	)
	return server
}

func newHTTPClient() *http.Client {
	/*
	 ** Better to do it like net/http/httptest:
	 ** https://golang.org/pkg/crypto/tls/#example_Dial
	 ** certpool := x509.NewCertPool()
	 ** certpool.AddCert(s.certificate)
	 ** s.client.Transport = &http.Transport{
	 **   TLSClientConfig: &tls.Config{
	 **     RootCAs: certpool,
	 ** }}
	 */
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 5 * time.Second,
	}
}

func TestMain(m *testing.M) {
	// Because os.Exit doesn't return, we need to call defer in separated function.
	code := testMain(m)
	os.Exit(code)
}

func testMain(m *testing.M) int {
	server := newServer()
	defer server.Stop(time.Second)
	go server.Run()
	waitForServer()
	return m.Run()
}

func waitForServer() {
	client := newHTTPClient()
	// Give up if the server does not come up in time, but do not keep the
	// watchdog running once it did: it would kill the whole test binary.
	watchdog := time.AfterFunc(timeout, func() {
		fmt.Fprintln(os.Stderr, "failed to get response")
		os.Exit(1)
	})
	defer watchdog.Stop()
	for {
		resp, err := client.Get(fmt.Sprintf("https://localhost:%d", port))
		if err == nil {
			resp.Body.Close()
			return
		}
	}
}

// startServer runs a server on its own port for the duration of the test.
func startServer(t *testing.T, serverPort int, opts ...Option) {
	t.Helper()
	opts = append([]Option{WithLogger(zap.NewNop()), WithPort(serverPort)}, opts...)
	server := NewServer(certPath, keyPath, opts...)
	go server.Run()
	client := newHTTPClient()
	require.Eventually(t, func() bool {
		resp, err := client.Get(fmt.Sprintf("https://localhost:%d", serverPort))
		if err != nil {
			return false
		}
		resp.Body.Close()
		return true
	}, timeout, 10*time.Millisecond)
	t.Cleanup(func() { server.Stop(time.Second) })
}

func TestNewServerDefaultTimeouts(t *testing.T) {
	t.Parallel()
	server := NewServer(certPath, keyPath, WithService(&dummyService{}))
	server.setup()

	assert.Equal(t, 10*time.Second, server.httpServer.ReadHeaderTimeout)
	assert.Equal(t, 30*time.Second, server.httpServer.ReadTimeout)
	assert.Equal(t, 30*time.Second, server.httpServer.WriteTimeout)
	assert.Equal(t, 60*time.Second, server.httpServer.IdleTimeout)
}

func TestWithTimeouts(t *testing.T) {
	t.Parallel()
	server := NewServer(
		certPath, keyPath,
		WithService(&dummyService{}),
		WithTimeouts(time.Second, 2*time.Second, 3*time.Second, 4*time.Second),
	)
	server.setup()

	assert.Equal(t, time.Second, server.httpServer.ReadHeaderTimeout)
	assert.Equal(t, 2*time.Second, server.httpServer.ReadTimeout)
	assert.Equal(t, 3*time.Second, server.httpServer.WriteTimeout)
	assert.Equal(t, 4*time.Second, server.httpServer.IdleTimeout)
}

// A connection whose request headers never complete must be closed by the
// server once ReadHeaderTimeout elapses, while complete requests on the same
// server keep being served.
func TestServerReadHeaderTimeout(t *testing.T) {
	const readHeaderTimeout = 200 * time.Millisecond
	startServer(t, slowHeadersPort,
		WithTimeouts(readHeaderTimeout, time.Second, time.Second, time.Second),
		WithService(&handlerService{
			path:    "/ok",
			handler: func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) },
		}),
	)

	t.Run("incomplete request headers: connection is closed", func(t *testing.T) {
		conn, err := tls.Dial(
			"tcp", fmt.Sprintf("localhost:%d", slowHeadersPort),
			&tls.Config{InsecureSkipVerify: true}, // self-signed test certificate
		)
		require.NoError(t, err)
		defer conn.Close()
		_, err = fmt.Fprintf(conn, "GET /ok HTTP/1.1\r\nHost: localhost\r\nX-Incomplete: ")
		require.NoError(t, err)

		start := time.Now()
		require.NoError(t, conn.SetReadDeadline(time.Now().Add(timeout)))
		_, err = conn.Read(make([]byte, 1))

		require.Error(t, err, "the server kept the connection open")
		var netErr net.Error
		assert.False(t, errors.As(err, &netErr) && netErr.Timeout(), "the server kept the connection open: %v", err)
		assert.Less(t, time.Since(start), timeout/2)
	})

	t.Run("complete request: served", func(t *testing.T) {
		resp, err := newHTTPClient().Get(fmt.Sprintf("https://localhost:%d/ok", slowHeadersPort))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// A response that is still being written when WriteTimeout elapses is cut
// off. Handlers that stream for longer (SSE) clear the deadline of their own
// response through http.ResponseController, which must work through the
// middleware chain of the server.
func TestServerWriteTimeoutStreamingResponse(t *testing.T) {
	const (
		writeTimeout  = 300 * time.Millisecond
		chunks        = 6
		chunkInterval = 100 * time.Millisecond // the whole stream outlasts writeTimeout
	)
	stream := func(clearWriteDeadline bool) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if clearWriteDeadline {
				if err := http.NewResponseController(w).SetWriteDeadline(time.Time{}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			flusher := w.(http.Flusher)
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			flusher.Flush()
			for i := 0; i < chunks; i++ {
				time.Sleep(chunkInterval)
				if _, err := fmt.Fprintf(w, "data: %d\n\n", i); err != nil {
					return
				}
				flusher.Flush()
			}
		}
	}
	startServer(t, streamingPort,
		WithTimeouts(time.Second, time.Second, writeTimeout, time.Second),
		WithService(&handlerService{path: "/stream", handler: stream(true)}),
		WithService(&handlerService{path: "/bounded", handler: stream(false)}),
	)

	patterns := []struct {
		desc         string
		path         string
		wantComplete bool
	}{
		{
			desc:         "handler that clears its write deadline streams past WriteTimeout",
			path:         "/stream",
			wantComplete: true,
		},
		{
			desc:         "handler that keeps the write deadline is cut off",
			path:         "/bounded",
			wantComplete: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			resp, err := newHTTPClient().Get(fmt.Sprintf("https://localhost:%d%s", streamingPort, p.path))
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			received := strings.Count(string(body), "data: ")
			if p.wantComplete {
				require.NoError(t, err)
				assert.Equal(t, chunks, received)
				return
			}
			assert.Error(t, err, "the response was not cut off")
			assert.Less(t, received, chunks)
		})
	}
}
