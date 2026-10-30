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

package metrics

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/v2/pkg/health"
)

func TestCheckHealthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	}))
	defer ts.Close()
	m := NewMetrics(9002, "/metrics", WithHealthCheckURL(ts.URL))
	assert.Equal(t, health.Healthy, m.Check(context.TODO()))
}

func TestCheckUnhealthy(t *testing.T) {
	m := NewMetrics(9002, "/metrics")
	assert.Equal(t, health.Unhealthy, m.Check(context.TODO()))
}

func TestNewMetricsServerTimeouts(t *testing.T) {
	m := NewMetrics(9002, "/metrics").(*metrics)

	assert.Equal(t, 10*time.Second, m.server.ReadHeaderTimeout)
	assert.Equal(t, 30*time.Second, m.server.ReadTimeout)
	assert.Equal(t, 30*time.Second, m.server.WriteTimeout)
	assert.Equal(t, 60*time.Second, m.server.IdleTimeout)
}

// The scrape and health endpoints keep being served with the timeouts in place.
func TestRunServesMetricsAndHealth(t *testing.T) {
	const port = 9003
	m := NewMetrics(port, "/metrics")
	go func() {
		if err := m.Run(); err != nil {
			t.Errorf("Run failed: %v", err)
		}
	}()
	defer m.Stop()
	client := &http.Client{Timeout: time.Second}

	patterns := []struct {
		path     string
		wantBody string
	}{
		{path: "/health", wantBody: "healthy"},
		{path: "/metrics", wantBody: "go_goroutines"},
	}
	for _, p := range patterns {
		t.Run(p.path, func(t *testing.T) {
			var resp *http.Response
			require.Eventually(t, func() bool {
				var err error
				resp, err = client.Get(fmt.Sprintf("http://localhost:%d%s", port, p.path))
				return err == nil
			}, 5*time.Second, 10*time.Millisecond)
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Contains(t, string(body), p.wantBody)
		})
	}
}
