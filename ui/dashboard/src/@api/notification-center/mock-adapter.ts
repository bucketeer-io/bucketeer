import type { AxiosInstance } from 'axios';
import MockAdapter from 'axios-mock-adapter';

// Most of the notification center backend has shipped (list, get, drafts,
// create, update, publish, delete — see proto/notification/service.proto).
// GetNotificationUnreadCount, MarkNotificationsAsRead, and
// MarkAllNotificationsAsRead are still stubs there (empty request/response
// messages; the handlers return Unimplemented), so this only intercepts
// those three to keep the bell badge and read/unread actions usable ahead of
// the backend. Delete a handler here (and eventually this whole file, plus
// its `use()` call in axios-client.ts) as each one ships for real.

// A rough, session-local stand-in for the bell badge count, since the real
// unread count isn't available without GetNotificationUnreadCount. Marking
// notifications read decrements it; there's no way to increment it back
// (e.g. a new notification arriving), since this mock has no visibility into
// the real list.
let unreadCount = 3;

export const installNotificationCenterMockAdapter = (client: AxiosInstance) => {
  const mock = new MockAdapter(client, { delayResponse: 300 });

  mock.onGet(/\/v1\/notifications\/unread_count/).reply(() => {
    return [200, { count: String(unreadCount) }];
  });

  mock.onPost('/v1/notifications/mark_as_read').reply(config => {
    const body = JSON.parse(config.data) as { ids: string[] };
    unreadCount = Math.max(0, unreadCount - body.ids.length);
    return [200, {}];
  });

  mock.onPost('/v1/notifications/mark_all_as_read').reply(() => {
    unreadCount = 0;
    return [200, {}];
  });

  // Anything else (list, get, drafts, create, update, publish, delete) goes
  // through to the real network — those endpoints are implemented now.
  mock.onAny().passThrough();
};
