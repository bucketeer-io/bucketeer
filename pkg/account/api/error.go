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
	statusEmailIsEmpty            = gstatus.New(codes.InvalidArgument, "account: email is empty")
	statusInvalidEmail            = gstatus.New(codes.InvalidArgument, "account: invalid email format")
	statusFirstNameIsEmpty        = gstatus.New(codes.InvalidArgument, "account: first name is empty")
	statusInvalidFirstName        = gstatus.New(codes.InvalidArgument, "account: invalid first name format")
	statusLastNameIsEmpty         = gstatus.New(codes.InvalidArgument, "account: last name is empty")
	statusInvalidLastName         = gstatus.New(codes.InvalidArgument, "account: invalid last name format")
	statusLanguageIsEmpty         = gstatus.New(codes.InvalidArgument, "account: language is empty")
	statusInvalidOrganizationRole = gstatus.New(codes.InvalidArgument, "account: invalid organization role")
	statusInvalidEnvironmentRole  = gstatus.New(
		codes.InvalidArgument,
		"account: environment roles must be specified",
	)
	statusInvalidUpdateEnvironmentRolesWriteType = gstatus.New(
		codes.InvalidArgument,
		"account: invalid update environment roles write type",
	)
	statusMissingAPIKeyID                  = gstatus.New(codes.InvalidArgument, "account: apikey id must be specified")
	statusMissingAPIKeyName                = gstatus.New(codes.InvalidArgument, "account: apikey name must be not empty")
	statusInvalidOrderBy                   = gstatus.New(codes.InvalidArgument, "account: order_by is invalid")
	statusNotFound                         = gstatus.New(codes.NotFound, "account: not found")
	statusAlreadyExists                    = gstatus.New(codes.AlreadyExists, "account: already exists")
	statusUnauthenticated                  = gstatus.New(codes.Unauthenticated, "account: unauthenticated")
	statusPermissionDenied                 = gstatus.New(codes.PermissionDenied, "account: permission denied")
	statusSearchFilterNameIsEmpty          = gstatus.New(codes.InvalidArgument, "account: search filter name is empty")
	statusSearchFilterQueryIsEmpty         = gstatus.New(codes.InvalidArgument, "account: search filter query is empty")
	statusSearchFilterTargetTypeIsRequired = gstatus.New(
		codes.InvalidArgument,
		"account: search filter target type is required",
	)
	statusSearchFilterIDIsEmpty    = gstatus.New(codes.InvalidArgument, "account: search filter ID is empty")
	statusSearchFilterIDNotFound   = gstatus.New(codes.InvalidArgument, "account: search filter ID not found")
	statusInvalidListAPIKeyRequest = gstatus.New(codes.InvalidArgument, "account: invalid list api key request")
)
