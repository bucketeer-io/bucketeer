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
	api "github.com/bucketeer-io/bucketeer/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
)

var (
	statusInternal                         = api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.NotificationPackageName, "internal"))
	statusIDRequired                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgument(pkgErr.NotificationPackageName, "id must be specified", pkgErr.InvalidTypeEmpty, "Id"))
	statusNameRequired                     = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgument(pkgErr.NotificationPackageName, "name must be specified", pkgErr.InvalidTypeEmpty, "Name"))
	statusSourceTypesRequired              = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgument(pkgErr.NotificationPackageName, "notification types must be specified", pkgErr.InvalidTypeEmpty, "SourceTypes"))
	statusUnknownRecipient                 = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgument(pkgErr.NotificationPackageName, "unknown recipient type", pkgErr.InvalidTypeNotMatchFormat, "Recipient"))
	statusRecipientRequired                = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgument(pkgErr.NotificationPackageName, "recipient must be specified", pkgErr.InvalidTypeEmpty, "Recipient"))
	statusSlackRecipientRequired           = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgument(pkgErr.NotificationPackageName, "slack recipient must be specified", pkgErr.InvalidTypeEmpty, "SlackRecipient"))
	statusSlackRecipientWebhookURLRequired = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgument(pkgErr.NotificationPackageName, "slack recipient webhook URL must be specified", pkgErr.InvalidTypeEmpty, "WebhookUrl"))
	statusInvalidCursor                    = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgument(pkgErr.NotificationPackageName, "cursor is invalid", pkgErr.InvalidTypeNotMatchFormat, "Cursor"))
	statusNoCommand                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgument(pkgErr.NotificationPackageName, "no command", pkgErr.InvalidTypeNotMatchFormat, "Command"))
	statusInvalidOrderBy                   = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgument(pkgErr.NotificationPackageName, "order_by is invalid", pkgErr.InvalidTypeNotMatchFormat, "OrderBy"))
	statusNotFound                         = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.NotificationPackageName, "not found"))
	statusAlreadyExists                    = api.NewGRPCStatus(pkgErr.NewErrorAlreadyExists(pkgErr.NotificationPackageName, "already exists"))
	statusUnauthenticated                  = api.NewGRPCStatus(pkgErr.NewErrorUnauthenticated(pkgErr.NotificationPackageName, "unauthenticated"))
	statusPermissionDenied                 = api.NewGRPCStatus(pkgErr.NewErrorPermissionDenied(pkgErr.NotificationPackageName, "permission denied"))
)
