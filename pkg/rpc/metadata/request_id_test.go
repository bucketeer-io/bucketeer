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
	"testing"

	"github.com/stretchr/testify/assert"
	gmetadata "google.golang.org/grpc/metadata"
)

func TestXGetRequestIDFromIncomingContext(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		ctx      context.Context
		expected string
	}{
		{
			desc:     "metadata doesn't exist",
			ctx:      context.Background(),
			expected: "",
		},
		{
			desc: "xRequestIDKey doesn't exist",
			ctx: gmetadata.NewIncomingContext(
				context.Background(),
				gmetadata.Pairs(),
			),
			expected: "",
		},
		{
			desc: "success",
			ctx: gmetadata.NewIncomingContext(
				context.Background(),
				gmetadata.Pairs(xRequestIDKey, "request-id-1"),
			),
			expected: "request-id-1",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := GetXRequestIDFromIncomingContext(p.ctx)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestXGetRequestIDFromOutgoingContext(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		ctx      context.Context
		expected string
	}{
		{
			desc:     "metadata doesn't exist",
			ctx:      context.Background(),
			expected: "",
		},
		{
			desc: "xRequestIDKey doesn't exist",
			ctx: gmetadata.NewOutgoingContext(
				context.Background(),
				gmetadata.Pairs(),
			),
			expected: "",
		},
		{
			desc: "success",
			ctx: gmetadata.NewOutgoingContext(
				context.Background(),
				gmetadata.Pairs(xRequestIDKey, "request-id-1"),
			),
			expected: "request-id-1",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := GetXRequestIDFromOutgoingContext(p.ctx)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestAppendXRequestIDToOutgoingContext(t *testing.T) {
	t.Parallel()
	ctx := gmetadata.NewOutgoingContext(
		context.Background(),
		gmetadata.Pairs(),
	)
	actualReqID := GetXRequestIDFromOutgoingContext(ctx)
	assert.Equal(t, "", actualReqID)
	expectedReqID := "request-id-1"
	ctx = AppendXRequestIDToOutgoingContext(ctx, expectedReqID)
	actualReqID = GetXRequestIDFromOutgoingContext(ctx)
	assert.Equal(t, expectedReqID, actualReqID)
}
