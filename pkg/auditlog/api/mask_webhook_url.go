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
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
)

// maskWebhookURLs masks the webhook URLs of the subscription audit logs saved before they
// were masked at creation. The rows saved since then are already masked.
func (s *auditlogService) maskWebhookURLs(auditlogs []*proto.AuditLog) {
	for i := range auditlogs {
		s.maskWebhookURL(auditlogs[i])
	}
}

func (s *auditlogService) maskWebhookURL(auditlog *proto.AuditLog) {
	if err := domain.MaskWebhookURL(auditlog); err != nil {
		s.logger.Error("Failed to mask the webhook url of the audit log",
			zap.Error(err),
			zap.String("id", auditlog.Id),
		)
	}
}
