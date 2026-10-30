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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// newGatewayMux mirrors the SDK gateway server: exact API routes only.
func newGatewayMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/gateway/ping", okHandler)
	mux.HandleFunc("/v1/gateway/evaluations", okHandler)
	mux.HandleFunc("/v1/gateway/health", okHandler)
	return mux
}

// newDashboardMux mirrors the dashboard server: subtree routes for static
// files, the SPA fallback on "/", and one exact API route.
func newDashboardMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/assets/", okHandler)
	mux.HandleFunc("/static/js/", okHandler)
	mux.HandleFunc("/v1/aichat/chat", okHandler)
	mux.HandleFunc("/", okHandler)
	return mux
}

func TestRouteLabels(t *testing.T) {
	t.Parallel()

	gatewayMux := newGatewayMux()
	dashboardMux := newDashboardMux()
	patterns := []struct {
		desc            string
		mux             *http.ServeMux
		target          string
		expectedVersion string
		expectedService string
		expectedMethod  string
	}{
		{
			desc:            "gateway: registered api route",
			mux:             gatewayMux,
			target:          "http://example.com/v1/gateway/evaluations",
			expectedVersion: "v1",
			expectedService: "gateway",
			expectedMethod:  "evaluations",
		},
		{
			desc:            "gateway: registered api route with query string",
			mux:             gatewayMux,
			target:          "http://example.com/v1/gateway/ping?tag=ios",
			expectedVersion: "v1",
			expectedService: "gateway",
			expectedMethod:  "ping",
		},
		{
			desc:            "gateway: unregistered method under a registered service is not found",
			mux:             gatewayMux,
			target:          "http://example.com/v1/gateway/no-such-method",
			expectedVersion: unknownLabelValue,
			expectedService: unknownLabelValue,
			expectedMethod:  unknownLabelValue,
		},
		{
			desc:            "gateway: unregistered three-segment path is not found",
			mux:             gatewayMux,
			target:          "http://example.com/x1/y1/z1",
			expectedVersion: unknownLabelValue,
			expectedService: unknownLabelValue,
			expectedMethod:  unknownLabelValue,
		},
		{
			desc:            "dashboard: registered api route",
			mux:             dashboardMux,
			target:          "http://example.com/v1/aichat/chat",
			expectedVersion: "v1",
			expectedService: "aichat",
			expectedMethod:  "chat",
		},
		{
			desc:            "dashboard: arbitrary three-segment path served by the spa fallback",
			mux:             dashboardMux,
			target:          "http://example.com/x1/y1/z1",
			expectedVersion: unknownLabelValue,
			expectedService: unknownLabelValue,
			expectedMethod:  unknownLabelValue,
		},
		{
			desc:            "dashboard: file below the assets subtree",
			mux:             dashboardMux,
			target:          "http://example.com/assets/x1/y1.js",
			expectedVersion: unknownLabelValue,
			expectedService: unknownLabelValue,
			expectedMethod:  unknownLabelValue,
		},
		{
			desc:            "dashboard: file below the static js subtree",
			mux:             dashboardMux,
			target:          "http://example.com/static/js/env.js",
			expectedVersion: unknownLabelValue,
			expectedService: unknownLabelValue,
			expectedMethod:  unknownLabelValue,
		},
		{
			desc:            "dashboard: spa root",
			mux:             dashboardMux,
			target:          "http://example.com/",
			expectedVersion: unknownLabelValue,
			expectedService: unknownLabelValue,
			expectedMethod:  unknownLabelValue,
		},
		{
			desc:            "dashboard: percent-encoded segment that is not valid utf-8",
			mux:             dashboardMux,
			target:          "http://example.com/v1/aichat/%85",
			expectedVersion: unknownLabelValue,
			expectedService: unknownLabelValue,
			expectedMethod:  unknownLabelValue,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, p.target, nil)
			version, service, method := routeLabels(p.mux, req)
			assert.Equal(t, p.expectedVersion, version)
			assert.Equal(t, p.expectedService, service)
			assert.Equal(t, p.expectedMethod, method)
		})
	}
}

// Requests to paths that no exact route serves must all land in the same
// "unknown" series, while registered routes keep their own series.
// Not parallel: it asserts exact deltas on the package-level counters.
func TestMetricsServerMiddlewareBoundsSeriesToRegisteredRoutes(t *testing.T) {
	const arbitraryRequests = 20
	mux := newDashboardMux()
	handler := MetricsServerMiddleware(mux)(mux)
	reg := prometheus.NewRegistry()
	reg.MustRegister(serverStartedCounter, serverHandledCounter, serverHandledHistogram)

	unknownStarted := serverStartedCounter.WithLabelValues(unknownLabelValue, unknownLabelValue, unknownLabelValue)
	unknownHandled := serverHandledCounter.WithLabelValues(unknownLabelValue, unknownLabelValue, unknownLabelValue, "200")
	chatStarted := serverStartedCounter.WithLabelValues("v1", "aichat", "chat")
	chatHandled := serverHandledCounter.WithLabelValues("v1", "aichat", "chat", "200")
	unknownStartedBefore := testutil.ToFloat64(unknownStarted)
	unknownHandledBefore := testutil.ToFloat64(unknownHandled)
	chatStartedBefore := testutil.ToFloat64(chatStarted)
	chatHandledBefore := testutil.ToFloat64(chatHandled)

	for i := 0; i < arbitraryRequests; i++ {
		target := fmt.Sprintf("http://example.com/arbitrary-v%d/arbitrary-s%d/arbitrary-m%d", i, i, i)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
		require.Equal(t, http.StatusOK, rec.Code)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "http://example.com/v1/aichat/chat", nil))
	require.Equal(t, http.StatusOK, rec.Code)

	assert.Equal(t, float64(arbitraryRequests), testutil.ToFloat64(unknownStarted)-unknownStartedBefore)
	assert.Equal(t, float64(arbitraryRequests), testutil.ToFloat64(unknownHandled)-unknownHandledBefore)
	assert.Equal(t, float64(1), testutil.ToFloat64(chatStarted)-chatStartedBefore)
	assert.Equal(t, float64(1), testutil.ToFloat64(chatHandled)-chatHandledBefore)

	// No series was created from the segments of the arbitrary request paths.
	families, err := reg.Gather()
	require.NoError(t, err)
	for _, family := range families {
		for _, metric := range family.GetMetric() {
			for _, label := range metric.GetLabel() {
				assert.False(t, strings.HasPrefix(label.GetValue(), "arbitrary-"),
					"%s has a series labelled with a request path segment: %s=%s",
					family.GetName(), label.GetName(), label.GetValue())
			}
		}
	}
}
