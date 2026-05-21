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
package v2

import (
	"context"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/domain"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
)

var (
	ErrAuditLogAlreadyExists = pkgErr.NewErrorAlreadyExists(pkgErr.AuditlogPackageName, "auditlog already exists")
	ErrAuditLogNotFound      = pkgErr.NewErrorNotFound(pkgErr.AuditlogPackageName, "auditlog not found", "auditlog")
)

// ListAuditLogsParams carries list intent for ListAuditLogs without database-specific types.
type ListAuditLogsParams struct {
	EnvironmentID  string
	EntityType     *int32
	EntityID       string
	From           int64
	To             int64
	SearchKeyword  string
	OrderBy        proto.ListAuditLogsRequest_OrderBy
	OrderDirection proto.ListAuditLogsRequest_OrderDirection
	PageSize       int
	Cursor         string
}

type AuditLogStorage interface {
	GetAuditLog(ctx context.Context, id string, environmentID string) (*proto.AuditLog, error)
	CreateAuditLogs(ctx context.Context, auditLogs []*domain.AuditLog) error
	CreateAuditLog(ctx context.Context, auditLog *domain.AuditLog) error
	ListAuditLogs(
		ctx context.Context,
		params ListAuditLogsParams,
	) ([]*proto.AuditLog, int, int64, error)
}
