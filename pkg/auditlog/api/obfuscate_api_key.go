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

// obfuscateAPIKeys obfuscates the API keys of the audit logs saved before they were obfuscated
// at creation. The rows saved since then are already obfuscated.
func (s *auditlogService) obfuscateAPIKeys(auditlogs []*proto.AuditLog) {
	for i := range auditlogs {
		s.obfuscateAPIKey(auditlogs[i])
	}
}

func (s *auditlogService) obfuscateAPIKey(auditlog *proto.AuditLog) {
	if err := domain.ObfuscateAPIKey(auditlog); err != nil {
		s.logger.Error("Failed to obfuscate the api key of the audit log",
			zap.Error(err),
			zap.String("id", auditlog.Id),
		)
	}
}
