import { pushFetcher, PushFetcherParams, PushResponse } from '@api/push';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<PushResponse> & {
  params?: PushFetcherParams;
};

export const PUSH_DETAILS_QUERY_KEY = 'push-details';

export const useQueryPushDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [PUSH_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return pushFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchPushDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [PUSH_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return pushFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchPushDetails = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [PUSH_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return pushFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidatePushDetails = (
  queryClient: QueryClient,
  params: PushFetcherParams
) => {
  queryClient.invalidateQueries({
    queryKey: [PUSH_DETAILS_QUERY_KEY, params]
  });
};
