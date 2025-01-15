import {
  userSegmentsFetcher,
  UserSegmentsFetcherParams
} from '@api/user-segment';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { UserSegmentCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<UserSegmentCollection> & {
  params?: UserSegmentsFetcherParams;
};

export const USER_SEGMENTS_QUERY_KEY = 'user-segments';

export const useQueryUserSegments = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [USER_SEGMENTS_QUERY_KEY, params],
    queryFn: async () => {
      return userSegmentsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchUserSegments = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [USER_SEGMENTS_QUERY_KEY, params],
    queryFn: async () => {
      return userSegmentsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchUserSegments = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [USER_SEGMENTS_QUERY_KEY, params],
    queryFn: async () => {
      return userSegmentsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateUserSegments = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [USER_SEGMENTS_QUERY_KEY]
  });
};
