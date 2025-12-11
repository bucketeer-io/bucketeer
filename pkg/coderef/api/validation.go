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
	proto "github.com/bucketeer-io/bucketeer/v2/proto/coderef"
)

func validateCreateCodeReferenceRequest(req *proto.CreateCodeReferenceRequest) error {
	if req.FeatureId == "" {
		return statusMissingFeatureID.Err()
	}
	if req.FilePath == "" {
		return statusMissingFilePath.Err()
	}
	if req.LineNumber <= 0 {
		return statusMissingLineNumber.Err()
	}
	if req.CodeSnippet == "" {
		return statusMissingCodeSnippet.Err()
	}
	if req.ContentHash == "" {
		return statusMissingContentHash.Err()
	}
	if req.RepositoryName == "" || req.RepositoryOwner == "" {
		return statusMissingRepositoryInfo.Err()
	}
	if req.RepositoryType == proto.CodeReference_REPOSITORY_TYPE_UNSPECIFIED {
		return statusInvalidRepositoryType.Err()
	}
	return nil
}

func validateUpdateCodeReferenceRequest(req *proto.UpdateCodeReferenceRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	if req.FilePath == "" {
		return statusMissingFilePath.Err()
	}
	if req.LineNumber <= 0 {
		return statusMissingLineNumber.Err()
	}
	if req.CodeSnippet == "" {
		return statusMissingCodeSnippet.Err()
	}
	if req.ContentHash == "" {
		return statusMissingContentHash.Err()
	}
	return nil
}

func validateDeleteCodeReferenceRequest(req *proto.DeleteCodeReferenceRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	return nil
}

func validateGetCodeReferenceRequest(req *proto.GetCodeReferenceRequest) error {
	if req.Id == "" {
		return statusMissingID.Err()
	}
	return nil
}

func validateListCodeReferencesRequest(req *proto.ListCodeReferencesRequest) error {
	if req.FeatureId == "" {
		return statusMissingFeatureID.Err()
	}
	return nil
}
