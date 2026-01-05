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
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	bkterr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

var (
	statusInvalidCursor = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgNotMatchFormat(bkterr.CoderefPackageName, "invalid cursor", "Cursor"),
	)
	statusMissingID = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(bkterr.CoderefPackageName, "id is required", "ID"),
	)
	statusMissingFeatureID = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(bkterr.CoderefPackageName, "feature_id is required", "FeatureFlagID"),
	)
	statusMissingFilePath = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(bkterr.CoderefPackageName, "file_path is required", "CodeReferenceFilePath"),
	)
	statusMissingLineNumber = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgNotMatchFormat(
			bkterr.CoderefPackageName, "line_number must be greater than 0", "CodeReferenceLineNumber"),
	)
	statusMissingCodeSnippet = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(
			bkterr.CoderefPackageName, "code_snippet is required", "CodeReferenceCodeSnippet"),
	)
	statusMissingContentHash = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(
			bkterr.CoderefPackageName, "content_hash is required", "CodeReferenceContentHash"),
	)
	statusMissingRepositoryInfo = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(
			bkterr.CoderefPackageName, "repository info is required", "CodeReferenceRepositoryInfo"),
	)
	statusInvalidRepositoryType = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgUnknown(
			bkterr.CoderefPackageName, "invalid repository type", "CodeReferenceRepositoryType"),
	)
	statusNotFound = api.NewGRPCStatus(
		bkterr.NewErrorNotFound(bkterr.CoderefPackageName, "code reference not found", "CodeReference"),
	)
	statusUnauthenticated = api.NewGRPCStatus(
		bkterr.NewErrorUnauthenticated(bkterr.CoderefPackageName, "unauthenticated"),
	)
	statusPermissionDenied = api.NewGRPCStatus(
		bkterr.NewErrorPermissionDenied(bkterr.CoderefPackageName, "permission denied"),
	)
)
