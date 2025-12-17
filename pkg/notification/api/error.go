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
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

var (
	statusIDRequired = api.NewGRPCStatus(
		err.NewErrorInvalidArgEmpty(err.NotificationPackageName, "id must be specified", "ID"),
	)
	statusNameRequired = api.NewGRPCStatus(
		err.NewErrorInvalidArgEmpty(err.NotificationPackageName, "name must be specified", "Name"),
	)
	statusSourceTypesRequired = api.NewGRPCStatus(err.NewErrorInvalidArgEmpty(
		err.NotificationPackageName,
		"notification types must be specified",
		"NotificationType",
	))
	statusUnknownRecipient = api.NewGRPCStatus(
		err.NewErrorInvalidArgUnknown(err.NotificationPackageName, "unknown recipient", "NotificationRecipient"),
	)
	statusRecipientRequired = api.NewGRPCStatus(err.NewErrorInvalidArgEmpty(
		err.NotificationPackageName,
		"recipient must be specified",
		"NotificationRecipient",
	))
	statusSlackRecipientRequired = api.NewGRPCStatus(err.NewErrorInvalidArgEmpty(
		err.NotificationPackageName,
		"slack recipient must be specified",
		"NotificationSlackRecipient",
	))
	statusSlackRecipientWebhookURLRequired = api.NewGRPCStatus(err.NewErrorInvalidArgEmpty(
		err.NotificationPackageName,
		"webhook URL must be specified",
		"WebhookURL",
	))
	statusInvalidCursor = api.NewGRPCStatus(
		err.NewErrorInvalidArgNotMatchFormat(err.NotificationPackageName, "cursor is invalid", "Cursor"),
	)
	statusNoCommand = api.NewGRPCStatus(
		err.NewErrorInvalidArgEmpty(err.NotificationPackageName, "no command", "Command"),
	)
	statusInvalidOrderBy = api.NewGRPCStatus(
		err.NewErrorInvalidArgNotMatchFormat(err.NotificationPackageName, "order_by is invalid", "OrderBy"),
	)
	statusNotFound = api.NewGRPCStatus(
		err.NewErrorNotFound(err.NotificationPackageName, "not found", "Notification"),
	)
	statusAlreadyExists = api.NewGRPCStatus(
		err.NewErrorAlreadyExists(err.NotificationPackageName, "already exists"),
	)
	statusUnauthenticated = api.NewGRPCStatus(
		err.NewErrorUnauthenticated(err.NotificationPackageName, "unauthenticated"),
	)
	statusPermissionDenied = api.NewGRPCStatus(
		err.NewErrorPermissionDenied(err.NotificationPackageName, "permission denied"),
	)
)
