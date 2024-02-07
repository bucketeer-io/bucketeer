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

package client

import (
	"context"
	"io/ioutil"
	"strings"

	"google.golang.org/grpc/credentials"
)

type perRPCCredentials struct {
	APIKey string
}

func NewPerRPCCredentials(apiKeyPath string) (credentials.PerRPCCredentials, error) {
	data, err := ioutil.ReadFile(apiKeyPath)
	if err != nil {
		return nil, err
	}
	return perRPCCredentials{
		APIKey: strings.TrimSpace(string(data)),
	}, nil
}

func (c perRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": c.APIKey,
	}, nil
}

func (c perRPCCredentials) RequireTransportSecurity() bool {
	return true
}
