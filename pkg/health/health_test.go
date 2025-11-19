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

func TestStopPreventsCheckFromSettingHealthy(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc             string
		setupFunc        func() *restChecker
		expectedStatus   int
		expectedStopped  bool
		expectedInternal Status
	}{
		{
			desc: "check() after Stop() does not override status to Healthy",
			setupFunc: func() *restChecker {
				healthyCheck := func(ctx context.Context) Status {
					return Healthy
				}
				c := NewRestChecker(version, service, WithCheck("healthy", healthyCheck))
				c.check(context.Background())
				return c
			},
			expectedStatus:   http.StatusServiceUnavailable,
			expectedStopped:  true,
			expectedInternal: Unhealthy,
		},
		{
			desc: "multiple check() calls after Stop() remain Unhealthy",
			setupFunc: func() *restChecker {
				healthyCheck := func(ctx context.Context) Status {
					return Healthy
				}
				c := NewRestChecker(version, service, WithCheck("healthy", healthyCheck))
				c.check(context.Background())
				return c
			},
			expectedStatus:   http.StatusServiceUnavailable,
			expectedStopped:  true,
			expectedInternal: Unhealthy,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			checker := p.setupFunc()

			// Verify initial state is healthy
			if checker.getStatus() != Healthy {
				t.Errorf("Initial state should be Healthy, got %v", checker.getStatus())
			}

			// Call Stop()
			checker.Stop()

			// Verify stopped flag is set
			if !checker.isStopped() {
				t.Error("Expected isStopped() to be true after Stop()")
			}

			// Verify status is Unhealthy
			if checker.getStatus() != Unhealthy {
				t.Errorf("Expected status to be Unhealthy after Stop(), got %v", checker.getStatus())
			}

			// Call check() to simulate the race condition
			// This should NOT override the status back to Healthy
			checker.check(context.Background())

			// Verify status remains Unhealthy
			if checker.getStatus() != p.expectedInternal {
				t.Errorf("Expected status to remain %v after check(), got %v",
					p.expectedInternal, checker.getStatus())
			}

			// Verify HTTP response is 503
			req := httptest.NewRequest("GET", fmt.Sprintf("%s%s%s", version, service, readyPath), nil)
			resp := httptest.NewRecorder()
			checker.ServeReadyHTTP(resp, req)
			if resp.Code != p.expectedStatus {
				t.Errorf("Expected HTTP status %d, got %d", p.expectedStatus, resp.Code)
			}

			// Call check() multiple times to ensure it stays Unhealthy
			for i := 0; i < 5; i++ {
				checker.check(context.Background())
				if checker.getStatus() != p.expectedInternal {
					t.Errorf("After %d check() calls, expected status %v, got %v",
						i+1, p.expectedInternal, checker.getStatus())
				}
			}
		})
	}
}

func TestStopWithRunningGoroutine(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc          string
		setupFunc     func() (*restChecker, context.CancelFunc)
		expectedErr   error
		expectedFinal int
	}{
		{
			desc: "Stop() prevents Run() goroutine from setting status to Healthy",
			setupFunc: func() (*restChecker, context.CancelFunc) {
				healthyCheck := func(ctx context.Context) Status {
					return Healthy
				}
				c := NewRestChecker(version, service,
					WithCheck("healthy", healthyCheck),
					WithInterval(10*time.Millisecond))
				ctx, cancel := context.WithCancel(context.Background())
				go c.Run(ctx)
				// Wait for first check to complete
				time.Sleep(50 * time.Millisecond)
				return c, cancel
			},
			expectedErr:   nil,
			expectedFinal: http.StatusServiceUnavailable,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			checker, cancel := p.setupFunc()
			defer cancel()

			// Verify initial state is healthy
			if checker.getStatus() != Healthy {
				t.Errorf("Initial state should be Healthy, got %v", checker.getStatus())
			}

			// Call Stop() while Run() goroutine is still running
			checker.Stop()

			// Verify status is immediately Unhealthy
			if checker.getStatus() != Unhealthy {
				t.Errorf("Expected status to be Unhealthy immediately after Stop(), got %v",
					checker.getStatus())
			}

			// Wait for multiple check intervals to pass
			// The Run() goroutine should NOT override status back to Healthy
			time.Sleep(100 * time.Millisecond)

			// Verify status remains Unhealthy
			if checker.getStatus() != Unhealthy {
				t.Errorf("Expected status to remain Unhealthy after waiting, got %v",
					checker.getStatus())
			}

			// Verify HTTP response is 503
			req := httptest.NewRequest("GET", fmt.Sprintf("%s%s%s", version, service, readyPath), nil)
			resp := httptest.NewRecorder()
			checker.ServeReadyHTTP(resp, req)
			if resp.Code != p.expectedFinal {
				t.Errorf("Expected HTTP status %d, got %d", p.expectedFinal, resp.Code)
			}
		})
	}
}
