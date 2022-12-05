// Copyright 2022 The Bucketeer Authors.
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
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/status"
)

var (
	statusInternal         = gstatus.New(codes.Internal, "auditlog: internal")
	statusUnauthenticated  = gstatus.New(codes.Unauthenticated, "auditlog: unauthenticated")
	statusPermissionDenied = gstatus.New(codes.PermissionDenied, "auditlog: permission denied")
	statusInvalidCursor    = gstatus.New(codes.InvalidArgument, "auditlog: cursor is invalid")
	statusInvalidOrderBy   = gstatus.New(codes.InvalidArgument, "auditlog: order_by is invalid")
)
