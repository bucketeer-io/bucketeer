import { pushesFetcher, PushesFetcherParams } from '@api/push';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { PushCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<PushCollection> & {
  params?: PushesFetcherParams;
};

export const PUSHES_QUERY_KEY = 'pushes';

export const useQueryPushes = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [PUSHES_QUERY_KEY, params],
    queryFn: async () => {
      return pushesFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchPushes = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [PUSHES_QUERY_KEY, params],
    queryFn: async () => {
      return pushesFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchPushes = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [PUSHES_QUERY_KEY, params],
    queryFn: async () => {
      return pushesFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidatePushes = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [PUSHES_QUERY_KEY]
  });
};
