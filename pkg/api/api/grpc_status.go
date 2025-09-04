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

const (
	bktDomain = ".bucketeer.io"
)

func NewGRPCStatus(err error) *status.Status {
	var bucketeerErr *pkgErr.BktError
	var bktFieldError *pkgErr.BktFieldError
	if err == nil {
		return status.New(codes.Unknown, "")
	} else if errors.As(err, &bktFieldError) {
		return convertFieldError(bktFieldError)
	} else if errors.As(err, &bucketeerErr) {
		return convertBktError(bucketeerErr)
	} else {
		return convertUnknownError(err)
	}
}

func convertBktError(bktError *pkgErr.BktError) *status.Status {
	st := status.New(convertStatusCode(bktError.ErrorType()), bktError.Error())
	metadata := map[string]string{
		"messageKey": bktError.PackageName() + "." + string(bktError.ErrorType()),
	}

	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Reason:   convertErrorReason(bktError.ErrorType()),
		Domain:   bktError.PackageName() + bktDomain,
		Metadata: metadata,
	})
	if err != nil {
		return status.New(codes.Internal, err.Error())
	}
	return st
}

func convertFieldError(fieldError *pkgErr.BktFieldError) *status.Status {
	st := status.New(convertStatusCode(fieldError.ErrorType()), fieldError.Error())
	metadata := map[string]string{
		"messageKey": fieldError.PackageName() + "." + string(fieldError.ErrorType()),
		"field":      fieldError.Field(),
	}

	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Reason:   convertErrorReason(fieldError.ErrorType()),
		Domain:   fieldError.PackageName() + bktDomain,
		Metadata: metadata,
	})
	if err != nil {
		return status.New(codes.Internal, err.Error())
	}
	return st
}

func convertUnknownError(err error) *status.Status {
	st := status.New(codes.Unknown, err.Error())
	metadata := map[string]string{
		"messageKey": "unknown",
		"message":    err.Error(),
	}
	st, detailErr := st.WithDetails(&errdetails.ErrorInfo{
		Reason:   "UNKNOWN",
		Domain:   "unknown" + bktDomain,
		Metadata: metadata,
	})
	if detailErr != nil {
		return status.New(codes.Internal, detailErr.Error())
	}
	return st
}

func convertErrorReason(errorType pkgErr.ErrorType) string {
	switch errorType {
	case pkgErr.ErrorTypeInvalidArgEmpty:
		return "INVALID_ARGUMENT_EMPTY"
	case pkgErr.ErrorTypeInvalidArgNil:
		return "INVALID_ARGUMENT_NIL"
	case pkgErr.ErrorTypeInvalidArgNotMatchFormat:
		return "INVALID_ARGUMENT_NOT_MATCH_FORMAT"
	case pkgErr.ErrorTypeInvalidArgUnknown:
		return "INVALID_ARGUMENT"
	case pkgErr.ErrorTypeInvalidArgDuplicated:
		return "INVALID_ARGUMENT_DUPLICATED"
	case pkgErr.ErrorTypeNotFound:
		return "NOT_FOUND"
	case pkgErr.ErrorTypeAlreadyExists:
		return "ALREADY_EXISTS"
	case pkgErr.ErrorTypeUnauthenticated:
		return "UNAUTHENTICATED"
	case pkgErr.ErrorTypePermissionDenied:
		return "PERMISSION_DENIED"
	case pkgErr.ErrorTypeUnexpectedAffectedRows:
		return "UNEXPECTED_AFFECTED_ROWS"
	case pkgErr.ErrorTypeInternal:
		return "INTERNAL"
	default:
		return "UNKNOWN"
	}
}

func convertStatusCode(errorType pkgErr.ErrorType) codes.Code {
	switch errorType {
	case pkgErr.ErrorTypeInvalidArgUnknown,
		pkgErr.ErrorTypeInvalidArgEmpty,
		pkgErr.ErrorTypeInvalidArgNil,
		pkgErr.ErrorTypeInvalidArgNotMatchFormat,
		pkgErr.ErrorTypeInvalidArgDuplicated:
		return codes.InvalidArgument
	case pkgErr.ErrorTypeNotFound:
		return codes.NotFound
	case pkgErr.ErrorTypeAlreadyExists:
		return codes.AlreadyExists
	case pkgErr.ErrorTypeUnauthenticated:
		return codes.Unauthenticated
	case pkgErr.ErrorTypePermissionDenied:
		return codes.PermissionDenied
	case pkgErr.ErrorTypeUnexpectedAffectedRows:
		return codes.Internal
	case pkgErr.ErrorTypeInternal:
		return codes.Internal
	default:
		return codes.Unknown
	}
}
