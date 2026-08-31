import {
  userSegmentsFetcher,
  UserSegmentsFetcherParams
} from '@api/user-segment';
import {
  QueryClient,
  useInfiniteQuery,
  useQuery,
  useQueryClient
} from '@tanstack/react-query';
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

type InfiniteQueryOptions = {
  params?: Omit<UserSegmentsFetcherParams, 'cursor'>;
  enabled?: boolean;
};

export const useInfiniteQueryUserSegments = ({
  params,
  enabled = true
}: InfiniteQueryOptions = {}) =>
  useInfiniteQuery({
    queryKey: [USER_SEGMENTS_QUERY_KEY, 'infinite', params],
    queryFn: ({ pageParam = '0' }) =>
      userSegmentsFetcher({ ...params, cursor: pageParam as string }),
    initialPageParam: '0',
    getNextPageParam: (lastPage: UserSegmentCollection) => {
      const nextCursor = Number(lastPage.cursor);
      const total = Number(lastPage.totalCount);
      return nextCursor < total ? String(nextCursor) : undefined;
    },
    enabled
  });
