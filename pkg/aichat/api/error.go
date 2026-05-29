// Copyright 2026 The Bucketeer Authors.
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

	api "github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

var (
	statusUnauthenticated = api.NewGRPCStatus(
		pkgErr.NewErrorUnauthenticated(pkgErr.AIChatPackageName, "unauthenticated"))
	statusPermissionDenied = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.AIChatPackageName, "permission denied"))
	statusMissingEnvironmentID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AIChatPackageName, "missing environment id", "EnvironmentId"))
	statusMissingMessages = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AIChatPackageName, "missing messages", "Messages"))
	statusTooManyMessages = api.NewGRPCStatus(
		pkgErr.NewErrorExceededMax(pkgErr.AIChatPackageName, "too many messages", "Messages", maxMessages))
	statusRateLimitExceeded = status.New(codes.ResourceExhausted, "aichat:rate limit exceeded")
	statusRequestCanceled   = status.New(codes.Canceled, "aichat:request canceled")
	statusChatFailed        = api.NewGRPCStatus(
		pkgErr.NewErrorInternal(pkgErr.AIChatPackageName, "chat generation failed"))
)
