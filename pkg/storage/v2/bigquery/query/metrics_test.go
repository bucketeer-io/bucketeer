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

package query

import (
	"testing"

	"cloud.google.com/go/bigquery/storage/apiv1/storagepb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// This test is based on 
func TestGetCodeFromError(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	a := &storagepb.StorageError{Code: storagepb.StorageError_INVALID_STREAM_STATE}
	b, err := status.New(codes.InvalidArgument, "invalid").WithDetails(a)
	require.NoError(t, err)

	patterns := []struct {
		desc     string
		input    error
		expected string
	}{
		{
			desc:     "codeOK",
			input:    nil,
			expected: codeOK,
		},
		// {
		// 	desc:     "codeUnknown: unknow error",
		// 	input:    errors.New("error"),
		// 	expected: codeUnknown,
		// },
		{
			desc:     "codeUnknown: error code is not unexpected",
			input:    b.Err(),
			expected: invalidStreamState,
		},
		// {
		// 	desc:     "codeBadRequest",
		// 	input:    &apierror.APIError{err: eS.Err(), status: eS, details: ErrDetails{ErrorInfo: ei}},
		// 	expected: codeBadRequest,
		// },
		// {
		// 	desc:     "codeForbidden",
		// 	input:    &googleapi.Error{Code: http.StatusForbidden},
		// 	expected: codeForbidden,
		// },
		// {
		// 	desc:     "codeNotFound",
		// 	input:    &googleapi.Error{Code: http.StatusNotFound},
		// 	expected: codeNotFound,
		// },
		// {
		// 	desc:     "codeConflict",
		// 	input:    &googleapi.Error{Code: http.StatusConflict},
		// 	expected: codeConflict,
		// },
		// {
		// 	desc:     "codeInternalServerError",
		// 	input:    &googleapi.Error{Code: http.StatusInternalServerError},
		// 	expected: codeInternalServerError,
		// },
		// {
		// 	desc:     "codeNotImplemented",
		// 	input:    &googleapi.Error{Code: http.StatusNotImplemented},
		// 	expected: codeNotImplemented,
		// },
		// {
		// 	desc:     "codeServiceUnavailable",
		// 	input:    &googleapi.Error{Code: http.StatusServiceUnavailable},
		// 	expected: codeServiceUnavailable,
		// },
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := getCodeFromError(p.input)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
		})
	}
}
