import {
  notificationDraftsFetcher,
  NotificationDraftsFetcherParams,
  notificationFetcher,
  NotificationFetcherParams,
  notificationsFetcher,
  NotificationsFetcherParams
} from '@api/notification-center';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'i18n';
import type {
  NotificationCenterDraftCollection,
  NotificationCenterFeedCollection,
  NotificationCenterFeedItem,
  QueryOptionsRespond
} from '@types';

export const NOTIFICATION_FEED_QUERY_KEY = 'notification-feed';
export const NOTIFICATION_QUERY_KEY = 'notification';
export const NOTIFICATION_DRAFTS_QUERY_KEY = 'notification-drafts';

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
