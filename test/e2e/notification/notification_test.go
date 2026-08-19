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

package notification

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

const timeout = 60 * time.Second

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
)

// TestAdminNotificationLifecycle drives a draft notification through its full
// lifecycle as the system admin: create, update, get (editor view), publish,
// then verify it becomes visible to a regular console user, can be marked as
// read, and is finally deleted.
func TestAdminNotificationLifecycle(t *testing.T) {
	requireAccessToken(t, *sysAdminAccessTokenPath, "system admin")
	requireAccessToken(t, *orgOwnerE2EAccessTokenPath, "org owner (e2e org)")
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	admin := newNotificationClient(t, *sysAdminAccessTokenPath)
	viewer := newNotificationClient(t, *orgOwnerE2EAccessTokenPath)

	// Create a draft with two localizations.
	title := fmt.Sprintf("e2e-notification-%s", uniqueSuffix())
	createResp, err := admin.CreateAdminNotification(ctx, &proto.CreateAdminNotificationRequest{
		Localizations: []*proto.NotificationLocalization{
			{Language: "en", Title: title, Content: "English content"},
			{Language: "ja", Title: title + "-ja", Content: "日本語のコンテンツ"},
		},
	})
	require.NoError(t, err)
	notificationID := createResp.Notification.Id
	assert.Equal(t, proto.Notification_DRAFT, createResp.Notification.Status)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), timeout)
		defer cleanupCancel()
		_, _ = admin.DeleteAdminNotification(cleanupCtx, &proto.DeleteAdminNotificationRequest{Id: notificationID})
	})

	// The draft must not be visible to a non-admin console user yet.
	_, err = viewer.GetNotification(ctx, &proto.GetNotificationRequest{Id: notificationID})
	requireStatusCode(t, err, codes.NotFound, "viewer GetNotification on draft")

	// The draft must be visible in the admin's draft listing.
	listDraftsResp, err := admin.ListDraftAdminNotifications(ctx, &proto.ListDraftAdminNotificationsRequest{
		PageSize:      0,
		SearchKeyword: title,
	})
	require.NoError(t, err)
	assert.True(t, containsNotificationID(listDraftsResp.Notifications, notificationID),
		"expected draft %q to be in ListDraftAdminNotifications", notificationID)

	// Update the draft's localizations.
	updatedTitle := title + "-updated"
	updateResp, err := admin.UpdateAdminNotification(ctx, &proto.UpdateAdminNotificationRequest{
		Id: notificationID,
		Localizations: []*proto.NotificationLocalization{
			{Language: "en", Title: updatedTitle, Content: "Updated English content"},
		},
	})
	require.NoError(t, err)
	assert.Len(t, updateResp.Notification.Localizations, 1)
	assert.Equal(t, updatedTitle, updateResp.Notification.Localizations[0].Title)

	// The admin can read back all localizations of the draft.
	getResp, err := admin.GetNotification(ctx, &proto.GetNotificationRequest{Id: notificationID, Language: "en"})
	require.NoError(t, err)
	assert.Equal(t, proto.Notification_DRAFT, getResp.Notification.Status)
	assert.Equal(t, updatedTitle, getResp.Notification.Localization.Title)

	// Publish the draft.
	publishResp, err := admin.PublishAdminNotification(ctx, &proto.PublishAdminNotificationRequest{Id: notificationID})
	require.NoError(t, err)
	assert.Equal(t, proto.Notification_PUBLISHED, publishResp.Notification.Status)
	assert.NotZero(t, publishResp.Notification.PublishedAt)

	// Publishing twice must fail with FailedPrecondition.
	_, err = admin.PublishAdminNotification(ctx, &proto.PublishAdminNotificationRequest{Id: notificationID})
	requireStatusCode(t, err, codes.FailedPrecondition, "re-publish already published notification")

	// A regular console user can now see it via Get, resolved to its language,
	// without the full localizations list, and unread.
	viewerGetResp, err := viewer.GetNotification(ctx, &proto.GetNotificationRequest{Id: notificationID, Language: "en"})
	require.NoError(t, err)
	assert.Equal(t, updatedTitle, viewerGetResp.Notification.Localization.Title)
	assert.Empty(t, viewerGetResp.Notification.Localizations, "non-admin viewers must not see the full localization list")
	assert.False(t, viewerGetResp.Notification.Read)

	// It also shows up in the viewer's published list.
	listResp, err := viewer.ListNotifications(ctx, &proto.ListNotificationsRequest{
		SearchKeyword: updatedTitle,
		Language:      "en",
	})
	require.NoError(t, err)
	assert.True(t, containsNotificationID(listResp.Notifications, notificationID),
		"expected published notification %q to be in the viewer's ListNotifications", notificationID)

	// Mark it as read and confirm the unread count drops accordingly.
	_, err = viewer.MarkNotificationsAsRead(ctx, &proto.MarkNotificationsAsReadRequest{Ids: []string{notificationID}})
	require.NoError(t, err)

	viewerGetResp, err = viewer.GetNotification(ctx, &proto.GetNotificationRequest{Id: notificationID, Language: "en"})
	require.NoError(t, err)
	assert.True(t, viewerGetResp.Notification.Read)

	// MarkAllNotificationsAsRead and GetNotificationUnreadCount are callable by
	// any authenticated user and must not error.
	_, err = viewer.MarkAllNotificationsAsRead(ctx, &proto.MarkAllNotificationsAsReadRequest{})
	require.NoError(t, err)

	unreadResp, err := viewer.GetNotificationUnreadCount(ctx, &proto.GetNotificationUnreadCountRequest{})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, unreadResp.Count, int64(0))

	// Delete the notification; a second delete must report NotFound.
	_, err = admin.DeleteAdminNotification(ctx, &proto.DeleteAdminNotificationRequest{Id: notificationID})
	require.NoError(t, err)

	_, err = admin.DeleteAdminNotification(ctx, &proto.DeleteAdminNotificationRequest{Id: notificationID})
	requireStatusCode(t, err, codes.NotFound, "delete already deleted notification")
}

// TestCreateAdminNotificationValidation checks the request validation rules
// enforced by CreateAdminNotification.
func TestCreateAdminNotificationValidation(t *testing.T) {
	requireAccessToken(t, *sysAdminAccessTokenPath, "system admin")
	t.Parallel()

	admin := newNotificationClient(t, *sysAdminAccessTokenPath)

	t.Run("no localizations", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := admin.CreateAdminNotification(ctx, &proto.CreateAdminNotificationRequest{})
		requireStatusCode(t, err, codes.InvalidArgument, "create with no localizations")
	})

	t.Run("duplicated language", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := admin.CreateAdminNotification(ctx, &proto.CreateAdminNotificationRequest{
			Localizations: []*proto.NotificationLocalization{
				{Language: "en", Title: "t1", Content: "c1"},
				{Language: "en", Title: "t2", Content: "c2"},
			},
		})
		requireStatusCode(t, err, codes.InvalidArgument, "create with duplicated language")
	})

	t.Run("missing title", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := admin.CreateAdminNotification(ctx, &proto.CreateAdminNotificationRequest{
			Localizations: []*proto.NotificationLocalization{
				{Language: "en", Content: "c1"},
			},
		})
		requireStatusCode(t, err, codes.InvalidArgument, "create with missing title")
	})
}

// TestAdminNotificationPermissions verifies that every system-admin-only RPC
// rejects non-system-admin callers with PermissionDenied/Unauthenticated.
func TestAdminNotificationPermissions(t *testing.T) {
	nonAdminTokens := []struct {
		name      string
		tokenPath string
	}{
		{"org owner", *orgOwnerE2EAccessTokenPath},
		{"environment editor", *envEditorAccessTokenPath},
		{"environment viewer", *envViewerAccessTokenPath},
	}

	for _, tc := range nonAdminTokens {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			requireAccessToken(t, tc.tokenPath, tc.name)
			c := newNotificationClient(t, tc.tokenPath)

			t.Run("ListDraftAdminNotifications", func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				_, err := c.ListDraftAdminNotifications(ctx, &proto.ListDraftAdminNotificationsRequest{})
				requirePermissionDenied(t, err, tc.name+" ListDraftAdminNotifications")
			})

			t.Run("CreateAdminNotification", func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				_, err := c.CreateAdminNotification(ctx, &proto.CreateAdminNotificationRequest{
					Localizations: []*proto.NotificationLocalization{
						{Language: "en", Title: "t", Content: "c"},
					},
				})
				requirePermissionDenied(t, err, tc.name+" CreateAdminNotification")
			})

			t.Run("UpdateAdminNotification", func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				_, err := c.UpdateAdminNotification(ctx, &proto.UpdateAdminNotificationRequest{
					Id: "non-existent-id",
					Localizations: []*proto.NotificationLocalization{
						{Language: "en", Title: "t", Content: "c"},
					},
				})
				requirePermissionDenied(t, err, tc.name+" UpdateAdminNotification")
			})

			t.Run("PublishAdminNotification", func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				_, err := c.PublishAdminNotification(ctx, &proto.PublishAdminNotificationRequest{Id: "non-existent-id"})
				requirePermissionDenied(t, err, tc.name+" PublishAdminNotification")
			})

			t.Run("DeleteAdminNotification", func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				_, err := c.DeleteAdminNotification(ctx, &proto.DeleteAdminNotificationRequest{Id: "non-existent-id"})
				requirePermissionDenied(t, err, tc.name+" DeleteAdminNotification")
			})
		})
	}
}

func containsNotificationID(notifications []*proto.Notification, id string) bool {
	for _, n := range notifications {
		if n.Id == id {
			return true
		}
	}
	return false
}

func requireStatusCode(t *testing.T, err error, code codes.Code, op string) {
	t.Helper()
	require.Error(t, err, "%s: expected an error", op)
	st, ok := status.FromError(err)
	require.True(t, ok, "%s: expected a gRPC status error, but got: %v", op, err)
	assert.Equal(t, code, st.Code(), "%s: unexpected status code, got: %v", op, err)
}

// requirePermissionDenied accepts PermissionDenied (authenticated but not a
// system admin) or Unauthenticated (no/invalid access token).
func requirePermissionDenied(t *testing.T, err error, op string) {
	t.Helper()
	require.Error(t, err, "%s: expected a permission error, but the call succeeded", op)
	st, ok := status.FromError(err)
	require.True(t, ok, "%s: expected a gRPC status error, but got: %v", op, err)
	switch st.Code() {
	case codes.PermissionDenied, codes.Unauthenticated:
		return
	default:
		t.Fatalf("%s: expected PermissionDenied or Unauthenticated, but got %s: %v", op, st.Code(), err)
	}
}

func requireAccessToken(t *testing.T, tokenPath, role string) {
	t.Helper()
	if tokenPath == "" {
		t.Skipf("skipping: no access token provided for the %s account", role)
	}
}

func newNotificationClient(t *testing.T, tokenPath string) proto.NotificationServiceClient {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(tokenPath)
	require.NoError(t, err)
	conn, err := rpcclient.NewClientConn(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	return proto.NewNotificationServiceClient(conn)
}

func uniqueSuffix() string {
	return fmt.Sprintf("%s-%d-%s", *testID, time.Now().UnixNano(), randomString())
}

func randomString() string {
	b := make([]byte, 10)
	_, _ = rand.Read(b)
	return strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b))
}
