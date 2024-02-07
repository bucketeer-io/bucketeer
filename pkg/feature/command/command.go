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

package command

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes"

	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	errBadCommand = errors.New("command: cannot handle command")
)

type Command interface{}

type Handler interface {
	Handle(ctx context.Context, cmd Command) error
}

func UnmarshalCommand(cmd *proto.Command) (Command, error) {
	var x ptypes.DynamicAny
	if err := ptypes.UnmarshalAny(cmd.Command, &x); err != nil {
		return nil, err
	}
	return x.Message, nil
}

// TODO: write test to unmarshal any and get correct type
