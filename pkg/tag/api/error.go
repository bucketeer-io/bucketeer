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
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

var (
	statusInternal     = api.NewGRPCStatus(err.NewErrorInternal(err.TagPackageName, "internal"))
	statusNameRequired = api.NewGRPCStatus(
		err.NewErrorInvalidArgEmpty(err.TagPackageName, "name must be specified", "name"),
	)
	statusEntityTypeRequired = api.NewGRPCStatus(
		err.NewErrorInvalidArgEmpty(err.TagPackageName, "entity_type must be specified", "entity_type"),
	)
	statusTagInUsed     = api.NewGRPCStatus(err.NewErrorFailedPrecondition(err.TagPackageName, "tag is in use"))
	statusInvalidCursor = api.NewGRPCStatus(
		err.NewErrorInvalidArgNotMatchFormat(err.TagPackageName, "cursor is invalid", "cursor"),
	)
	statusInvalidOrderBy = api.NewGRPCStatus(
		err.NewErrorInvalidArgNotMatchFormat(err.TagPackageName, "order_by is invalid", "order_by"),
	)
	statusUnauthenticated  = api.NewGRPCStatus(err.NewErrorUnauthenticated(err.TagPackageName, "unauthenticated"))
	statusPermissionDenied = api.NewGRPCStatus(err.NewErrorPermissionDenied(err.TagPackageName, "permission denied"))
)
