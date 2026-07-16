import {
  notificationDraftsFetcher,
  NotificationDraftsFetcherParams,
  notificationsFetcher,
  NotificationsFetcherParams,
  notificationUnreadCountFetcher,
  NotificationUnreadCountFetcherParams
} from '@api/notification-center';
import { useQuery } from '@tanstack/react-query';
import type {
  NotificationCenterDraftCollection,
  NotificationCenterFeedCollection,
  NotificationCenterUnreadCount,
  QueryOptionsRespond
} from '@types';

export const NOTIFICATION_FEED_QUERY_KEY = 'notification-feed';
export const NOTIFICATION_DRAFTS_QUERY_KEY = 'notification-drafts';
export const NOTIFICATION_UNREAD_COUNT_QUERY_KEY = 'notification-unread-count';

type FeedQueryOptions =
  QueryOptionsRespond<NotificationCenterFeedCollection> & {
    params?: NotificationsFetcherParams;
  };

export const useQueryNotificationFeed = (options?: FeedQueryOptions) => {
  const { params, ...queryOptions } = options || {};
  return useQuery({
    queryKey: [NOTIFICATION_FEED_QUERY_KEY, params],
    queryFn: () => notificationsFetcher(params),
    ...queryOptions
  });
};

type DraftsQueryOptions =
  QueryOptionsRespond<NotificationCenterDraftCollection> & {
    params?: NotificationDraftsFetcherParams;
  };

export const useQueryNotificationDrafts = (options?: DraftsQueryOptions) => {
  const { params, ...queryOptions } = options || {};
  return useQuery({
    queryKey: [NOTIFICATION_DRAFTS_QUERY_KEY, params],
    queryFn: () => notificationDraftsFetcher(params),
    ...queryOptions
  });
};

type UnreadCountQueryOptions =
  QueryOptionsRespond<NotificationCenterUnreadCount> & {
    params?: NotificationUnreadCountFetcherParams;
  };

export const useQueryNotificationUnreadCount = (
  options?: UnreadCountQueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  return useQuery({
    queryKey: [NOTIFICATION_UNREAD_COUNT_QUERY_KEY, params],
    queryFn: () => notificationUnreadCountFetcher(params),
    ...queryOptions
  });
};
