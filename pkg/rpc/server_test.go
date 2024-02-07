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

package rpc

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	pb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	proto "github.com/bucketeer-io/bucketeer/proto/test"
)

const certPath = "testdata/server.crt"
const keyPath = "testdata/server.key"

type testService struct {
}

func (s *testService) Register(server *grpc.Server) {
	proto.RegisterTestServiceServer(server, s)
}

func (s *testService) Test(ctx context.Context, req *proto.TestRequest) (*proto.TestResponse, error) {
	return &proto.TestResponse{Message: "test"}, nil
}

type dummyService struct {
}

func (s *dummyService) Register(server *grpc.Server) {
}

type dummyRegisterer struct {
}

func (s *dummyRegisterer) MustRegister(...prometheus.Collector) {

}

func (s *dummyRegisterer) Unregister(prometheus.Collector) bool {
	return true
}

type dummyVerifier struct {
}

func (v *dummyVerifier) Verify(rawIDToken string) (*token.IDToken, error) {
	return &token.IDToken{
		Email: "test@email",
	}, nil
}

type dummyPerRPCCredentials struct {
	Metadata map[string]string
}

func (c dummyPerRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return c.Metadata, nil
}

func (c dummyPerRPCCredentials) RequireTransportSecurity() bool {
	return true
}

func newServer(ctx context.Context) *Server {
	logger := zap.NewExample()
	health := health.NewGrpcChecker()
	server := NewServer(
		&testService{},
		certPath,
		keyPath,
		"test-server",
		WithService(health),
		WithVerifier(&dummyVerifier{}),
		WithMetrics(&dummyRegisterer{}),
		WithLogger(logger),
		WithPort(4443),
		WithHandler("/health", health),
	)
	return server
}

func newRPCClient(t *testing.T, rpcCreds credentials.PerRPCCredentials) proto.TestServiceClient {
	creds, err := credentials.NewClientTLSFromFile(certPath, "localhost")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := grpc.Dial("localhost:4443",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(creds),
		grpc.WithTimeout(5*time.Second),
		grpc.WithPerRPCCredentials(rpcCreds),
	)
	if err != nil {
		t.Fatal(err)
	}
	return proto.NewTestServiceClient(conn)
}

func newHealthClient(t *testing.T) pb.HealthClient {
	creds, err := credentials.NewClientTLSFromFile(certPath, "localhost")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := grpc.Dial("localhost:4443",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(creds),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}
	return pb.NewHealthClient(conn)
}

func TestGRPCHealthHandler(t *testing.T) {
	client := newHealthClient(t)
	resp, err := client.Check(context.TODO(), &pb.HealthCheckRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != pb.HealthCheckResponse_NOT_SERVING {
		t.Fatal(resp)
	}
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

func TestHTTPHealthHandler(t *testing.T) {
	client := newHTTPClient()
	resp, err := client.Get("https://localhost:4443/health")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatal(resp)
	}
}

func TestRPCHandlerOK(t *testing.T) {
	client := newRPCClient(t, &dummyPerRPCCredentials{
		Metadata: map[string]string{
			"authorization": "bearer dummy-token",
		},
	})
	resp, err := client.Test(context.TODO(), &proto.TestRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "test" {
		t.Fatal(resp)
	}
}

func TestRPCHandlerUnauthenticated(t *testing.T) {
	client := newRPCClient(t, &dummyPerRPCCredentials{})
	_, err := client.Test(context.TODO(), &proto.TestRequest{})
	if err == nil {
		t.Fatal("expected an error")
	}
	status, ok := status.FromError(err)
	if !ok {
		t.FailNow()
	}
	if status.Code() != codes.Unauthenticated {
		t.Fatal("code should be Unauthenticated")
	}
}

func TestMain(m *testing.M) {
	// Because os.Exit doesn't return, we need to call defer in separated function.
	code := testMain(m)
	os.Exit(code)
}

func testMain(m *testing.M) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server := newServer(ctx)
	defer server.Stop(time.Second)
	go server.Run()
	waitForServer()
	return m.Run()
}

func waitForServer() {
	client := newHTTPClient()
	for {
		resp, err := client.Get("https://localhost:4443")
		if err == nil {
			resp.Body.Close()
			return
		}
	}
}
