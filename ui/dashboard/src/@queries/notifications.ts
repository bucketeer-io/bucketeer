import {
  notificationsFetcher,
  NotificationsFetcherParams
} from '@api/notification';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { NotificationsCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<NotificationsCollection> & {
  params?: NotificationsFetcherParams;
};

export const NOTIFICATIONS_QUERY_KEY = 'notifications';

export const useQueryNotifications = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [NOTIFICATIONS_QUERY_KEY, params],
    queryFn: async () => {
      return notificationsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchNotifications = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [NOTIFICATIONS_QUERY_KEY, params],
    queryFn: async () => {
      return notificationsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchNotifications = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [NOTIFICATIONS_QUERY_KEY, params],
    queryFn: async () => {
      return notificationsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateNotifications = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [NOTIFICATIONS_QUERY_KEY]
  });
};
