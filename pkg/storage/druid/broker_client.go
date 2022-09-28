// Copyright 2022 The Bucketeer Authors.
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

package druid

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ca-dp/godruid"
)

type BrokerClient struct {
	*godruid.Client
}

func NewBrokerClient(
	ctx context.Context,
	url,
	username,
	password string) (*BrokerClient, error) {
	return &BrokerClient{
		Client: &godruid.Client{
			Url:        fmt.Sprintf("http://%s:%s@%s", username, password, url),
			HttpClient: createHTTPClient(),
		},
	}, nil
}

func (c *BrokerClient) Close() {}

// Depending on the request timing for long requests, it could get an EOF error
// because the HTTP client reuses the connection for concurrent requests.
// We set the client transport manually, so the keep-alive setting can be disabled.
func createHTTPClient() *http.Client {
	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     true,
	}
	return &http.Client{
		Transport: t,
	}
}
