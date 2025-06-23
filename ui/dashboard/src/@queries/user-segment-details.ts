import {
  userSegmentFetcher,
  UserSegmentFetcherParams,
  UserSegmentResponse
} from '@api/user-segment';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<UserSegmentResponse> & {
  params?: UserSegmentFetcherParams;
};

export const SEGMENT_DETAILS_QUERY_KEY = 'segment-details';

export const useQueryUserSegment = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [SEGMENT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return userSegmentFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchUserSegment = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [SEGMENT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return userSegmentFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchUserSegment = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [SEGMENT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return userSegmentFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateUserSegment = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [SEGMENT_DETAILS_QUERY_KEY]
  });
};
