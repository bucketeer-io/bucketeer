import {
  notificationCreator,
  notificationDelete,
  notificationMarkAllAsRead,
  notificationMarkAsRead,
  NotificationReadStatus,
  notificationPublisher,
  notificationUpdater
} from '@api/notification-center';
import {
  useQueryNotification,
  useQueryNotificationDrafts,
  useQueryNotificationFeed,
  useQueryNotificationUnreadCount
} from '@queries/notification-center';
import { useMutation } from '@tanstack/react-query';
import { LIST_PAGE_SIZE } from 'constants/app';
import { getLanguage } from 'i18n';
import {
  NotificationCenterPublishPayload,
  NotificationCenterStatus
} from '@types';
import { NotificationFilters } from '../types';

const DEFAULT_PAGE_SIZE = 10;

export const useFetchFeed = (
  readStatus: NotificationReadStatus,
  page: number,
  filters: NotificationFilters
) => {
  const cursor = (page - 1) * DEFAULT_PAGE_SIZE;
  return useQueryNotificationFeed({
    params: {
      readStatus,
      cursor: String(cursor),
      pageSize: DEFAULT_PAGE_SIZE,
      searchKeyword: filters.searchQuery,
      orderDirection: filters.sort === 'oldest' ? 'ASC' : 'DESC',
      publishedAtFrom: filters.from,
      publishedAtTo: filters.to,
      language: getLanguage()
    }
  });
};

export const useFetchNotification = (id?: string) => {
  return useQueryNotification({
    params: { id: id ?? '', language: getLanguage() },
    enabled: !!id
  });
};

export const useFetchDrafts = (enabled = true) => {
  return useQueryNotificationDrafts({
    params: {
      cursor: '0',
      pageSize: LIST_PAGE_SIZE
    },
    enabled
  });
};

export const useFetchUnreadCount = () => {
  return useQueryNotificationUnreadCount();
};

export const useFetchTabCounts = () => {
  const { data: unreadCount } = useFetchUnreadCount();
  const { data: readFeed } = useQueryNotificationFeed({
    params: { readStatus: 'READ', cursor: '0', pageSize: 1 }
  });

  return {
    unreadCount: unreadCount ?? 0,
    readCount: Number(readFeed?.totalCount ?? 0)
  };
};

export const useMarkAsRead = () => {
  return useMutation({
    mutationFn: (id: string) => notificationMarkAsRead({ ids: [id] })
  });
};

export const useMarkManyAsRead = () => {
  return useMutation({
    mutationFn: (ids: string[]) => notificationMarkAsRead({ ids })
  });
};

export const useMarkAllAsRead = () => {
  return useMutation({
    mutationFn: () => notificationMarkAllAsRead()
  });
};

export const usePublishNotification = () => {
  return useMutation({
    mutationFn: async (input: NotificationCenterPublishPayload) => {
      const created = await notificationCreator({
        localizations: input.localizations
      });
      return notificationPublisher({ id: created.notification.id });
    }
  });
};

export const usePublishDraft = () => {
  return useMutation({
    mutationFn: (id: string) => notificationPublisher({ id })
  });
};

export const useSaveDraft = () => {
  return useMutation({
    mutationFn: (input: NotificationCenterPublishPayload) =>
      notificationCreator({
        localizations: input.localizations
      })
  });
};

export const useUpdateNotification = () => {
  return useMutation({
    mutationFn: async ({
      id,
      input
    }: {
      id: string;
      input: NotificationCenterPublishPayload;
    }) => {
      if (input.status === NotificationCenterStatus.PUBLISHED) {
        await notificationUpdater({ id, localizations: input.localizations });
        return notificationPublisher({ id });
      }
      return notificationUpdater({
        id,
        localizations: input.localizations
      });
    }
  });
};

export const useDeleteNotification = () => {
  return useMutation({
    mutationFn: (id: string) => notificationDelete({ id })
  });
};
