// Copyright 2022 The Bucketeer Authors.
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
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/status"
)

var (
	statusInternal            = gstatus.New(codes.Internal, "push: internal")
	statusIDRequired          = gstatus.New(codes.InvalidArgument, "push: id must be specified")
	statusNameRequired        = gstatus.New(codes.InvalidArgument, "push: name must be specified")
	statusFCMAPIKeyRequired   = gstatus.New(codes.InvalidArgument, "push: fcm api key must be specified")
	statusTagsRequired        = gstatus.New(codes.InvalidArgument, "push: tags must be specified")
	statusInvalidCursor       = gstatus.New(codes.InvalidArgument, "push: cursor is invalid")
	statusNoCommand           = gstatus.New(codes.InvalidArgument, "push: no command")
	statusInvalidOrderBy      = gstatus.New(codes.InvalidArgument, "push: order_by is invalid")
	statusNotFound            = gstatus.New(codes.NotFound, "push: not found")
	statusAlreadyDeleted      = gstatus.New(codes.NotFound, "push: already deleted")
	statusAlreadyExists       = gstatus.New(codes.AlreadyExists, "push: already exists")
	statusFCMKeyAlreadyExists = gstatus.New(codes.AlreadyExists, "push: fcm key already exists")
	statusTagAlreadyExists    = gstatus.New(codes.AlreadyExists, "push: tag already exists")
	statusUnauthenticated     = gstatus.New(codes.Unauthenticated, "push: unauthenticated")
	statusPermissionDenied    = gstatus.New(codes.PermissionDenied, "push: permission denied")

	errInternalJaJP = status.MustWithDetails(
		statusInternal,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "内部エラーが発生しました",
		},
	)
	errIDRequiredJaJP = status.MustWithDetails(
		statusIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "idは必須です",
		},
	)
	errNameRequiredJaJP = status.MustWithDetails(
		statusNameRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "nameは必須です",
		},
	)
	errFCMAPIKeyRequiredJaJP = status.MustWithDetails(
		statusFCMAPIKeyRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "fcm api keyは必須です",
		},
	)
	errTagsRequiredJaJP = status.MustWithDetails(
		statusTagsRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "tagsは必須です",
		},
	)
	errInvalidCursorJaJP = status.MustWithDetails(
		statusInvalidCursor,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なcursorです",
		},
	)
	errNoCommandJaJP = status.MustWithDetails(
		statusNoCommand,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "commandは必須です",
		},
	)
	errInvalidOrderByJaJP = status.MustWithDetails(
		statusInvalidOrderBy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なソート順の指定です",
		},
	)
	errNotFoundJaJP = status.MustWithDetails(
		statusNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "データが存在しません",
		},
	)
	errAlreadyDeletedJaJP = status.MustWithDetails(
		statusAlreadyDeleted,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "データがすでに削除済みです",
		},
	)
	errAlreadyExistsJaJP = status.MustWithDetails(
		statusAlreadyExists,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "同じidのデータがすでに存在します",
		},
	)
	errFCMKeyAlreadyExistsJaJP = status.MustWithDetails(
		statusFCMKeyAlreadyExists,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "同じfcm keyがすでに存在します",
		},
	)
	errTagAlreadyExistsJaJP = status.MustWithDetails(
		statusTagAlreadyExists,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "同じtagがすでに存在します",
		},
	)
	errUnauthenticatedJaJP = status.MustWithDetails(
		statusUnauthenticated,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "認証されていません",
		},
	)
	errPermissionDeniedJaJP = status.MustWithDetails(
		statusPermissionDenied,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "権限がありません",
		},
	)
)

func localizedError(s *gstatus.Status, loc string) error {
	// handle loc if multi-lang is necessary
	switch s {
	case statusInternal:
		return errInternalJaJP
	case statusIDRequired:
		return errIDRequiredJaJP
	case statusNameRequired:
		return errNameRequiredJaJP
	case statusFCMAPIKeyRequired:
		return errFCMAPIKeyRequiredJaJP
	case statusTagsRequired:
		return errTagsRequiredJaJP
	case statusInvalidCursor:
		return errInvalidCursorJaJP
	case statusNoCommand:
		return errNoCommandJaJP
	case statusInvalidOrderBy:
		return errInvalidOrderByJaJP
	case statusNotFound:
		return errNotFoundJaJP
	case statusAlreadyDeleted:
		return errAlreadyDeletedJaJP
	case statusAlreadyExists:
		return errAlreadyExistsJaJP
	case statusFCMKeyAlreadyExists:
		return errFCMKeyAlreadyExistsJaJP
	case statusTagAlreadyExists:
		return errTagAlreadyExistsJaJP
	case statusUnauthenticated:
		return errUnauthenticatedJaJP
	case statusPermissionDenied:
		return errPermissionDeniedJaJP
	default:
		return errInternalJaJP
	}
}
