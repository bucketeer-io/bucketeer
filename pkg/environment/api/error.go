// Copyright 2023 The Bucketeer Authors.
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
	statusInternal                   = gstatus.New(codes.Internal, "environment: internal")
	statusNoCommand                  = gstatus.New(codes.InvalidArgument, "environment: no command")
	statusInvalidCursor              = gstatus.New(codes.InvalidArgument, "environment: cursor is invalid")
	statusEnvironmentIDRequired      = gstatus.New(codes.InvalidArgument, "environment: environment id must be specified")
	statusInvalidEnvironmentID       = gstatus.New(codes.InvalidArgument, "environment: invalid environment id")
	statusProjectIDRequired          = gstatus.New(codes.InvalidArgument, "environment: project id must be specified")
	statusInvalidProjectName         = gstatus.New(codes.InvalidArgument, "environment: invalid project name")
	statusInvalidProjectUrlCode      = gstatus.New(codes.InvalidArgument, "environment: invalid project url code")
	statusInvalidProjectCreatorEmail = gstatus.New(codes.InvalidArgument, "environment: invalid project creator email")
	statusInvalidOrderBy             = gstatus.New(codes.InvalidArgument, "environment: order_by is invalid")
	statusEnvironmentNotFound        = gstatus.New(codes.NotFound, "environment: environment not found")
	statusProjectNotFound            = gstatus.New(codes.NotFound, "environment: project not found")
	statusEnvironmentAlreadyDeleted  = gstatus.New(codes.NotFound, "environment: environment already deleted")
	statusEnvironmentAlreadyExists   = gstatus.New(codes.AlreadyExists, "environment: environment already exists")
	statusProjectAlreadyExists       = gstatus.New(codes.AlreadyExists, "environment: project already exists")
	statusProjectDisabled            = gstatus.New(codes.FailedPrecondition, "environment: project disabled")
	statusUnauthenticated            = gstatus.New(codes.Unauthenticated, "environment: unauthenticated")
	statusPermissionDenied           = gstatus.New(codes.PermissionDenied, "environment: permission denied")
)
