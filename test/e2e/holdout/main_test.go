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

// Package holdout holds the end-to-end tests for RFC 0048 "Universal Holdout
// Groups" (docs/rfcs/0048-universal-holdout-groups.md), phases 1 and 2.
//
// Isolation. A holdout group changes the evaluation of every target flag in
// its environment and bumps the UserEvaluationsID of every client there, so
// these tests never run in the shared e2e environment. TestMain creates a
// dedicated environment under the e2e project with its own client and server
// SDK keys, grants the env-editor and env-viewer accounts their usual roles
// on it, and archives it when the package finishes. Tests that start a group
// hold servingMu, because at most one group may serve per environment and
// the WAITING windows of other groups may not overlap it (RFC §4.5/§4.6).
// Control-plane tests that only need a WAITING group run in parallel on
// disjoint far-future windows handed out by farFutureWindow.
package holdout

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	gatewayclient "github.com/bucketeer-io/bucketeer/v2/pkg/api/client"
	auditlogclient "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/client"
	btclient "github.com/bucketeer-io/bucketeer/v2/pkg/batch/client"
	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	btproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
	"github.com/bucketeer-io/bucketeer/v2/test/e2e/util"
)

const (
	prefixTestName        = "e2e-test"
	timeout               = 60 * time.Second
	deadlockRetryAttempts = 3
	pollInterval          = time.Second

	// Dev cron (docker-compose/config/batch-jobs/crontab) runs the
	// HoldoutGroupCacher, FeatureFlagCacher, SegmentUserCacher and
	// HoldoutGroupStatusUpdater every minute, and the gateway's in-memory
	// caches expire after 1 minute (api --features-memory-cache-ttl). Holdout
	// writes also refresh the Redis snapshot synchronously and evict the
	// gateway L1 through domain events, so a start/stop is normally visible
	// within seconds. The worst case is one missed cron run plus one L1 TTL:
	// 1 min + 1 min, with 30 s of slack.
	cachePropagationTimeout = 150 * time.Second
	// A scheduled transition is taken by the next HoldoutGroupStatusUpdater
	// cron run (≤ 60 s away) plus the job's own duration.
	statusUpdaterTimeout = 90 * time.Second
	// Audit logs are written asynchronously by the subscriber from pubsub.
	auditLogTimeout = 120 * time.Second
	// Segment user uploads are processed asynchronously by the subscriber;
	// see test/e2e/feature/segment_user_test.go for the dev-container budget.
	segmentUploadTimeout = 3 * time.Minute
	// GetEvaluations can transiently fail with codes.Internal while the
	// feature cache snapshot is being rebuilt; outlast the 1 minute L1 TTL.
	getEvaluationsRetryAttempts = 10
	getEvaluationsRetryInterval = 10 * time.Second
	// Serving groups get a one hour window; every serving test finishes well
	// inside it and stops the group explicitly.
	servingWindow = time.Hour
)

var (
	webGatewayAddr                 = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort                 = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert                 = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath                     = flag.String("api-key", "", "Client SDK API key for api-gateway")
	apiKeyServerPath               = flag.String("api-key-server", "", "Server SDK API key for api-gateway")
	gatewayAddr                    = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort                    = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert                    = flag.String("gateway-cert", "", "Gateway crt file")
	sysAdminAccessTokenPath        = flag.String("sys-admin-access-token", "", "System admin access token path")
	orgOwnerDefaultAccessTokenPath = flag.String("org-owner-default-access-token", "", "Organization admin access token path")
	orgOwnerE2EAccessTokenPath     = flag.String("org-owner-e2e-access-token", "", "Organization admin (e2e org) access token path")
	envEditorAccessTokenPath       = flag.String("env-editor-access-token", "", "Environment editor access token path")
	envViewerAccessTokenPath       = flag.String("env-viewer-access-token", "", "Environment viewer access token path")
	environmentID                  = flag.String("environment-id", "", "Environment id")
	organizationID                 = flag.String("organization-id", "", "Organization ID")
	testID                         = flag.String("test-id", "", "test ID")

	// holdoutEnv is the dedicated environment created by TestMain.
	holdoutEnv struct {
		id            string
		clientKeyPath string
		serverKeyPath string
	}
	// servingMu serializes every test that starts a group or creates a group
	// whose window is near now (RFC §4.6: one serving group per environment,
	// no overlapping WAITING windows).
	servingMu sync.Mutex
	// farFutureBase is a per-process random day offset so WAITING groups left
	// behind by a crashed run cannot collide with the next run's windows.
	farFutureBase = 400 + rand.Intn(1000) // nolint:gosec
	windowCounter atomic.Int64
)

func TestMain(m *testing.M) {
	flag.Parse()
	code, err := run(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "holdout e2e setup failed: %v\n", err)
		if code == 0 {
			code = 1
		}
	}
	os.Exit(code)
}

func run(m *testing.M) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	setup, err := newEnvironmentSetup(ctx)
	if err != nil {
		return 1, err
	}
	defer setup.cleanup()
	if err := setup.waitForGateway(ctx); err != nil {
		return 1, err
	}
	return m.Run(), nil
}

// environmentSetup owns the resources TestMain creates and reverses them in
// cleanup: the environment is archived, the role granted to the editor and
// viewer accounts on it is revoked, and the SDK keys are disabled.
type environmentSetup struct {
	envClient     environmentclient.Client
	accountClient accountclient.Client
	featureClient featureclient.Client
	keyIDs        []string
	// grantedTokens are the access tokens of the accounts that received a
	// role on the holdout environment; cleanup revokes exactly that role.
	grantedTokens []string
	tempFiles     []string
}

func newEnvironmentSetup(ctx context.Context) (*environmentSetup, error) {
	s := &environmentSetup{}
	var err error
	s.envClient, err = newEnvironmentClient(*orgOwnerDefaultAccessTokenPath)
	if err != nil {
		return nil, err
	}
	s.accountClient, err = newAccountClient(*orgOwnerDefaultAccessTokenPath)
	if err != nil {
		return nil, err
	}
	s.featureClient, err = newFeatureClientWithToken(*orgOwnerDefaultAccessTokenPath)
	if err != nil {
		return nil, err
	}
	// The holdout environment lives in the same project as the shared e2e
	// environment, so the org owner token keeps Organization_ADMIN there.
	shared, err := s.envClient.GetEnvironmentV2(ctx, &environmentproto.GetEnvironmentV2Request{Id: *environmentID})
	if err != nil {
		s.cleanup()
		return nil, fmt.Errorf("get shared environment %q: %w", *environmentID, err)
	}
	name, suffix, err := holdoutEnvironmentName()
	if err != nil {
		s.cleanup()
		return nil, err
	}
	created, err := s.envClient.CreateEnvironmentV2(ctx, &environmentproto.CreateEnvironmentV2Request{
		Name:        name,
		UrlCode:     "e2e-test-holdout-" + suffix,
		Description: "Dedicated environment for the holdout e2e suite. Archived by TestMain.",
		ProjectId:   shared.Environment.ProjectId,
	})
	if err != nil {
		s.cleanup()
		return nil, fmt.Errorf("create holdout environment: %w", err)
	}
	holdoutEnv.id = created.Environment.Id
	holdoutEnv.clientKeyPath, err = s.createAPIKey(ctx, accountproto.APIKey_SDK_CLIENT, "client")
	if err != nil {
		s.cleanup()
		return nil, err
	}
	holdoutEnv.serverKeyPath, err = s.createAPIKey(ctx, accountproto.APIKey_SDK_SERVER, "server")
	if err != nil {
		s.cleanup()
		return nil, err
	}
	for tokenPath, role := range map[string]accountproto.AccountV2_Role_Environment{
		*envEditorAccessTokenPath: accountproto.AccountV2_Role_Environment_EDITOR,
		*envViewerAccessTokenPath: accountproto.AccountV2_Role_Environment_VIEWER,
	} {
		if err := s.grantEnvironmentRole(ctx, tokenPath, role); err != nil {
			s.cleanup()
			return nil, err
		}
	}
	return s, nil
}

const (
	// Server-side caps (pkg/environment/api): name ≤ 50 bytes, url code
	// ^[a-z0-9-]{1,50}$.
	maxEnvironmentNameLength = 50
	// Fewest random hex characters accepted before the run is refused
	// rather than risking a name collision.
	minEnvironmentSuffixLength = 8
)

// holdoutEnvironmentName returns "<prefix>[-<test-id>]-holdout-env-<hex>",
// shrinking the random suffix so the name fits the server cap whatever the
// test ID is. hack/delete-e2e-data-* match the same "-holdout-env-" shape.
func holdoutEnvironmentName() (name, suffix string, err error) {
	base := prefixed("holdout-env", "")
	room := maxEnvironmentNameLength - len(base)
	if room < minEnvironmentSuffixLength {
		return "", "", fmt.Errorf(
			"test ID %q is too long: the holdout environment name %q leaves %d bytes for a suffix, need %d",
			*testID, base, room, minEnvironmentSuffixLength,
		)
	}
	suffix = strings.ReplaceAll(mustUUID(), "-", "")
	if room < len(suffix) {
		suffix = suffix[:room]
	}
	return base + suffix, suffix, nil
}

func (s *environmentSetup) createAPIKey(
	ctx context.Context,
	role accountproto.APIKey_Role,
	kind string,
) (string, error) {
	resp, err := s.accountClient.CreateAPIKey(ctx, &accountproto.CreateAPIKeyRequest{
		EnvironmentId: holdoutEnv.id,
		Name:          prefixed("holdout-"+kind+"-key", mustUUID()),
		Role:          role,
		Description:   "Holdout e2e suite SDK key",
	})
	if err != nil {
		return "", fmt.Errorf("create %s api key: %w", kind, err)
	}
	s.keyIDs = append(s.keyIDs, resp.ApiKey.Id)
	// The gateway client and the REST helpers read the key from a file.
	f, err := os.CreateTemp("", "holdout-e2e-"+kind+"-key-*")
	if err != nil {
		return "", err
	}
	s.tempFiles = append(s.tempFiles, f.Name())
	if _, err := f.WriteString(resp.ApiKey.ApiKey); err != nil {
		f.Close()
		return "", err
	}
	return f.Name(), f.Close()
}

// grantEnvironmentRole adds the holdout environment to the environment roles
// of the account behind tokenPath. The e2e editor and viewer accounts only
// carry a role on the shared e2e environment; the authorization tests need
// them to reach the holdout environment with the same roles.
func (s *environmentSetup) grantEnvironmentRole(
	ctx context.Context,
	tokenPath string,
	role accountproto.AccountV2_Role_Environment,
) error {
	if err := s.editEnvironmentRoles(ctx, tokenPath, func(roles []*accountproto.AccountV2_EnvironmentRole) []*accountproto.AccountV2_EnvironmentRole {
		return append(roles, &accountproto.AccountV2_EnvironmentRole{EnvironmentId: holdoutEnv.id, Role: role})
	}); err != nil {
		return err
	}
	s.grantedTokens = append(s.grantedTokens, tokenPath)
	return nil
}

// revokeEnvironmentRole removes only this run's holdout environment from the
// account's roles. The list is re-read at revoke time instead of restored
// from a snapshot, so a concurrent run that granted its own environment in
// the meantime keeps its grant.
func (s *environmentSetup) revokeEnvironmentRole(ctx context.Context, tokenPath string) error {
	return s.editEnvironmentRoles(ctx, tokenPath, func(roles []*accountproto.AccountV2_EnvironmentRole) []*accountproto.AccountV2_EnvironmentRole {
		kept := roles[:0]
		for _, r := range roles {
			if r.EnvironmentId != holdoutEnv.id {
				kept = append(kept, r)
			}
		}
		return kept
	})
}

// editEnvironmentRoles reads the account behind tokenPath, applies edit to
// its current environment roles and writes them back. The account domain
// requires a language on every update, so the account's own language is
// sent along unchanged ("en" only when the bootstrapped account has none).
func (s *environmentSetup) editEnvironmentRoles(
	ctx context.Context,
	tokenPath string,
	edit func([]*accountproto.AccountV2_EnvironmentRole) []*accountproto.AccountV2_EnvironmentRole,
) error {
	own, err := newAccountClient(tokenPath)
	if err != nil {
		return err
	}
	defer own.Close()
	me, err := own.GetMe(ctx, &accountproto.GetMeRequest{OrganizationId: *organizationID})
	if err != nil {
		return fmt.Errorf("get account for %s: %w", tokenPath, err)
	}
	current := make([]*accountproto.AccountV2_EnvironmentRole, 0, len(me.Account.EnvironmentRoles))
	for _, r := range me.Account.EnvironmentRoles {
		current = append(current, &accountproto.AccountV2_EnvironmentRole{
			EnvironmentId: r.Environment.Id,
			Role:          r.Role,
		})
	}
	language := me.Account.Language
	if language == "" {
		language = "en"
	}
	_, err = s.accountClient.UpdateAccountV2(ctx, &accountproto.UpdateAccountV2Request{
		Email:            me.Account.Email,
		OrganizationId:   *organizationID,
		Language:         wrapperspb.String(language),
		EnvironmentRoles: edit(current),
	})
	if err != nil {
		return fmt.Errorf("update environment roles of %s: %w", me.Account.Email, err)
	}
	return nil
}

// waitForGateway blocks until the api-gateway accepts both new SDK keys. The
// gateway falls back to the database when a key misses both caches, so this
// normally succeeds on the first attempt; the bound covers the ApiKeyCacher
// cron (1 min) plus the in-memory TTL (1 min) in case a negative result was
// cached.
func (s *environmentSetup) waitForGateway(ctx context.Context) error {
	deadline := time.Now().Add(cachePropagationTimeout)
	for {
		err := s.probeGateway(ctx)
		if err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("gateway did not accept the holdout environment keys: %w", err)
		}
		time.Sleep(pollInterval)
	}
}

func (s *environmentSetup) probeGateway(ctx context.Context) error {
	client, err := newGatewayClient(holdoutEnv.clientKeyPath)
	if err != nil {
		return err
	}
	defer client.Close()
	server, err := newGatewayClient(holdoutEnv.serverKeyPath)
	if err != nil {
		return err
	}
	defer server.Close()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if _, err := client.GetEvaluations(ctx, &gatewayproto.GetEvaluationsRequest{
		Tag:  "probe",
		User: &userproto.User{Id: "probe"},
	}); err != nil {
		return fmt.Errorf("client key: %w", err)
	}
	if _, err := server.GetFeatureFlags(ctx, &gatewayproto.GetFeatureFlagsRequest{
		SourceId:   eventproto.SourceId_GO_SERVER,
		SdkVersion: "v0.0.1-e2e-test",
	}); err != nil {
		return fmt.Errorf("server key: %w", err)
	}
	return nil
}

func (s *environmentSetup) cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if holdoutEnv.id != "" && s.featureClient != nil {
		s.sweepHoldoutGroups(ctx)
	}
	for _, tokenPath := range s.grantedTokens {
		if err := s.revokeEnvironmentRole(ctx, tokenPath); err != nil {
			fmt.Fprintf(os.Stderr, "holdout e2e cleanup: %v\n", err)
		}
	}
	for _, id := range s.keyIDs {
		if _, err := s.accountClient.UpdateAPIKey(ctx, &accountproto.UpdateAPIKeyRequest{
			Id:            id,
			EnvironmentId: holdoutEnv.id,
			Disabled:      wrapperspb.Bool(true),
		}); err != nil {
			fmt.Fprintf(os.Stderr, "holdout e2e cleanup: disable api key %s: %v\n", id, err)
		}
	}
	if holdoutEnv.id != "" {
		if _, err := s.envClient.ArchiveEnvironmentV2(ctx, &environmentproto.ArchiveEnvironmentV2Request{
			Id: holdoutEnv.id,
		}); err != nil {
			fmt.Fprintf(os.Stderr, "holdout e2e cleanup: archive environment %s: %v\n", holdoutEnv.id, err)
		}
	}
	for _, f := range s.tempFiles {
		os.Remove(f)
	}
	for _, c := range []interface{ Close() }{s.envClient, s.accountClient, s.featureClient} {
		if c != nil {
			c.Close()
		}
	}
}

// sweepHoldoutGroups stops and archives any group a failed test left serving
// and deletes leftover WAITING groups, so the archived environment holds no
// live holdout and the next run of the cleanup tooling finds only inert rows.
func (s *environmentSetup) sweepHoldoutGroups(ctx context.Context) {
	resp, err := s.featureClient.ListHoldoutGroups(ctx, &featureproto.ListHoldoutGroupsRequest{
		EnvironmentId: holdoutEnv.id,
		PageSize:      100,
		Statuses: []featureproto.HoldoutGroup_Status{
			featureproto.HoldoutGroup_WAITING,
			featureproto.HoldoutGroup_RUNNING,
			featureproto.HoldoutGroup_ANALYSIS,
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "holdout e2e cleanup: list holdout groups: %v\n", err)
		return
	}
	for _, g := range resp.HoldoutGroups {
		fmt.Fprintf(os.Stderr, "holdout e2e cleanup: sweeping leftover group %s (%s)\n", g.Id, g.Status)
		var err error
		switch g.Status {
		case featureproto.HoldoutGroup_WAITING:
			_, err = s.featureClient.DeleteHoldoutGroup(ctx, &featureproto.DeleteHoldoutGroupRequest{
				Id: g.Id, EnvironmentId: holdoutEnv.id,
			})
		default:
			_, err = s.featureClient.StopHoldoutGroup(ctx, &featureproto.StopHoldoutGroupRequest{
				Id: g.Id, EnvironmentId: holdoutEnv.id,
			})
			if err == nil {
				_, err = s.featureClient.ArchiveHoldoutGroup(ctx, &featureproto.ArchiveHoldoutGroupRequest{
					Id: g.Id, EnvironmentId: holdoutEnv.id,
				})
			}
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "holdout e2e cleanup: sweep group %s: %v\n", g.Id, err)
		}
	}
}

// ---- clients ----

func dialOptions(tokenPath string) ([]rpcclient.Option, error) {
	creds, err := rpcclient.NewPerRPCCredentials(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("credentials %s: %w", tokenPath, err)
	}
	return []rpcclient.Option{
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30 * time.Second),
		rpcclient.WithBlock(),
	}, nil
}

func webGatewayTarget() string {
	return fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort)
}

func newEnvironmentClient(tokenPath string) (environmentclient.Client, error) {
	opts, err := dialOptions(tokenPath)
	if err != nil {
		return nil, err
	}
	return environmentclient.NewClient(webGatewayTarget(), *webGatewayCert, opts...)
}

func newAccountClient(tokenPath string) (accountclient.Client, error) {
	opts, err := dialOptions(tokenPath)
	if err != nil {
		return nil, err
	}
	return accountclient.NewClient(webGatewayTarget(), *webGatewayCert, opts...)
}

func newFeatureClientWithToken(tokenPath string) (featureclient.Client, error) {
	opts, err := dialOptions(tokenPath)
	if err != nil {
		return nil, err
	}
	return featureclient.NewClient(webGatewayTarget(), *webGatewayCert, opts...)
}

// newFeatureClient returns a FeatureService client authenticated as the
// organization owner, which holds Organization_ADMIN on the holdout
// environment and may therefore call every holdout mutation.
func newFeatureClient(t *testing.T) featureclient.Client {
	t.Helper()
	return newFeatureClientAs(t, *orgOwnerDefaultAccessTokenPath)
}

func newFeatureClientAs(t *testing.T, tokenPath string) featureclient.Client {
	t.Helper()
	c, err := newFeatureClientWithToken(tokenPath)
	require.NoError(t, err)
	t.Cleanup(func() { c.Close() })
	return c
}

func newBatchClient(t *testing.T) btclient.Client {
	t.Helper()
	opts, err := dialOptions(*orgOwnerDefaultAccessTokenPath)
	require.NoError(t, err)
	c, err := btclient.NewClient(webGatewayTarget(), *webGatewayCert, opts...)
	require.NoError(t, err)
	t.Cleanup(func() { c.Close() })
	return c
}

func newAuditLogClient(t *testing.T) auditlogclient.Client {
	t.Helper()
	opts, err := dialOptions(*orgOwnerDefaultAccessTokenPath)
	require.NoError(t, err)
	c, err := auditlogclient.NewClient(webGatewayTarget(), *webGatewayCert, opts...)
	require.NoError(t, err)
	t.Cleanup(func() { c.Close() })
	return c
}

func newGatewayClient(apiKeyPath string) (gatewayclient.Client, error) {
	creds, err := gatewayclient.NewPerRPCCredentials(apiKeyPath)
	if err != nil {
		return nil, err
	}
	return gatewayclient.NewClient(
		fmt.Sprintf("%s:%d", *gatewayAddr, *gatewayPort),
		*gatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
}

// ---- naming ----

func mustUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return id.String()
}

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	require.NoError(t, err)
	return id.String()
}

// prefixed builds "e2e-test[-<test-id>]-<kind>-<id>" so the e2e cleanup
// tooling (hack/delete-e2e-data-*) finds every row this suite creates.
func prefixed(kind, id string) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-%s-%s", prefixTestName, *testID, kind, id)
	}
	return fmt.Sprintf("%s-%s-%s", prefixTestName, kind, id)
}

func newFeatureID(t *testing.T) string {
	t.Helper()
	return prefixed("feature-id", newUUID(t))
}

func newTag(t *testing.T) string {
	t.Helper()
	return prefixed("tag", newUUID(t))
}

func newUserID(t *testing.T) string {
	t.Helper()
	return prefixed("user", newUUID(t))
}

func newGroupName(t *testing.T) string {
	t.Helper()
	return prefixed("holdout", newUUID(t))
}

// ---- time windows ----

// farFutureWindow returns a one-day window at least 400 days from now, on a
// 10-day grid indexed by a package counter, so parallel control-plane tests
// never trip the overlap validation of RFC §4.5 against each other.
func farFutureWindow() (startAt, stopAt int64) {
	i := windowCounter.Add(1)
	start := time.Now().AddDate(0, 0, farFutureBase+int(i)*10).Truncate(24 * time.Hour)
	return start.Unix(), start.Add(24 * time.Hour).Unix()
}

// servingWindowNow returns a window for a group the test starts manually:
// [now+servingWindow, now+2*servingWindow). start_at is kept in the future
// so the HoldoutGroupStatusUpdater cron (every minute) cannot start the group
// first; the manual Start clamps start_at to now (RFC §4.6). Callers must
// hold servingMu.
func servingWindowNow() (startAt, stopAt int64) {
	now := time.Now()
	return now.Add(servingWindow).Unix(), now.Add(2 * servingWindow).Unix()
}

// waitUntilUnix blocks until the wall clock is strictly past ts. It is used
// to let updated_at (whole seconds) advance between two mutations and to let
// a scheduled timestamp pass; the wait is bounded by the distance to ts.
func waitUntilUnix(t *testing.T, ts int64) {
	t.Helper()
	target := time.Unix(ts+1, 0)
	if d := time.Until(target); d > 0 {
		require.Less(t, d, 2*time.Minute, "waitUntilUnix: target too far in the future")
		time.Sleep(d)
	}
}

// ---- generic helpers ----

// withDeadlockRetry runs a write, retrying MySQL deadlocks (Error 1213) like
// the feature suite does.
func withDeadlockRetry(t *testing.T, name string, fn func(ctx context.Context) error) error {
	t.Helper()
	var err error
	for i := 0; i < deadlockRetryAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		err = fn(ctx)
		cancel()
		if err == nil || !util.IsDeadlockError(err) || i == deadlockRetryAttempts-1 {
			return err
		}
		t.Logf("Retrying %s (attempt %d/%d): %v", name, i+1, deadlockRetryAttempts, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return err
}

func requireCode(t *testing.T, err error, want codes.Code, msgAndArgs ...interface{}) {
	t.Helper()
	require.Error(t, err, msgAndArgs...)
	st, ok := status.FromError(err)
	require.True(t, ok, "not a gRPC status: %v", err)
	require.Equalf(t, want, st.Code(), "unexpected code for %v (%v)", err, msgAndArgs)
}

func runBatchJob(t *testing.T, job btproto.BatchJob) {
	t.Helper()
	client := newBatchClient(t)
	var lastErr error
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		_, lastErr = client.ExecuteBatchJob(ctx, &btproto.BatchJobRequest{Job: job})
		cancel()
		if lastErr == nil {
			return
		}
		if status.Code(lastErr) == codes.Internal {
			break
		}
		time.Sleep(2 * time.Second)
	}
	require.NoError(t, lastErr, "ExecuteBatchJob %s", job)
}
