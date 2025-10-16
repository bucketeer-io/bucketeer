// Copyright 2025 The Bucketeer Authors.
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

package health

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	pb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	version = "/v1"
	service = "/gateway"
)

func TestHTTPHealthyNoCheck(t *testing.T) {
	t.Parallel()
	checker := NewRestChecker(version, service)
	checker.check(context.Background())
	req := httptest.NewRequest("GET", getTargetPath(t), nil)
	resp := httptest.NewRecorder()
	checker.ServeLiveHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHTTPHealthy(t *testing.T) {
	t.Parallel()
	healthyCheck := func(ctx context.Context) Status {
		return Healthy
	}
	checker := NewRestChecker(version, service, WithCheck("healthy", healthyCheck))
	checker.check(context.Background())
	req := httptest.NewRequest("GET", getTargetPath(t), nil)
	resp := httptest.NewRecorder()
	checker.ServeLiveHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fail()
	}
}

func TestGRPCHealthyNoCheck(t *testing.T) {
	t.Parallel()
	checker := NewGrpcChecker(WithInterval(time.Millisecond))
	checker.check(context.Background())
	resp, err := checker.Check(context.Background(), &pb.HealthCheckRequest{})
	if err != nil {
		t.Fail()
	}
	if resp.Status != pb.HealthCheckResponse_SERVING {
		t.Fail()
	}
}

func TestGRPCHealthy(t *testing.T) {
	t.Parallel()
	healthyCheck := func(ctx context.Context) Status {
		return Healthy
	}
	checker := NewGrpcChecker(WithCheck("healthy", healthyCheck))
	checker.check(context.Background())
	resp, err := checker.Check(context.Background(), &pb.HealthCheckRequest{})
	if err != nil {
		t.Fail()
	}
	if resp.Status != pb.HealthCheckResponse_SERVING {
		t.Fail()
	}
}

func TestGRPCUnhealthy(t *testing.T) {
	t.Parallel()
	unhealthyCheck := func(ctx context.Context) Status {
		return Unhealthy
	}
	checker := NewGrpcChecker(WithCheck("unhealthy", unhealthyCheck))
	checker.check(context.Background())
	resp, err := checker.Check(context.Background(), &pb.HealthCheckRequest{})
	if err != nil {
		t.Fail()
	}
	if resp.Status != pb.HealthCheckResponse_NOT_SERVING {
		t.Fail()
	}
}

func getTargetPath(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("%s%s%s", version, service, healthPath)
}

func TestHTTPReadyHealthy(t *testing.T) {
	t.Parallel()
	healthyCheck := func(ctx context.Context) Status {
		return Healthy
	}
	checker := NewRestChecker(version, service, WithCheck("healthy", healthyCheck))
	checker.check(context.Background())
	req := httptest.NewRequest("GET", fmt.Sprintf("%s%s%s", version, service, readyPath), nil)
	resp := httptest.NewRecorder()
	checker.ServeReadyHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}
}

func TestHTTPReadyUnhealthy(t *testing.T) {
	t.Parallel()
	unhealthyCheck := func(ctx context.Context) Status {
		return Unhealthy
	}
	checker := NewRestChecker(version, service, WithCheck("unhealthy", unhealthyCheck))
	checker.check(context.Background())
	req := httptest.NewRequest("GET", fmt.Sprintf("%s%s%s", version, service, readyPath), nil)
	resp := httptest.NewRecorder()
	checker.ServeReadyHTTP(resp, req)
	if resp.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, resp.Code)
	}
}

func TestHTTPHealthAffectedByStop(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		setupFunc      func() *restChecker
		expectedBefore int
		expectedAfter  int
	}{
		{
			desc: "health endpoint returns 503 after Stop() with healthy check",
			setupFunc: func() *restChecker {
				healthyCheck := func(ctx context.Context) Status {
					return Healthy
				}
				c := NewRestChecker(version, service, WithCheck("healthy", healthyCheck))
				c.check(context.Background())
				return c
			},
			expectedBefore: http.StatusOK,
			expectedAfter:  http.StatusServiceUnavailable,
		},
		{
			desc: "health endpoint returns 503 after Stop() with no checks",
			setupFunc: func() *restChecker {
				c := NewRestChecker(version, service)
				c.check(context.Background())
				return c
			},
			expectedBefore: http.StatusOK,
			expectedAfter:  http.StatusServiceUnavailable,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			checker := p.setupFunc()

			// Test before Stop()
			req := httptest.NewRequest("GET", fmt.Sprintf("%s%s%s", version, service, healthPath), nil)
			resp := httptest.NewRecorder()
			checker.ServeLiveHTTP(resp, req)
			if resp.Code != p.expectedBefore {
				t.Errorf("Before Stop(): Expected status %d, got %d", p.expectedBefore, resp.Code)
			}

			// Call Stop()
			checker.Stop()

			// Test after Stop() - health should return 503 to stop GCLB routing
			req = httptest.NewRequest("GET", fmt.Sprintf("%s%s%s", version, service, healthPath), nil)
			resp = httptest.NewRecorder()
			checker.ServeLiveHTTP(resp, req)
			if resp.Code != p.expectedAfter {
				t.Errorf("After Stop(): Expected status %d, got %d", p.expectedAfter, resp.Code)
			}
		})
	}
}

func TestHTTPReadyAffectedByStop(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		setupFunc      func() *restChecker
		expectedBefore int
		expectedAfter  int
	}{
		{
			desc: "ready endpoint returns 503 after Stop() with healthy check",
			setupFunc: func() *restChecker {
				healthyCheck := func(ctx context.Context) Status {
					return Healthy
				}
				c := NewRestChecker(version, service, WithCheck("healthy", healthyCheck))
				c.check(context.Background())
				return c
			},
			expectedBefore: http.StatusOK,
			expectedAfter:  http.StatusServiceUnavailable,
		},
		{
			desc: "ready endpoint returns 503 after Stop() with no checks",
			setupFunc: func() *restChecker {
				c := NewRestChecker(version, service)
				c.check(context.Background())
				return c
			},
			expectedBefore: http.StatusOK,
			expectedAfter:  http.StatusServiceUnavailable,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			checker := p.setupFunc()

			// Test before Stop()
			req := httptest.NewRequest("GET", fmt.Sprintf("%s%s%s", version, service, readyPath), nil)
			resp := httptest.NewRecorder()
			checker.ServeReadyHTTP(resp, req)
			if resp.Code != p.expectedBefore {
				t.Errorf("Before Stop(): Expected status %d, got %d", p.expectedBefore, resp.Code)
			}

			// Call Stop()
			checker.Stop()

			// Test after Stop() - ready should return 503
			req = httptest.NewRequest("GET", fmt.Sprintf("%s%s%s", version, service, readyPath), nil)
			resp = httptest.NewRecorder()
			checker.ServeReadyHTTP(resp, req)
			if resp.Code != p.expectedAfter {
				t.Errorf("After Stop(): Expected status %d, got %d", p.expectedAfter, resp.Code)
			}
		})
	}
}

func TestGRPCHealthAffectedByStop(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		setupFunc      func() *grpcChecker
		expectedBefore int
		expectedAfter  int
	}{
		{
			desc: "gRPC health endpoint returns 503 after Stop() with healthy check",
			setupFunc: func() *grpcChecker {
				healthyCheck := func(ctx context.Context) Status {
					return Healthy
				}
				c := NewGrpcChecker(WithCheck("healthy", healthyCheck))
				c.check(context.Background())
				return c
			},
			expectedBefore: http.StatusOK,
			expectedAfter:  http.StatusServiceUnavailable,
		},
		{
			desc: "gRPC health endpoint returns 503 after Stop() with no checks",
			setupFunc: func() *grpcChecker {
				c := NewGrpcChecker()
				c.check(context.Background())
				return c
			},
			expectedBefore: http.StatusOK,
			expectedAfter:  http.StatusServiceUnavailable,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			checker := p.setupFunc()

			// Test before Stop()
			req := httptest.NewRequest("GET", "/health", nil)
			resp := httptest.NewRecorder()
			checker.ServeHTTP(resp, req)
			if resp.Code != p.expectedBefore {
				t.Errorf("Before Stop(): Expected status %d, got %d", p.expectedBefore, resp.Code)
			}

			// Call Stop()
			checker.Stop()

			// Test after Stop() - health should return 503 to stop GCLB routing
			req = httptest.NewRequest("GET", "/health", nil)
			resp = httptest.NewRecorder()
			checker.ServeHTTP(resp, req)
			if resp.Code != p.expectedAfter {
				t.Errorf("After Stop(): Expected status %d, got %d", p.expectedAfter, resp.Code)
			}
		})
	}
}

func TestGRPCReadyAffectedByStop(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		setupFunc      func() *grpcChecker
		expectedBefore int
		expectedAfter  int
	}{
		{
			desc: "gRPC ready endpoint returns 503 after Stop() with healthy check",
			setupFunc: func() *grpcChecker {
				healthyCheck := func(ctx context.Context) Status {
					return Healthy
				}
				c := NewGrpcChecker(WithCheck("healthy", healthyCheck))
				c.check(context.Background())
				return c
			},
			expectedBefore: http.StatusOK,
			expectedAfter:  http.StatusServiceUnavailable,
		},
		{
			desc: "gRPC ready endpoint returns 503 after Stop() with no checks",
			setupFunc: func() *grpcChecker {
				c := NewGrpcChecker()
				c.check(context.Background())
				return c
			},
			expectedBefore: http.StatusOK,
			expectedAfter:  http.StatusServiceUnavailable,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			checker := p.setupFunc()

			// Test before Stop()
			req := httptest.NewRequest("GET", "/ready", nil)
			resp := httptest.NewRecorder()
			checker.ServeHTTP(resp, req)
			if resp.Code != p.expectedBefore {
				t.Errorf("Before Stop(): Expected status %d, got %d", p.expectedBefore, resp.Code)
			}

			// Call Stop()
			checker.Stop()

			// Test after Stop() - ready should return 503
			req = httptest.NewRequest("GET", "/ready", nil)
			resp = httptest.NewRecorder()
			checker.ServeHTTP(resp, req)
			if resp.Code != p.expectedAfter {
				t.Errorf("After Stop(): Expected status %d, got %d", p.expectedAfter, resp.Code)
			}
		})
	}
}
