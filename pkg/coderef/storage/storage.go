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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package storage

import (
	"context"
	"errors"

	"github.com/bucketeer-io/bucketeer/v2/pkg/coderef/domain"
	bkterr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	coderefproto "github.com/bucketeer-io/bucketeer/v2/proto/coderef"
)

var (
	ErrCodeReferenceNotFound = bkterr.NewErrorNotFound(
		bkterr.CoderefPackageName,
		"code reference not found", "code_reference",
	)
	ErrCodeReferenceUnexpectedAffectedRows = bkterr.NewErrorUnexpectedAffectedRows(
		bkterr.CoderefPackageName,
		"code reference unexpected affected rows",
	)
)

// Shared list-query errors returned by CodeReferenceStorage implementations.
var (
	ErrInvalidOrderBy = errors.New("coderef/storage: invalid order by")
	ErrInvalidCursor  = errors.New("coderef/storage: invalid cursor")
)

type CodeReferenceStorage interface {
	CreateCodeReference(ctx context.Context, codeRef *domain.CodeReference) error
	UpdateCodeReference(ctx context.Context, codeRef *domain.CodeReference) error
	GetCodeReference(ctx context.Context, id string) (*domain.CodeReference, error)
	ListCodeReferences(
		ctx context.Context,
		params ListCodeReferencesParams,
	) ([]*domain.CodeReference, int, int64, error)
	DeleteCodeReference(ctx context.Context, id string) error
	// GetCodeReferenceCountsByFeatureIDs returns a map of feature ID to code reference count
	// for the given environment. This is used for bulk archivability evaluation.
	GetCodeReferenceCountsByFeatureIDs(
		ctx context.Context,
		environmentID string,
		featureIDs []string,
	) (map[string]int64, error)
}

type ListCodeReferencesParams struct {
	EnvironmentID    string
	FeatureID        string
	RepositoryName   string
	RepositoryOwner  string
	RepositoryType   coderefproto.CodeReference_RepositoryType
	RepositoryBranch string
	FileExtension    string
	OrderBy          coderefproto.ListCodeReferencesRequest_OrderBy
	OrderDirection   coderefproto.ListCodeReferencesRequest_OrderDirection
	PageSize         int
	Cursor           string
}
