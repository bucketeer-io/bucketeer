import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useAuth } from 'auth';
import { getLanguage } from 'i18n';
import { NotificationFilters, PublishNotificationInput } from '../types';
import {
  fetchDrafts,
  fetchFeed,
  markAllAsRead,
  markAsRead,
  markManyAsRead,
  publishNotification,
  saveDraft,
  updateNotification
} from './mock-service';

const DEFAULT_PAGE_SIZE = 10;

const useViewer = () => {
  const { consoleAccount } = useAuth();
  return { email: consoleAccount?.email ?? '', lang: getLanguage() };
};

const feedKey = (
  environmentId: string,
  email: string,
  lang: string,
  read: boolean,
  page: number,
  filters: NotificationFilters
) => [
  'notification-feed',
  environmentId,
  email,
  lang,
  read,
  page,
  filters.searchQuery,
  filters.sort,
  filters.from,
  filters.to
];

const draftsKey = (environmentId: string, email: string, lang: string) => [
  'notification-drafts',
  environmentId,
  email,
  lang
];

const invalidateFeed = (
  queryClient: ReturnType<typeof useQueryClient>,
  environmentId: string
) =>
  queryClient.invalidateQueries({
    queryKey: ['notification-feed', environmentId]
  });

export const useFetchFeed = (
  environmentId: string,
  read: boolean,
  page: number,
  filters: NotificationFilters
) => {
  const { email, lang } = useViewer();
  return useQuery({
    queryKey: feedKey(environmentId, email, lang, read, page, filters),
    queryFn: () =>
      fetchFeed(environmentId, email, lang, {
        read,
        page,
        pageSize: DEFAULT_PAGE_SIZE,
        searchQuery: filters.searchQuery,
        sort: filters.sort,
        from: filters.from,
        to: filters.to
      })
  });
};

export const useFetchDrafts = (environmentId: string) => {
  const { email, lang } = useViewer();
  return useQuery({
    queryKey: draftsKey(environmentId, email, lang),
    queryFn: () => fetchDrafts(environmentId, email, lang)
  });
};

export const useMarkAsRead = (environmentId: string) => {
  const { email } = useViewer();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => markAsRead(environmentId, id, email),
    onSuccess: () => invalidateFeed(queryClient, environmentId)
  });
};

export const useMarkManyAsRead = (environmentId: string) => {
  const { email } = useViewer();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (ids: string[]) => markManyAsRead(environmentId, ids, email),
    onSuccess: () => invalidateFeed(queryClient, environmentId)
  });
};

export const useMarkAllAsRead = (environmentId: string) => {
  const { email } = useViewer();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => markAllAsRead(environmentId, email),
    onSuccess: () => invalidateFeed(queryClient, environmentId)
  });
};

export const usePublishNotification = (environmentId: string) => {
  const { email } = useViewer();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: PublishNotificationInput) =>
      publishNotification(environmentId, email, input),
    onSuccess: () => invalidateFeed(queryClient, environmentId)
  });
};

export const useSaveDraft = (environmentId: string) => {
  const { email } = useViewer();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: PublishNotificationInput) =>
      saveDraft(environmentId, email, input),
    onSuccess: () =>
      queryClient.invalidateQueries({
        queryKey: ['notification-drafts', environmentId]
      })
  });
};

export const useUpdateNotification = (environmentId: string) => {
  const { email } = useViewer();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      input
    }: {
      id: string;
      input: PublishNotificationInput;
    }) => updateNotification(environmentId, email, id, input),
    onSuccess: () => {
      invalidateFeed(queryClient, environmentId);
      queryClient.invalidateQueries({
        queryKey: ['notification-drafts', environmentId]
      });
    }
  });
};
