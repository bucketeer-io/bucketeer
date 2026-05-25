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
	ErrAdminAuditLogAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.AuditlogPackageName, "admin auditlog already exists")
)

// ListAdminAuditLogsParams carries list intent for ListAdminAuditLogs without database-specific types.
type ListAdminAuditLogsParams struct {
	EntityType     *int32
	From           int64
	To             int64
	SearchKeyword  string
	OrderBy        proto.ListAdminAuditLogsRequest_OrderBy
	OrderDirection proto.ListAdminAuditLogsRequest_OrderDirection
	PageSize       int
	Cursor         string
}

type AdminAuditLogStorage interface {
	CreateAdminAuditLogs(ctx context.Context, auditLogs []*domain.AuditLog) error
	CreateAdminAuditLog(ctx context.Context, auditLog *domain.AuditLog) error
	ListAdminAuditLogs(
		ctx context.Context,
		params ListAdminAuditLogsParams,
	) ([]*proto.AuditLog, int, int64, error)
}
