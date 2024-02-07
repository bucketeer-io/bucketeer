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
//

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package client

import (
	"google.golang.org/grpc"

	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	proto "github.com/bucketeer-io/bucketeer/proto/experimentcalculator"
)

type Client interface {
	proto.ExperimentCalculatorServiceClient
	Close()
}

type client struct {
	proto.ExperimentCalculatorServiceClient
	address    string
	connection *grpc.ClientConn
}

func NewClient(addr, certPath string, opts ...rpcclient.Option) (Client, error) {
	conn, err := rpcclient.NewClientConn(addr, certPath, opts...)
	if err != nil {
		return nil, err
	}
	return &client{
		ExperimentCalculatorServiceClient: proto.NewExperimentCalculatorServiceClient(conn),
		address:                           addr,
		connection:                        conn,
	}, nil
}

func (c *client) Close() {
	c.connection.Close()
}
