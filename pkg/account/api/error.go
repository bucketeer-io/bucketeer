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
	statusInternal              = gstatus.New(codes.Internal, "account: internal")
	statusInvalidCursor         = gstatus.New(codes.InvalidArgument, "account: cursor is invalid")
	statusNoCommand             = gstatus.New(codes.InvalidArgument, "account: command must not be empty")
	statusMissingOrganizationID = gstatus.New(
		codes.InvalidArgument,
		"account: organization id must be specified",
	)
	statusEmailIsEmpty                           = gstatus.New(codes.InvalidArgument, "account: email is empty")
	statusInvalidEmail                           = gstatus.New(codes.InvalidArgument, "account: invalid email format")
	statusNameIsEmpty                            = gstatus.New(codes.InvalidArgument, "account: name is empty")
	statusInvalidName                            = gstatus.New(codes.InvalidArgument, "account: invalid name format")
	statusInvalidOrganizationRole                = gstatus.New(codes.InvalidArgument, "account: invalid organization role")
	statusInvalidUpdateEnvironmentRolesWriteType = gstatus.New(
		codes.InvalidArgument,
		"account: invalid update environment roles write type",
	)
	statusMissingAPIKeyID   = gstatus.New(codes.InvalidArgument, "account: apikey id must be specified")
	statusMissingAPIKeyName = gstatus.New(codes.InvalidArgument, "account: apikey name must be not empty")
	statusInvalidOrderBy    = gstatus.New(codes.InvalidArgument, "account: order_by is invalid")
	statusNotFound          = gstatus.New(codes.NotFound, "account: not found")
	statusAlreadyExists     = gstatus.New(codes.AlreadyExists, "account: already exists")
	statusUnauthenticated   = gstatus.New(codes.Unauthenticated, "account: unauthenticated")
	statusPermissionDenied  = gstatus.New(codes.PermissionDenied, "account: permission denied")
)
