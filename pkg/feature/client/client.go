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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package client

import (
	"google.golang.org/grpc"

	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// MaxSegmentUsersResponseBytes is the size of the largest response with
// segment users that a caller accepts. ListSegmentUsers returns all users of
// a segment in one response, and a gRPC client accepts 4MiB by default, which
// is about 41,000 users with 6-character IDs and about 25,700 with UUIDs.
// 64MiB holds more than 200,000 users with IDs of the maximum length of 100
// characters.
const MaxSegmentUsersResponseBytes = 64 * 1024 * 1024

// WithMaxSegmentUsersResponse makes a call accept a response of up to
// MaxSegmentUsersResponseBytes.
func WithMaxSegmentUsersResponse() grpc.CallOption {
	return grpc.MaxCallRecvMsgSize(MaxSegmentUsersResponseBytes)
}

type Client interface {
	proto.FeatureServiceClient
	Close()
}

type client struct {
	proto.FeatureServiceClient
	address    string
	connection *grpc.ClientConn
}

func NewClient(addr, certPath string, opts ...rpcclient.Option) (Client, error) {
	conn, err := rpcclient.NewClientConn(addr, certPath, opts...)
	if err != nil {
		return nil, err
	}
	return &client{
		FeatureServiceClient: proto.NewFeatureServiceClient(conn),
		address:              addr,
		connection:           conn,
	}, nil
}

func (c *client) Close() {
	c.connection.Close()
}
