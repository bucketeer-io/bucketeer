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

package domain

import (
	"time"

	"github.com/jinzhu/copier"

	proto "github.com/bucketeer-io/bucketeer/proto/coderef"
)

type CodeReference struct {
	proto.CodeReference
}

func NewCodeReference(
	id string,
	featureId string,
	filePath string,
	lineNumber int32,
	codeSnippet string,
	contentHash string,
	aliases []string,
	repositoryName string,
	repositoryOwner string,
	repositoryType proto.CodeReference_RepositoryType,
	repositoryBranch string,
	commitHash string,
	environmentId string,
) *CodeReference {
	now := time.Now().Unix()
	return &CodeReference{
		CodeReference: proto.CodeReference{
			Id:               id,
			FeatureId:        featureId,
			FilePath:         filePath,
			LineNumber:       lineNumber,
			CodeSnippet:      codeSnippet,
			ContentHash:      contentHash,
			Aliases:          aliases,
			RepositoryName:   repositoryName,
			RepositoryOwner:  repositoryOwner,
			RepositoryType:   repositoryType,
			RepositoryBranch: repositoryBranch,
			CommitHash:       commitHash,
			EnvironmentId:    environmentId,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
	}
}

func (c *CodeReference) Update(
	filePath string,
	lineNumber int32,
	codeSnippet string,
	contentHash string,
	aliases []string,
	repositoryBranch string,
	commitHash string,
) (*CodeReference, error) {
	updated := &CodeReference{}
	if err := copier.Copy(updated, c); err != nil {
		return nil, err
	}
	if filePath != "" {
		updated.FilePath = filePath
	}
	if lineNumber != 0 {
		updated.LineNumber = lineNumber
	}
	if codeSnippet != "" {
		updated.CodeSnippet = codeSnippet
	}
	if contentHash != "" {
		updated.ContentHash = contentHash
	}
	if len(aliases) > 0 {
		updated.Aliases = aliases
	}
	if repositoryBranch != "" {
		updated.RepositoryBranch = repositoryBranch
	}
	if commitHash != "" {
		updated.CommitHash = commitHash
	}
	updated.UpdatedAt = time.Now().Unix()
	return updated, nil
}
