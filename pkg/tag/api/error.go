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
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

var (
	statusInternal           = gstatus.New(codes.Internal, "tag: internal")
	statusNameRequired       = gstatus.New(codes.InvalidArgument, "tag: name must be specified")
	statusEntityTypeRequired = gstatus.New(codes.InvalidArgument, "tag: entity_type must be specified")
	statusTagInUsed          = gstatus.New(codes.FailedPrecondition, "tag: tag is in use")
	statusInvalidCursor      = gstatus.New(codes.InvalidArgument, "tag: cursor is invalid")
	statusInvalidOrderBy     = gstatus.New(codes.InvalidArgument, "tag: order_by is invalid")
	statusUnauthenticated    = gstatus.New(codes.Unauthenticated, "tag: unauthenticated")
	statusPermissionDenied   = gstatus.New(codes.PermissionDenied, "tag: permission denied")
)
