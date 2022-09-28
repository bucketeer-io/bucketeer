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
	"net/http"

	"github.com/ca-dp/godruid"
)

type OverlordClient struct {
	*godruid.OverlordClient
}

func NewOverlordClient(
	ctx context.Context,
	url,
	username,
	password string) (*OverlordClient, error) {

	return &OverlordClient{
		OverlordClient: &godruid.OverlordClient{
			Url:        fmt.Sprintf("http://%s:%s@%s", username, password, url),
			HttpClient: &http.Client{},
		},
	}, nil
}

func (c *OverlordClient) Close() {}
