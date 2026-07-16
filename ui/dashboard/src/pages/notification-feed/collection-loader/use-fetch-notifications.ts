import {
  notificationCreator,
  notificationDelete,
  notificationMarkAllAsRead,
  notificationMarkAsRead,
  notificationPublisher,
  notificationUpdater
} from '@api/notification-center';
import {
  useQueryNotificationDrafts,
  useQueryNotificationFeed,
  useQueryNotificationUnreadCount
} from '@queries/notification-center';
import { useMutation } from '@tanstack/react-query';
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

export const useMarkAsRead = (environmentId: string) => {
  return useMutation({
    mutationFn: (id: string) =>
      notificationMarkAsRead({ environmentId, ids: [id] })
  });
};

export const useMarkManyAsRead = (environmentId: string) => {
  return useMutation({
    mutationFn: (ids: string[]) =>
      notificationMarkAsRead({ environmentId, ids })
  });
};

export const useMarkAllAsRead = (environmentId: string) => {
  return useMutation({
    mutationFn: () => notificationMarkAllAsRead({ environmentId })
  });
};

export const usePublishNotification = (environmentId: string) => {
  return useMutation({
    mutationFn: (input: NotificationCenterPublishPayload) =>
      notificationPublisher({
        environmentId,
        localizations: input.localizations
      })
  });
};

export const useSaveDraft = (environmentId: string) => {
  return useMutation({
    mutationFn: (input: NotificationCenterPublishPayload) =>
      notificationCreator({
        environmentId,
        status: NotificationCenterStatus.DRAFT,
        localizations: input.localizations
      })
  });
};

export const useUpdateNotification = (environmentId: string) => {
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
    }
  });
};

export const useDeleteNotification = (environmentId: string) => {
  return useMutation({
    mutationFn: (id: string) => notificationDelete({ id, environmentId })
  });
};
