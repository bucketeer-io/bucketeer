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

package metadata

import (
	"context"

	gmetadata "google.golang.org/grpc/metadata"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
)

const xRequestIDKey = "x-request-id"

func GetXRequestIDFromIncomingContext(ctx context.Context) string {
	md, ok := gmetadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	reqIDs, ok := md[xRequestIDKey]
	if !ok || len(reqIDs) == 0 {
		return ""
	}
	return reqIDs[0]
}

func GetXRequestIDFromOutgoingContext(ctx context.Context) string {
	md, ok := gmetadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	reqIDs, ok := md[xRequestIDKey]
	if !ok || len(reqIDs) == 0 {
		return ""
	}
	return reqIDs[0]
}

func AppendXRequestIDToOutgoingContext(ctx context.Context, xRequestID string) context.Context {
	return gmetadata.AppendToOutgoingContext(ctx, xRequestIDKey, xRequestID)
}

func GenerateXRequestID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		return ""
	}
	return id.String()
}
