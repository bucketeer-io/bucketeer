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
	statusInternal            = gstatus.New(codes.Internal, "push: internal")
	statusIDRequired          = gstatus.New(codes.InvalidArgument, "push: id must be specified")
	statusNameRequired        = gstatus.New(codes.InvalidArgument, "push: name must be specified")
	statusFCMAPIKeyRequired   = gstatus.New(codes.InvalidArgument, "push: fcm api key must be specified")
	statusTagsRequired        = gstatus.New(codes.InvalidArgument, "push: tags must be specified")
	statusInvalidCursor       = gstatus.New(codes.InvalidArgument, "push: cursor is invalid")
	statusNoCommand           = gstatus.New(codes.InvalidArgument, "push: no command")
	statusInvalidOrderBy      = gstatus.New(codes.InvalidArgument, "push: order_by is invalid")
	statusNotFound            = gstatus.New(codes.NotFound, "push: not found")
	statusAlreadyExists       = gstatus.New(codes.AlreadyExists, "push: already exists")
	statusFCMKeyAlreadyExists = gstatus.New(codes.AlreadyExists, "push: fcm key already exists")
	statusTagAlreadyExists    = gstatus.New(codes.AlreadyExists, "push: tag already exists")
	statusUnauthenticated     = gstatus.New(codes.Unauthenticated, "push: unauthenticated")
	statusPermissionDenied    = gstatus.New(codes.PermissionDenied, "push: permission denied")
)
