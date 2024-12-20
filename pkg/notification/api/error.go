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

package api

import (
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

var (
	statusInternal            = gstatus.New(codes.Internal, "notification: internal")
	statusIDRequired          = gstatus.New(codes.InvalidArgument, "notification: id must be specified")
	statusNameRequired        = gstatus.New(codes.InvalidArgument, "notification: name must be specified")
	statusSourceTypesRequired = gstatus.New(
		codes.InvalidArgument,
		"notification: notification types must be specified",
	)
	statusUnknownRecipient  = gstatus.New(codes.InvalidArgument, "notification: unknown recipient")
	statusRecipientRequired = gstatus.New(
		codes.InvalidArgument,
		"notification: recipient must be specified",
	)
	statusSlackRecipientRequired = gstatus.New(
		codes.InvalidArgument,
		"notification: slack recipient must be specified",
	)
	statusSlackRecipientWebhookURLRequired = gstatus.New(
		codes.InvalidArgument,
		"notification: webhook URL must be specified",
	)
	statusInvalidCursor    = gstatus.New(codes.InvalidArgument, "notification: cursor is invalid")
	statusNoCommand        = gstatus.New(codes.InvalidArgument, "notification: no command")
	statusInvalidOrderBy   = gstatus.New(codes.InvalidArgument, "environment: order_by is invalid")
	statusNotFound         = gstatus.New(codes.NotFound, "notification: not found")
	statusAlreadyExists    = gstatus.New(codes.AlreadyExists, "notification: already exists")
	statusUnauthenticated  = gstatus.New(codes.Unauthenticated, "notification: unauthenticated")
	statusPermissionDenied = gstatus.New(codes.PermissionDenied, "notification: permission denied")
)
