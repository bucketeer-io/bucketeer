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

package api

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
)

func NewGRPCStatus(err error, anotherDetailData ...map[string]string) *status.Status {
	var pkg string
	var reason string
	var metadatas []map[string]string
	var st *status.Status
	var bucketeerErr *pkgErr.BucketeerError
	if err == nil {
		return status.New(codes.Unknown, "")
	}
	if errors.As(err, &bucketeerErr) {
		pkg = bucketeerErr.PackageName()
		bucketeerErr.AddMetadata(anotherDetailData...)
		var stCode codes.Code

		switch bucketeerErr.ErrorType() {
		case pkgErr.ErrorTypeInvalidAugment:
			reason = "INVALID_AUGMENT"
			stCode = codes.InvalidArgument
		case pkgErr.ErrorTypeNotFound:
			reason = "NOT_FOUND"
			stCode = codes.NotFound
		case pkgErr.ErrorTypeAlreadyExists:
			reason = "ALREADY_EXISTS"
			stCode = codes.AlreadyExists
		case pkgErr.ErrorTypeUnauthenticated:
			reason = "UNAUTHENTICATED"
			stCode = codes.Unauthenticated
		case pkgErr.ErrorTypePermissionDenied:
			reason = "PERMISSION_DENIED"
			stCode = codes.PermissionDenied
		case pkgErr.ErrorTypeUnexpectedAffectedRows:
			reason = "UNEXPECTED_AFFECTED_ROWS"
			stCode = codes.Internal
		case pkgErr.ErrorTypeInternal:
			reason = "INTERNAL"
			stCode = codes.Internal
		default:
			reason = "UNKNOWN"
			stCode = codes.Unknown
		}

		st = status.New(stCode, bucketeerErr.Message())
		metadatas = append(metadatas, bucketeerErr.Metadatas()...)

		//	if bucketeerErr, ok = err.(pkgErr.BucketeerError); ok {
		// pkg = bucketeerErr.PackageName()
		// bucketeerErr.AddMetadata(anotherDetailData...)
		// var stCode codes.Code

		// if errors.Is(bucketeerErr, &pkgErr.ErrorInvalidAugment{}) {
		// 	reason = "INVALID_AUGMENT"
		// 	stCode = codes.InvalidArgument
		// } else if errors.Is(bucketeerErr, &pkgErr.ErrorNotFound{}) {
		// 	reason = "NOT_FOUND"
		// 	stCode = codes.NotFound
		// } else if errors.Is(bucketeerErr, &pkgErr.ErrorAlreadyExists{}) {
		// 	reason = "ALREADY_EXISTS"
		// 	stCode = codes.AlreadyExists
		// } else if errors.Is(bucketeerErr, &pkgErr.ErrorUnauthenticated{}) {
		// 	reason = "UNAUTHENTICATED"
		// 	stCode = codes.Unauthenticated
		// } else if errors.Is(bucketeerErr, &pkgErr.ErrorPermissionDenied{}) {
		// 	reason = "PERMISSION_DENIED"
		// 	stCode = codes.PermissionDenied
		// } else if errors.Is(bucketeerErr, &pkgErr.ErrorUnexpectedAffectedRows{}) {
		// 	reason = "UNEXPECTED_AFFECTED_ROWS"
		// 	stCode = codes.Internal
		// } else if errors.Is(bucketeerErr, &pkgErr.ErrorInternal{}) {
		// 	reason = "INTERNAL"
		// 	stCode = codes.Internal
		// } else {
		// 	reason = "UNKNOWN"
		// 	stCode = codes.Unknown
		// }
		// st = status.New(stCode, bucketeerErr.Message())
	} else {
		pkg = "unknown"
		reason = "UNKNOWN"
		st = status.New(codes.Unknown, err.Error())
		if len(anotherDetailData) > 0 {
			metadatas = append(metadatas, anotherDetailData...)
		}
	}

	for _, md := range metadatas {
		st, err = st.WithDetails(&errdetails.ErrorInfo{
			Reason:   reason,
			Domain:   pkg + ".bucketeer.io",
			Metadata: md,
		})
		if err != nil {
			return status.New(codes.Internal, err.Error())
		}
	}
	return st
}
