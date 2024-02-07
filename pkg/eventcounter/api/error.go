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
	statusInternal               = gstatus.New(codes.Internal, "eventcounter: internal")
	statusFeatureIDRequired      = gstatus.New(codes.InvalidArgument, "eventcounter: feature id is required")
	statusFeatureVersionRequired = gstatus.New(codes.InvalidArgument, "eventcounter: feature version is required")
	statusVariationIDRequired    = gstatus.New(codes.InvalidArgument, "eventcounter: variation id is required")
	statusExperimentIDRequired   = gstatus.New(codes.InvalidArgument, "eventcounter: experiment id is required")
	statusMAUYearMonthRequired   = gstatus.New(codes.InvalidArgument, "eventcounter: mau year month is required")
	statusGoalIDRequired         = gstatus.New(codes.InvalidArgument, "eventcounter: goal id is required")
	statusStartAtRequired        = gstatus.New(codes.InvalidArgument, "eventcounter: start at is required")
	statusEndAtRequired          = gstatus.New(codes.InvalidArgument, "eventcounter: end at is required")
	statusStartAtIsAfterEndAt    = gstatus.New(codes.InvalidArgument, "eventcounter: start at is after end at")
	statusAutoOpsRuleIDRequired  = gstatus.New(codes.InvalidArgument, "eventcounter: auto ops rule id is required")
	statusClauseIDRequired       = gstatus.New(codes.InvalidArgument, "eventcounter: clause id is required")
	statusNotFound               = gstatus.New(codes.NotFound, "eventcounter: not found")
	statusUnauthenticated        = gstatus.New(codes.Unauthenticated, "feature: unauthenticated")
	statusPermissionDenied       = gstatus.New(codes.PermissionDenied, "feature: permission denied")
	statusUnknownTimeRange       = gstatus.New(codes.Internal, "eventcounter: unknown time range")
)
