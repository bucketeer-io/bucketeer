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
	"context"
	"strings"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

// maintainerReadScope is what a public API read will be judged against once
// it acts as the API key's maintainer (RFC 0050 step 12).
type maintainerReadScope int

const (
	maintainerReadEnvironment maintainerReadScope = iota
	maintainerReadOrganization

	maintainerReadTTL = 5 * time.Minute

	maintainerReadReasonNoAccount         = "no_account"
	maintainerReadReasonDisabled          = "disabled"
	maintainerReadReasonNoEnvironmentRole = "no_environment_role"
)

// maintainerReads are the public API reads relayed with the service token
// alone, so the maintainer is not checked today. The reads that carry the
// maintainer metadata are already judged as the maintainer.
var maintainerReads = map[string]maintainerReadScope{
	"GetAuditLog":                  maintainerReadEnvironment,
	"ListAuditLogs":                maintainerReadEnvironment,
	"ListFeatureHistory":           maintainerReadEnvironment,
	"GetAutoOpsRule":               maintainerReadEnvironment,
	"ListAutoOpsRules":             maintainerReadEnvironment,
	"GetCodeReference":             maintainerReadEnvironment,
	"ListCodeReferences":           maintainerReadEnvironment,
	"GetEvaluationTimeseriesCount": maintainerReadEnvironment,
	"GetExperimentEvaluationCount": maintainerReadEnvironment,
	"GetExperimentGoalCount":       maintainerReadEnvironment,
	"GetExperimentResult":          maintainerReadEnvironment,
	"GetOpsEvaluationUserCount":    maintainerReadEnvironment,
	"GetOpsGoalUserCount":          maintainerReadEnvironment,
	"ListExperimentResults":        maintainerReadEnvironment,
	"GetExperiment":                maintainerReadEnvironment,
	"ListExperiments":              maintainerReadEnvironment,
	"DebugEvaluateFeatures":        maintainerReadEnvironment,
	"GetFeature":                   maintainerReadEnvironment,
	"ListFeatures":                 maintainerReadEnvironment,
	"GetFlagTrigger":               maintainerReadEnvironment,
	"ListFlagTriggers":             maintainerReadEnvironment,
	"GetGoal":                      maintainerReadEnvironment,
	"ListGoals":                    maintainerReadEnvironment,
	"GetProgressiveRollout":        maintainerReadEnvironment,
	"ListProgressiveRollouts":      maintainerReadEnvironment,
	"GetPush":                      maintainerReadEnvironment,
	"ListPushes":                   maintainerReadEnvironment,
	"GetSubscription":              maintainerReadEnvironment,
	"GetSegment":                   maintainerReadEnvironment,
	"ListSegments":                 maintainerReadEnvironment,
	"GetProject":                   maintainerReadOrganization,
	"ListProjects":                 maintainerReadOrganization,
	"GetEnvironmentV2":             maintainerReadOrganization,
	"ListEnvironmentsV2":           maintainerReadOrganization,
}

// maintainerAccount is the maintainer's account in the key's organization, nil when there is none.
type maintainerAccount struct {
	account   *accountproto.AccountV2
	expiresAt time.Time
	logged    atomic.Bool
}

// reportMaintainerRead counts a public API read that would be denied as the
// API key's maintainer (RFC 0050 step 11). It runs off the response path;
// the lookup is bounded by singleflightFetchTimeout.
func (s *grpcGatewayService) reportMaintainerRead(
	ctx context.Context,
	envAPIKey *accountproto.EnvironmentAPIKey,
) {
	method := methodName(ctx)
	scope, ok := maintainerReads[method]
	if !ok {
		return
	}
	cached, err := s.getMaintainerAccount(ctx, envAPIKey)
	if err != nil {
		s.logger.Warn("Failed to check the API key maintainer",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("apiKey", obfuscateString(envAPIKey.ApiKey.Id, obfuscateAPIKeyLength)),
				zap.String("method", method),
			)...,
		)
		return
	}
	reason := maintainerReadDenial(cached.account, scope, envAPIKey.Environment.Id)
	if reason == "" {
		return
	}
	maintainerReadDeniedCounter.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.Environment.Id, method, reason,
	).Inc()
	if cached.logged.CompareAndSwap(false, true) {
		s.logger.Warn("Public API read would be denied as the API key maintainer",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("apiKey", obfuscateString(envAPIKey.ApiKey.Id, obfuscateAPIKeyLength)),
				zap.String("maintainer", envAPIKey.ApiKey.Maintainer),
				zap.String("organizationId", envAPIKey.Environment.OrganizationId),
				zap.String("environmentId", envAPIKey.Environment.Id),
				zap.String("method", method),
				zap.String("reason", reason),
			)...,
		)
	}
}

// maintainerReadDenial is why the maintainer would be denied the read, or
// empty when it would pass: the same membership and role checks that
// pkg/role applies to the maintainer on writes.
func maintainerReadDenial(
	account *accountproto.AccountV2,
	scope maintainerReadScope,
	environmentID string,
) string {
	if account == nil {
		return maintainerReadReasonNoAccount
	}
	if account.Disabled {
		return maintainerReadReasonDisabled
	}
	if scope == maintainerReadOrganization ||
		account.OrganizationRole >= accountproto.AccountV2_Role_Organization_ADMIN {
		return ""
	}
	for _, role := range account.EnvironmentRoles {
		if role.EnvironmentId == environmentID && role.Role >= accountproto.AccountV2_Role_Environment_VIEWER {
			return ""
		}
	}
	return maintainerReadReasonNoEnvironmentRole
}

func (s *grpcGatewayService) getMaintainerAccount(
	ctx context.Context,
	envAPIKey *accountproto.EnvironmentAPIKey,
) (*maintainerAccount, error) {
	keyID := envAPIKey.ApiKey.Id
	if entry := s.cachedMaintainerAccount(keyID); entry != nil {
		return entry, nil
	}
	v, err := s.singleflightFetch(ctx, "maintainer-read:"+keyID, func(ctx context.Context) (interface{}, error) {
		// A fetch that finished while this caller was waiting already filled the cache.
		if entry := s.cachedMaintainerAccount(keyID); entry != nil {
			return entry, nil
		}
		resp, err := s.accountClient.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
			Email:          envAPIKey.ApiKey.Maintainer,
			OrganizationId: envAPIKey.Environment.OrganizationId,
		})
		entry := &maintainerAccount{expiresAt: time.Now().Add(maintainerReadTTL)}
		switch status.Code(err) {
		case codes.OK:
			entry.account = resp.Account
		case codes.NotFound, codes.InvalidArgument:
			// An empty or malformed maintainer has no account either.
		default:
			return nil, err
		}
		s.maintainerAccountCache.Store(keyID, entry)
		s.evictExpiredMaintainerAccounts()
		return entry, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*maintainerAccount), nil
}

func (s *grpcGatewayService) cachedMaintainerAccount(keyID string) *maintainerAccount {
	cached, ok := s.maintainerAccountCache.Load(keyID)
	if !ok {
		return nil
	}
	if entry := cached.(*maintainerAccount); time.Now().Before(entry.expiresAt) {
		return entry
	}
	return nil
}

// evictExpiredMaintainerAccounts drops the expired entries, so the cache
// follows the keys in use. An entry replaced meanwhile is left alone.
func (s *grpcGatewayService) evictExpiredMaintainerAccounts() {
	now := time.Now()
	s.maintainerAccountCache.Range(func(key, value any) bool {
		if entry := value.(*maintainerAccount); !now.Before(entry.expiresAt) {
			s.maintainerAccountCache.CompareAndDelete(key, value)
		}
		return true
	})
}

func methodName(ctx context.Context) string {
	fullMethod, _ := grpc.Method(ctx)
	return fullMethod[strings.LastIndex(fullMethod, "/")+1:]
}
