// Copyright 2021, Google Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// The following codes based on https://github.com/googleapis/gax-go/blob/v2.7.0/v2/apierror/apierror_test.go.
// See the file headers for detail informations.

package writer

import (
	"testing"

	"cloud.google.com/go/bigquery/storage/apiv1/storagepb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetCodeFromError(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		expected string
		status   func() *status.Status
	}{
		{
			desc:     "codeOK",
			expected: codeOK,
		},
		{
			desc: "codeUnknown",
			status: func() *status.Status {
				status := status.New(codes.InvalidArgument, "invalid")
				return status
			},
			expected: codeUnknown,
		},
		{
			desc: "storageErrorCodeUnspecified",
			status: func() *status.Status {
				status, err := status.New(codes.InvalidArgument, "invalid").WithDetails(
					&storagepb.StorageError{Code: storagepb.StorageError_STORAGE_ERROR_CODE_UNSPECIFIED},
				)
				require.NoError(t, err)
				return status
			},
			expected: storageErrorCodeUnspecified,
		},
		{
			desc: "tableNotFound",
			status: func() *status.Status {
				status, err := status.New(codes.InvalidArgument, "invalid").WithDetails(
					&storagepb.StorageError{Code: storagepb.StorageError_TABLE_NOT_FOUND},
				)
				require.NoError(t, err)
				return status
			},
			expected: tableNotFound,
		},
		{
			desc: "streamAlreadyCommitted",
			status: func() *status.Status {
				status, err := status.New(codes.InvalidArgument, "invalid").WithDetails(
					&storagepb.StorageError{Code: storagepb.StorageError_STREAM_ALREADY_COMMITTED},
				)
				require.NoError(t, err)
				return status
			},
			expected: streamAlreadyCommitted,
		},
		{
			desc: "sreamNotFound",
			status: func() *status.Status {
				status, err := status.New(codes.InvalidArgument, "invalid").WithDetails(
					&storagepb.StorageError{Code: storagepb.StorageError_STREAM_NOT_FOUND},
				)
				require.NoError(t, err)
				return status
			},
			expected: sreamNotFound,
		},
		{
			desc: "invalidStreamType",
			status: func() *status.Status {
				status, err := status.New(codes.InvalidArgument, "invalid").WithDetails(
					&storagepb.StorageError{Code: storagepb.StorageError_INVALID_STREAM_TYPE},
				)
				require.NoError(t, err)
				return status
			},
			expected: invalidStreamType,
		},
		{
			desc: "streamFinalized",
			status: func() *status.Status {
				status, err := status.New(codes.InvalidArgument, "invalid").WithDetails(
					&storagepb.StorageError{Code: storagepb.StorageError_STREAM_FINALIZED},
				)
				require.NoError(t, err)
				return status
			},
			expected: streamFinalized,
		},
		{
			desc: "schemaMismatchExtraFields",
			status: func() *status.Status {
				status, err := status.New(codes.InvalidArgument, "invalid").WithDetails(
					&storagepb.StorageError{Code: storagepb.StorageError_SCHEMA_MISMATCH_EXTRA_FIELDS},
				)
				require.NoError(t, err)
				return status
			},
			expected: schemaMismatchExtraFields,
		},
		{
			desc: "offsetAlreadyExists",
			status: func() *status.Status {
				status, err := status.New(codes.InvalidArgument, "invalid").WithDetails(
					&storagepb.StorageError{Code: storagepb.StorageError_OFFSET_ALREADY_EXISTS},
				)
				require.NoError(t, err)
				return status
			},
			expected: offsetAlreadyExists,
		},
		{
			desc: "offsetOutOfRange",
			status: func() *status.Status {
				status, err := status.New(codes.InvalidArgument, "invalid").WithDetails(
					&storagepb.StorageError{Code: storagepb.StorageError_OFFSET_OUT_OF_RANGE},
				)
				require.NoError(t, err)
				return status
			},
			expected: offsetOutOfRange,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			var err error
			if p.status != nil {
				err = p.status().Err()
			}
			actual := getCodeFromError(err)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
		})
	}
}
