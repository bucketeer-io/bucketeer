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
	"github.com/bucketeer-io/bucketeer/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/pkg/error"
)

var (
	statusInternal      = api.NewGRPCStatus(error.NewErrorInternal(error.CoderefPackageName, "internal"))
	statusInvalidCursor = api.NewGRPCStatus(
		error.NewErrorInvalidArgEmpty(error.CoderefPackageName, "invalid cursor", "cursor"),
	)
	statusMissingID = api.NewGRPCStatus(
		error.NewErrorInvalidArgEmpty(error.CoderefPackageName, "missing id", "id"),
	)
	statusMissingFeatureID = api.NewGRPCStatus(
		error.NewErrorInvalidArgEmpty(error.CoderefPackageName, "missing feature_id", "feature_id"),
	)
	statusMissingFilePath = api.NewGRPCStatus(
		error.NewErrorInvalidArgEmpty(error.CoderefPackageName, "missing file_path", "file_path"),
	)
	statusMissingLineNumber = api.NewGRPCStatus(
		error.NewErrorInvalidArgNotMatchFormat(error.CoderefPackageName, "missing line_number", "line_number"),
	)
	statusMissingCodeSnippet = api.NewGRPCStatus(
		error.NewErrorInvalidArgEmpty(error.CoderefPackageName, "missing code_snippet", "code_snippet"),
	)
	statusMissingContentHash = api.NewGRPCStatus(
		error.NewErrorInvalidArgEmpty(error.CoderefPackageName, "missing content_hash", "content_hash"),
	)
	statusMissingRepositoryInfo = api.NewGRPCStatus(
		error.NewErrorInvalidArgEmpty(error.CoderefPackageName, "missing repository info", "repository_info"),
	)
	statusInvalidRepositoryType = api.NewGRPCStatus(
		error.NewErrorInvalidArgUnknown(error.CoderefPackageName, "invalid repository type", "repository_type"),
	)
	statusNotFound = api.NewGRPCStatus(
		error.NewErrorNotFound(error.CoderefPackageName, "not found", "coderef"),
	)
	statusUnauthenticated = api.NewGRPCStatus(
		error.NewErrorUnauthenticated(error.CoderefPackageName, "unauthenticated"),
	)
	statusPermissionDenied = api.NewGRPCStatus(
		error.NewErrorPermissionDenied(error.CoderefPackageName, "permission denied"),
	)
)
