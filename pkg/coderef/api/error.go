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
	bkterr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

var (
	statusInternal      = api.NewGRPCStatus(bkterr.NewErrorInternal(bkterr.CoderefPackageName, "internal"))
	statusInvalidCursor = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(bkterr.CoderefPackageName, "invalid cursor", "cursor"),
	)
	statusMissingID = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(bkterr.CoderefPackageName, "missing id", "id"),
	)
	statusMissingFeatureID = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(bkterr.CoderefPackageName, "missing feature_id", "feature_id"),
	)
	statusMissingFilePath = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(bkterr.CoderefPackageName, "missing file_path", "file_path"),
	)
	statusMissingLineNumber = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgNotMatchFormat(bkterr.CoderefPackageName, "missing line_number", "line_number"),
	)
	statusMissingCodeSnippet = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(bkterr.CoderefPackageName, "missing code_snippet", "code_snippet"),
	)
	statusMissingContentHash = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(bkterr.CoderefPackageName, "missing content_hash", "content_hash"),
	)
	statusMissingRepositoryInfo = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(bkterr.CoderefPackageName, "missing repository info", "repository_info"),
	)
	statusInvalidRepositoryType = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgUnknown(bkterr.CoderefPackageName, "invalid repository type", "repository_type"),
	)
	statusNotFound = api.NewGRPCStatus(
		bkterr.NewErrorNotFound(bkterr.CoderefPackageName, "not found", "coderef"),
	)
	statusUnauthenticated = api.NewGRPCStatus(
		bkterr.NewErrorUnauthenticated(bkterr.CoderefPackageName, "unauthenticated"),
	)
	statusPermissionDenied = api.NewGRPCStatus(
		bkterr.NewErrorPermissionDenied(bkterr.CoderefPackageName, "permission denied"),
	)
)
