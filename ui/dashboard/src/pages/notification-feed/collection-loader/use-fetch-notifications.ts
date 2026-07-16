import {
  notificationCreator,
  notificationDelete,
  notificationMarkAllAsRead,
  notificationMarkAsRead,
  notificationPublisher,
  notificationUpdater
} from '@api/notification-center';
import {
  NOTIFICATION_DRAFTS_QUERY_KEY,
  NOTIFICATION_FEED_QUERY_KEY,
  NOTIFICATION_UNREAD_COUNT_QUERY_KEY,
  useQueryNotificationDrafts,
  useQueryNotificationFeed,
  useQueryNotificationUnreadCount
} from '@queries/notification-center';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { LIST_PAGE_SIZE } from 'constants/app';
import {
  NotificationCenterPublishPayload,
  NotificationCenterStatus
} from '@types';
import { NotificationFilters } from '../types';

const DEFAULT_PAGE_SIZE = 10;

export const useFetchFeed = (
  environmentId: string,
  read: boolean,
  page: number,
  filters: NotificationFilters
) => {
  const cursor = (page - 1) * DEFAULT_PAGE_SIZE;
  return useQueryNotificationFeed({
    params: {
      environmentId,
      read,
      cursor: String(cursor),
      pageSize: DEFAULT_PAGE_SIZE,
      searchKeyword: filters.searchQuery,
      orderDirection: filters.sort === 'oldest' ? 'ASC' : 'DESC',
      from: filters.from ? String(Math.floor(filters.from / 1000)) : undefined,
      to: filters.to ? String(Math.floor(filters.to / 1000)) : undefined
    }
  });
};

export const useFetchDrafts = (environmentId: string, enabled = true) => {
  return useQueryNotificationDrafts({
    params: {
      environmentId,
      cursor: '0',
      pageSize: LIST_PAGE_SIZE
    },
    enabled
  });
};

export const useFetchUnreadCount = (environmentId: string) => {
  return useQueryNotificationUnreadCount({ params: { environmentId } });
};

export const useFetchTabCounts = (environmentId: string) => {
  const { data: unread } = useFetchUnreadCount(environmentId);
  const { data: readFeed } = useQueryNotificationFeed({
    params: { environmentId, read: true, cursor: '0', pageSize: 1 }
  });

  return {
    unreadCount: Number(unread?.count ?? 0),
    readCount: Number(readFeed?.readCount ?? 0)
  };
};

const invalidateFeed = (queryClient: ReturnType<typeof useQueryClient>) => {
  queryClient.invalidateQueries({ queryKey: [NOTIFICATION_FEED_QUERY_KEY] });
  queryClient.invalidateQueries({
    queryKey: [NOTIFICATION_UNREAD_COUNT_QUERY_KEY]
  });
};

export const useMarkAsRead = (environmentId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) =>
      notificationMarkAsRead({ environmentId, ids: [id] }),
    onSuccess: () => invalidateFeed(queryClient)
  });
};

export const useMarkManyAsRead = (environmentId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (ids: string[]) =>
      notificationMarkAsRead({ environmentId, ids }),
    onSuccess: () => invalidateFeed(queryClient)
  });
};

export const useMarkAllAsRead = (environmentId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => notificationMarkAllAsRead({ environmentId }),
    onSuccess: () => invalidateFeed(queryClient)
  });
};

export const usePublishNotification = (environmentId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: NotificationCenterPublishPayload) =>
      notificationPublisher({
        environmentId,
        localizations: input.localizations
      }),
    onSuccess: () => invalidateFeed(queryClient)
  });
};

export const useSaveDraft = (environmentId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: NotificationCenterPublishPayload) =>
      notificationCreator({
        environmentId,
        status: NotificationCenterStatus.DRAFT,
        localizations: input.localizations
      }),
    onSuccess: () =>
      queryClient.invalidateQueries({
        queryKey: [NOTIFICATION_DRAFTS_QUERY_KEY]
      })
  });
};

export const useUpdateNotification = (environmentId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      input
    }: {
      id: string;
      input: NotificationCenterPublishPayload;
    }) => {
      // Publishing an edited draft promotes it in place via the publish
      // endpoint (with its id); saving it as a draft again just updates it.
      if (input.status === NotificationCenterStatus.PUBLISHED) {
        return notificationPublisher({
          environmentId,
          id,
          localizations: input.localizations
        });
      }
      return notificationUpdater({
        id,
        environmentId,
        status: input.status,
        localizations: input.localizations
      });
    },
    onSuccess: () => {
      invalidateFeed(queryClient);
      queryClient.invalidateQueries({
        queryKey: [NOTIFICATION_DRAFTS_QUERY_KEY]
      });
    }
  });
};

export const useDeleteNotification = (environmentId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => notificationDelete({ id, environmentId }),
    onSuccess: () => {
      invalidateFeed(queryClient);
      queryClient.invalidateQueries({
        queryKey: [NOTIFICATION_DRAFTS_QUERY_KEY]
      });
    }
  });
};
