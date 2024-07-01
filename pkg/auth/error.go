//  Copyright 2024 The Bucketeer Authors.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package auth

import (
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

var (
	StatusInternal                = gstatus.New(codes.Internal, "auth: internal")
	StatusMissingCode             = gstatus.New(codes.InvalidArgument, "auth: code must not be empty")
	StatusMissingState            = gstatus.New(codes.InvalidArgument, "auth: state must not be empty")
	StatusMissingAuthType         = gstatus.New(codes.InvalidArgument, "auth: missing authType")
	StatusUnknownAuthType         = gstatus.New(codes.InvalidArgument, "auth: unknown authType")
	StatusMissingRedirectURL      = gstatus.New(codes.InvalidArgument, "auth: missing redirectURL")
	StatusUnregisteredRedirectURL = gstatus.New(codes.InvalidArgument, "auth: unregistered redirectURL")
	StatusMissingRefreshToken     = gstatus.New(codes.InvalidArgument, "auth: refreshToken must not be empty")
	StatusInvalidCode             = gstatus.New(codes.InvalidArgument, "auth: invalid code")
	StatusInvalidRefreshToken     = gstatus.New(codes.InvalidArgument, "auth: invalid refresh token")
	StatusUnapprovedAccount       = gstatus.New(codes.PermissionDenied, "auth: unapproved account")
	StatusAccessDeniedEmail       = gstatus.New(codes.PermissionDenied, "auth: access denied email")
)
