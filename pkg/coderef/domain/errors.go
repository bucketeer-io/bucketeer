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

package domain

import "errors"

var (
	ErrInvalidFeatureID        = errors.New("coderef: invalid feature id")
	ErrInvalidFilePath         = errors.New("coderef: invalid file path")
	ErrInvalidLineNumber       = errors.New("coderef: invalid line number")
	ErrInvalidCodeSnippet      = errors.New("coderef: invalid code snippet")
	ErrInvalidContentHash      = errors.New("coderef: invalid content hash")
	ErrInvalidRepositoryName   = errors.New("coderef: invalid repository name")
	ErrInvalidRepositoryOwner  = errors.New("coderef: invalid repository owner")
	ErrInvalidRepositoryType   = errors.New("coderef: invalid repository type")
	ErrInvalidRepositoryBranch = errors.New("coderef: invalid repository branch")
	ErrInvalidCommitHash       = errors.New("coderef: invalid commit hash")
	ErrInvalidEnvironmentID    = errors.New("coderef: invalid environment id")
)
