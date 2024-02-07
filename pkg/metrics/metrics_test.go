// Copyright 2024 The Bucketeer Authors.
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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/pkg/health"
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
