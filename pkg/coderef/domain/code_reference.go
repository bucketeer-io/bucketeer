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

	proto "github.com/bucketeer-io/bucketeer/proto/coderef"
)

type CodeReference struct {
	Id               string
	FeatureId        string
	FilePath         string
	LineNumber       int32
	CodeSnippet      string
	ContentHash      string
	Aliases          []string
	RepositoryName   string
	RepositoryOwner  string
	RepositoryType   proto.CodeReference_RepositoryType
	RepositoryBranch string
	CommitHash       string
	EnvironmentId    string
	CreatedAt        int64
	UpdatedAt        int64
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
	}
}

func NewCodeReferenceWithTime(
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
	now time.Time,
) *CodeReference {
	timestamp := now.Unix()
	return &CodeReference{
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
		CreatedAt:        timestamp,
		UpdatedAt:        timestamp,
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
) {
	c.FilePath = filePath
	c.LineNumber = lineNumber
	c.CodeSnippet = codeSnippet
	c.ContentHash = contentHash
	c.Aliases = aliases
	c.RepositoryBranch = repositoryBranch
	c.CommitHash = commitHash
	c.UpdatedAt = time.Now().Unix()
}

func (c *CodeReference) ToProto() *proto.CodeReference {
	return &proto.CodeReference{
		Id:               c.Id,
		FeatureId:        c.FeatureId,
		FilePath:         c.FilePath,
		LineNumber:       c.LineNumber,
		CodeSnippet:      c.CodeSnippet,
		ContentHash:      c.ContentHash,
		Aliases:          c.Aliases,
		RepositoryName:   c.RepositoryName,
		RepositoryOwner:  c.RepositoryOwner,
		RepositoryType:   c.RepositoryType,
		RepositoryBranch: c.RepositoryBranch,
		CommitHash:       c.CommitHash,
		EnvironmentId:    c.EnvironmentId,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}
