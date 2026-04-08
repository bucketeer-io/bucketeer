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

package status

import (
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
)

func MustWithDetails(s *status.Status, details ...proto.Message) error {
	v1Messages := make([]protoadapt.MessageV1, 0, 10)
	for _, v2Message := range details {
		v1Messages = append(v1Messages, protoadapt.MessageV1Of(v2Message))
	}

	dt, err := s.WithDetails(v1Messages...)
	if err != nil {
		panic(err)
	}
	return dt.Err()
}
