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

package client

import (
	"google.golang.org/grpc"

	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	proto "github.com/bucketeer-io/bucketeer/proto/gateway"
)

type Client interface {
	proto.GatewayClient
	Close()
}

type client struct {
	proto.GatewayClient
	address    string
	connection *grpc.ClientConn
}

func NewClient(addr, certPath string, opts ...rpcclient.Option) (Client, error) {
	conn, err := rpcclient.NewClientConn(addr, certPath, opts...)
	if err != nil {
		return nil, err
	}
	return &client{
		GatewayClient: proto.NewGatewayClient(conn),
		address:       addr,
		connection:    conn,
	}, nil
}

func (c *client) Close() {
	c.connection.Close()
}
