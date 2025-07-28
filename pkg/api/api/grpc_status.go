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

	if invalidAugmentError := (*pkgErr.ErrorInvalidAugment)(nil); errors.As(err, &invalidAugmentError) {
		pkg = invalidAugmentError.PackageName
		st = status.New(codes.InvalidArgument, invalidAugmentError.Message)
		reason = "INVALID_AUGMENT"
		metadatas = append(metadatas, invalidAugmentError.Metadatas...)
	} else if notFoundError := (*pkgErr.ErrorNotFound)(nil); errors.As(err, &notFoundError) {
		pkg = notFoundError.PackageName
		st = status.New(codes.NotFound, notFoundError.Message)
		reason = "NOT_FOUND"
		metadatas = append(metadatas, notFoundError.Metadatas...)
	} else if alreadyExistsError := (*pkgErr.ErrorAlreadyExists)(nil); errors.As(err, &alreadyExistsError) {
		pkg = alreadyExistsError.PackageName
		st = status.New(codes.AlreadyExists, alreadyExistsError.Message)
		reason = "ALREADY_EXISTS"
		metadatas = append(metadatas, alreadyExistsError.Metadatas...)
	} else if unauthenticatedError := (*pkgErr.ErrorUnauthenticated)(nil); errors.As(err, &unauthenticatedError) {
		pkg = unauthenticatedError.PackageName
		st = status.New(codes.Unauthenticated, unauthenticatedError.Message)
		reason = "UNAUTHENTICATED"
		metadatas = append(metadatas, unauthenticatedError.Metadatas...)
	} else if permissionDeniedError := (*pkgErr.ErrorPermissionDenied)(nil); errors.As(err, &permissionDeniedError) {
		pkg = permissionDeniedError.PackageName
		st = status.New(codes.PermissionDenied, permissionDeniedError.Message)
		reason = "PERMISSION_DENIED"
		metadatas = append(metadatas, permissionDeniedError.Metadatas...)
	} else if unexpectedAffectedRowsError := (*pkgErr.ErrorUnexpectedAffectedRows)(nil); errors.As(err, &unexpectedAffectedRowsError) {
		pkg = unexpectedAffectedRowsError.PackageName
		st = status.New(codes.Internal, unexpectedAffectedRowsError.Message)
		reason = "UNEXPECTED_AFFECTED_ROWS"
		metadatas = append(metadatas, unexpectedAffectedRowsError.Metadatas...)
	} else if internalError := (*pkgErr.ErrorInternal)(nil); errors.As(err, &internalError) {
		pkg = internalError.PackageName
		st = status.New(codes.Internal, internalError.Message)
		reason = "INTERNAL"
		metadatas = append(metadatas, internalError.Metadatas...)
	} else {
		pkg = "unknown"
		st = status.New(codes.Unknown, err.Error())
		reason = "UNKNOWN"
	}
	// when adding multiple details
	for _, md := range anotherDetailData {
		for k, v := range md {
			metadatas = append(metadatas, map[string]string{
				k: v,
			})
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
