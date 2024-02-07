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
	statusInternal             = gstatus.New(codes.Internal, "experiment: internal")
	statusInvalidCursor        = gstatus.New(codes.InvalidArgument, "experiment: cursor is invalid")
	statusNoCommand            = gstatus.New(codes.InvalidArgument, "experiment: must contain at least one command")
	statusFeatureIDRequired    = gstatus.New(codes.InvalidArgument, "experiment: feature id must be specified")
	statusExperimentIDRequired = gstatus.New(codes.InvalidArgument, "experiment: experiment id must be specified")
	statusGoalIDRequired       = gstatus.New(codes.InvalidArgument, "experiment: goal id must be specified")
	statusInvalidGoalID        = gstatus.New(codes.InvalidArgument, "experiment: invalid goal id")
	statusGoalNameRequired     = gstatus.New(codes.InvalidArgument, "experiment: goal name must be specified")
	statusPeriodTooLong        = gstatus.New(codes.InvalidArgument, "experiment: period too long")
	statusInvalidOrderBy       = gstatus.New(codes.InvalidArgument, "expriment: order_by is invalid")
	statusNotFound             = gstatus.New(codes.NotFound, "experiment: not found")
	statusGoalNotFound         = gstatus.New(codes.NotFound, "experiment: goal not found")
	statusFeatureNotFound      = gstatus.New(codes.NotFound, "experiment: feature not found")
	statusAlreadyExists        = gstatus.New(codes.AlreadyExists, "experiment: already exists")
	statusUnauthenticated      = gstatus.New(codes.Unauthenticated, "experiment: unauthenticated")
	statusPermissionDenied     = gstatus.New(codes.PermissionDenied, "experiment: permission denied")
)
