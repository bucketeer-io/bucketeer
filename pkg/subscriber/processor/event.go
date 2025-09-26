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

package processor

import (
	"context"

	"github.com/golang/protobuf/proto"

	uproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

type eventDWHMap map[string]proto.Message
type environmentEventDWHMap map[string]eventDWHMap

type eventOPSMap map[string]proto.Message
type environmentEventOPSMap map[string]eventOPSMap

type Writer interface {
	Write(ctx context.Context, evt environmentEventDWHMap) map[string]bool
}

type Updater interface {
	UpdateUserCounts(ctx context.Context, events environmentEventOPSMap) map[string]bool
}

// Because the `userId` field in the EvaluationEvent proto message is already in the `User` field,
// we should remove it to avoid sending the same value twice.
// To keep compatibility, we must check both fields until all the SDKs are updated
func getUserID(userID string, user *uproto.User) string {
	if userID == "" {
		return user.Id
	}
	return userID
}
