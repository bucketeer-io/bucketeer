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
	checker := NewRestChecker(version, service)
	checker.check(context.Background())
	req := httptest.NewRequest("GET", getTargetPath(t), nil)
	resp := httptest.NewRecorder()
	checker.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHTTPHealthy(t *testing.T) {
	healthyCheck := func(ctx context.Context) Status {
		return Healthy
	}
	checker := NewRestChecker(version, service, WithCheck("healthy", healthyCheck))
	checker.check(context.Background())
	req := httptest.NewRequest("GET", getTargetPath(t), nil)
	resp := httptest.NewRecorder()
	checker.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHTTPUnhealthy(t *testing.T) {
	unhealthyCheck := func(ctx context.Context) Status {
		return Unhealthy
	}
	checker := NewRestChecker(version, service, WithCheck("unhealthy", unhealthyCheck))
	checker.check(context.Background())
	req := httptest.NewRequest("GET", getTargetPath(t), nil)
	resp := httptest.NewRecorder()
	checker.ServeHTTP(resp, req)
	if resp.Code != http.StatusServiceUnavailable {
		t.Fail()
	}
}

func TestGRPCHealthyNoCheck(t *testing.T) {
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
