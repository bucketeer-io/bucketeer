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
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// WithDetails attaches detail messages to the given Status and returns the updated Status.
// grpc/status.WithDetails is inefficient when attaching v2 messages because it goes through
// a v2 -> v1 -> v2 conversion path.
// To avoid that extra conversion, this helper operates directly on the Status proto
// representation and uses anypb.New with v2 messages.
// This implementation may be revisited once gRPC provides a v2-native API.
func WithDetails(s *status.Status, details ...proto.Message) (*status.Status, error) {
	if s.Code() == codes.OK {
		return nil, errors.New("no error details for status with code OK")
	}
	p := s.Proto()
	for _, detail := range details {
		m, err := anypb.New(detail)
		if err != nil {
			return nil, err
		}
		p.Details = append(p.Details, m)
	}
	return status.FromProto(p), nil
}

// MustWithDetails is like WithDetails, but panics if attaching details fails.
func MustWithDetails(s *status.Status, details ...proto.Message) error {
	dt, err := WithDetails(s, details...)
	if err != nil {
		panic(err)
	}

	return dt.Err()
}
