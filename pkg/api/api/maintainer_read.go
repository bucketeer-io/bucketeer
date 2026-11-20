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
	account *accountproto.AccountV2
	logged  atomic.Bool
}

// reportMaintainerRead counts a public API read that would be denied as the
// API key's maintainer (RFC 0050 step 11). A cached maintainer is judged in
// place; otherwise the lookup runs once per key off the response path and
// counts the read that started it.
func (s *grpcGatewayService) reportMaintainerRead(
	ctx context.Context,
	envAPIKey *accountproto.EnvironmentAPIKey,
) {
	method := methodName(ctx)
	scope, ok := maintainerReads[method]
	if !ok {
		return
	}
	key := maintainerAccountCacheKey(envAPIKey)
	if entry := s.cachedMaintainerAccount(key); entry != nil {
		s.countMaintainerRead(ctx, envAPIKey, entry, method, scope)
		return
	}
	// The result is not awaited: singleflight runs one lookup per key and the
	// reads that arrive meanwhile are left to the next one.
	s.flightgroup.DoChan(key, func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), singleflightFetchTimeout)
		defer cancel()
		// A lookup that finished while this read was on its way already filled the cache.
		entry := s.cachedMaintainerAccount(key)
		if entry == nil {
			var err error
			if entry, err = s.lookupMaintainerAccount(ctx, envAPIKey); err != nil {
				s.logger.Warn("Failed to check the API key maintainer",
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("apiKey", obfuscateString(envAPIKey.ApiKey.Id, obfuscateAPIKeyLength)),
						zap.String("method", method),
					)...,
				)
				return nil, err
			}
			if err := s.maintainerAccountCache.Put(key, entry, maintainerReadTTL); err != nil {
				return nil, err
			}
		}
		s.countMaintainerRead(ctx, envAPIKey, entry, method, scope)
		return nil, nil
	})
}

// maintainerAccountCacheKey names the key's maintainer, so a changed
// maintainer is looked up again and each key warns on its own.
func maintainerAccountCacheKey(envAPIKey *accountproto.EnvironmentAPIKey) string {
	return "maintainer-read:" + envAPIKey.Environment.OrganizationId + ":" +
		envAPIKey.ApiKey.Maintainer + ":" + envAPIKey.ApiKey.Id
}

func (s *grpcGatewayService) countMaintainerRead(
	ctx context.Context,
	envAPIKey *accountproto.EnvironmentAPIKey,
	entry *maintainerAccount,
	method string,
	scope maintainerReadScope,
) {
	reason := maintainerReadDenial(entry.account, scope, envAPIKey.Environment.Id)
	if reason == "" {
		return
	}
	maintainerReadDeniedCounter.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.Environment.Id, method, reason,
	).Inc()
	if entry.logged.CompareAndSwap(false, true) {
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

func (s *grpcGatewayService) lookupMaintainerAccount(
	ctx context.Context,
	envAPIKey *accountproto.EnvironmentAPIKey,
) (*maintainerAccount, error) {
	resp, err := s.accountClient.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
		Email:          envAPIKey.ApiKey.Maintainer,
		OrganizationId: envAPIKey.Environment.OrganizationId,
	})
	switch status.Code(err) {
	case codes.OK:
		return &maintainerAccount{account: resp.Account}, nil
	case codes.NotFound, codes.InvalidArgument:
		// An empty or malformed maintainer has no account either.
		return &maintainerAccount{}, nil
	default:
		return nil, err
	}
}

func (s *grpcGatewayService) cachedMaintainerAccount(key string) *maintainerAccount {
	cached, err := s.maintainerAccountCache.Get(key)
	if err != nil {
		return nil
	}
	entry, _ := cached.(*maintainerAccount)
	return entry
}

func methodName(ctx context.Context) string {
	fullMethod, _ := grpc.Method(ctx)
	return fullMethod[strings.LastIndex(fullMethod, "/")+1:]
}
