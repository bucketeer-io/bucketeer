import {
  notificationDraftsFetcher,
  NotificationDraftsFetcherParams,
  notificationFetcher,
  NotificationFetcherParams,
  notificationsFetcher,
  NotificationsFetcherParams,
  notificationUnreadCountFetcher,
  NotificationUnreadCountFetcherParams
} from '@api/notification-center';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'i18n';
import type {
  NotificationCenterDraftCollection,
  NotificationCenterFeedCollection,
  NotificationCenterFeedItem,
  NotificationCenterUnreadCount,
  QueryOptionsRespond
} from '@types';

export const NOTIFICATION_FEED_QUERY_KEY = 'notification-feed';
export const NOTIFICATION_QUERY_KEY = 'notification';
export const NOTIFICATION_DRAFTS_QUERY_KEY = 'notification-drafts';
export const NOTIFICATION_UNREAD_COUNT_QUERY_KEY = 'notification-unread-count';

type FeedQueryOptions =
  QueryOptionsRespond<NotificationCenterFeedCollection> & {
    params?: NotificationsFetcherParams;
  };

export const useQueryNotificationFeed = (options?: FeedQueryOptions) => {
  const { params, ...queryOptions } = options || {};
  // i18n.language is included in the query key so switching languages
  // refetches instead of serving the previously cached localization.
  const { i18n } = useTranslation('common');
  return useQuery({
    queryKey: [NOTIFICATION_FEED_QUERY_KEY, params, i18n.language],
    queryFn: () => notificationsFetcher(params),
    ...queryOptions
  });
};

type NotificationQueryOptions =
  QueryOptionsRespond<NotificationCenterFeedItem> & {
    params: NotificationFetcherParams;
  };

export const useQueryNotification = (options: NotificationQueryOptions) => {
  const { params, ...queryOptions } = options;
  const { i18n } = useTranslation('common');
  return useQuery({
    queryKey: [NOTIFICATION_QUERY_KEY, params, i18n.language],
    queryFn: () => notificationFetcher(params),
    ...queryOptions
  });
};

type DraftsQueryOptions =
  QueryOptionsRespond<NotificationCenterDraftCollection> & {
    params?: NotificationDraftsFetcherParams;
  };

export const useQueryNotificationDrafts = (options?: DraftsQueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const { i18n } = useTranslation('common');
  return useQuery({
    queryKey: [NOTIFICATION_DRAFTS_QUERY_KEY, params, i18n.language],
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
