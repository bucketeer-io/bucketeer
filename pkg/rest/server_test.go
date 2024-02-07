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

package rest

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
)

const (
	certPath = "testdata/server.crt"
	keyPath  = "testdata/server.key"
	timeout  = 10 * time.Second
	port     = 9222
)

type dummyService struct{}

func (*dummyService) Register(mux *http.ServeMux) {}

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
	go func() {
		<-time.After(timeout)
		fmt.Fprintln(os.Stderr, "failed to get response")
		os.Exit(1)
	}()
	for {
		resp, err := client.Get(fmt.Sprintf("https://localhost:%d", port))
		if err == nil {
			resp.Body.Close()
			return
		}
	}
}
