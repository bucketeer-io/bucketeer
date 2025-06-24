import { pushFetcher, PushFetcherParams, PushResponse } from '@api/push';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<PushResponse> & {
  params?: PushFetcherParams;
};

export const PUSH_DETAILS_QUERY_KEY = 'push-details';

export const useQueryPush = (options?: QueryOptions) => {
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

export const usePrefetchPush = (options?: QueryOptions) => {
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

export const prefetchPush = (
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

export const invalidatePush = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [PUSH_DETAILS_QUERY_KEY]
  });
};
