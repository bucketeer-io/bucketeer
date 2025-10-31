import {
  notificationFetcher,
  NotificationFetcherPayload,
  NotificationResponse
} from '@api/notification';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<NotificationResponse> & {
  params?: NotificationFetcherPayload;
};

export const NOTIFICATION_DETAILS_QUERY_KEY = 'notification-details';

export const useQueryNotification = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [NOTIFICATION_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return notificationFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchNotification = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [NOTIFICATION_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return notificationFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchNotification = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [NOTIFICATION_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return notificationFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateNotificationDetails = (
  queryClient: QueryClient,
  params: NotificationFetcherPayload
) => {
  queryClient.invalidateQueries({
    queryKey: [NOTIFICATION_DETAILS_QUERY_KEY, params]
  });
};

export const invalidateNotification = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [NOTIFICATION_DETAILS_QUERY_KEY]
  });
};
