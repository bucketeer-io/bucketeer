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

package health

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	pb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type grpcChecker struct {
	*checker
}

func NewGrpcChecker(opts ...option) *grpcChecker {
	checker := &grpcChecker{
		checker: newChecker(opts...),
	}
	return checker
}

// ServeHTTP implements http.Handler for backward compatibility when grpcChecker is used as HTTP handler
// This allows the grpcChecker to respond to /health HTTP requests
// Routes: /health -> liveness, /ready -> readiness
func (hc *grpcChecker) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// Route to appropriate handler based on path
	switch req.URL.Path {
	case readyPath:
		hc.ServeReadyHTTP(resp, req)
	default:
		// Default to liveness for unknown paths
		hc.ServeLiveHTTP(resp, req)
	}
}

func (hc *grpcChecker) Register(server *grpc.Server) {
	pb.RegisterHealthServer(server, hc)
}

func (hc *grpcChecker) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	if hc.getStatus() == Unhealthy {
		return &pb.HealthCheckResponse{
			Status: pb.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}

func (hc *grpcChecker) Watch(*pb.HealthCheckRequest, pb.Health_WatchServer) error {
	// TODO: Implements here when needed.
	return status.Errorf(codes.Unimplemented, "unsupported method")
}

func (hc *grpcChecker) List(ctx context.Context, req *pb.HealthListRequest) (*pb.HealthListResponse, error) {
	statuses := make(map[string]*pb.HealthCheckResponse)

	// Add overall server health status
	if hc.getStatus() == Unhealthy {
		statuses[""] = &pb.HealthCheckResponse{
			Status: pb.HealthCheckResponse_NOT_SERVING,
		}
	} else {
		statuses[""] = &pb.HealthCheckResponse{
			Status: pb.HealthCheckResponse_SERVING,
		}
	}

	return &pb.HealthListResponse{
		Statuses: statuses,
	}, nil
}
