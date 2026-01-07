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

package api

import (
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

var (
	statusInternal   = api.NewGRPCStatus(err.NewErrorInternal(err.PushPackageName, "internal"))
	statusIDRequired = api.NewGRPCStatus(
		err.NewErrorInvalidArgEmpty(err.PushPackageName, "id must be specified", "ID"),
	)
	statusNameRequired = api.NewGRPCStatus(
		err.NewErrorInvalidArgEmpty(err.PushPackageName, "name must be specified", "Name"),
	)
	statusFCMServiceAccountRequired = api.NewGRPCStatus(
		err.NewErrorInvalidArgEmpty(err.PushPackageName, "fcm service account must be specified", "FCMServiceAccount"),
	)
	statusFCMServiceAccountInvalid = api.NewGRPCStatus(
		err.NewErrorInvalidArgNotMatchFormat(err.PushPackageName, "fcm service account is invalid", "FCMServiceAccount"),
	)
	statusInvalidCursor = api.NewGRPCStatus(
		err.NewErrorInvalidArgNotMatchFormat(err.PushPackageName, "cursor is invalid", "Cursor"),
	)
	statusInvalidOrderBy = api.NewGRPCStatus(
		err.NewErrorInvalidArgUnknown(err.PushPackageName, "order_by is invalid", "OrderBy"),
	)
	statusNotFound = api.NewGRPCStatus(
		err.NewErrorNotFound(err.PushPackageName, "not found", "Push"),
	)
	statusAlreadyExists = api.NewGRPCStatus(
		err.NewErrorAlreadyExists(err.PushPackageName, "already exists"),
	)
	statusFCMServiceAccountAlreadyExists = api.NewGRPCStatus(
		err.NewErrorAlreadyExists(err.PushPackageName, "fcm service account already exists"),
	)
	statusTagAlreadyExists = api.NewGRPCStatus(
		err.NewErrorAlreadyExists(err.PushPackageName, "tag already exists"),
	)
	statusUnauthenticated = api.NewGRPCStatus(
		err.NewErrorUnauthenticated(err.PushPackageName, "unauthenticated"),
	)
	statusPermissionDenied = api.NewGRPCStatus(err.NewErrorPermissionDenied(err.PushPackageName, "permission denied"))
)
