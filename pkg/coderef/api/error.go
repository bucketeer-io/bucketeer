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
	"google.golang.org/grpc/status"
)

var (
	statusInternal              = status.New(codes.Internal, "coderef: internal")
	statusInvalidCursor         = status.New(codes.InvalidArgument, "coderef: invalid cursor")
	statusInvalidOrderBy        = status.New(codes.InvalidArgument, "coderef: invalid order_by")
	statusMissingID             = status.New(codes.InvalidArgument, "coderef: missing id")
	statusMissingEnvironmentID  = status.New(codes.InvalidArgument, "coderef: missing environment_id")
	statusMissingFeatureID      = status.New(codes.InvalidArgument, "coderef: missing feature_id")
	statusMissingFilePath       = status.New(codes.InvalidArgument, "coderef: missing file_path")
	statusMissingLineNumber     = status.New(codes.InvalidArgument, "coderef: missing line_number")
	statusMissingCodeSnippet    = status.New(codes.InvalidArgument, "coderef: missing code_snippet")
	statusMissingContentHash    = status.New(codes.InvalidArgument, "coderef: missing content_hash")
	statusMissingRepositoryInfo = status.New(codes.InvalidArgument, "coderef: missing repository info")
	statusInvalidRepositoryType = status.New(codes.InvalidArgument, "coderef: invalid repository type")
	statusNotFound              = status.New(codes.NotFound, "coderef: not found")
	statusUnauthenticated       = status.New(codes.Unauthenticated, "coderef: unauthenticated")
	statusPermissionDenied      = status.New(codes.PermissionDenied, "coderef: permission denied")
)
