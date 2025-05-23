import { autoOpsCountFetcher, AutoOpsCountFetcherParams } from '@api/auto-ops';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { AutoOpsCountCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<AutoOpsCountCollection> & {
  params: AutoOpsCountFetcherParams;
};

export const AUTO_OPS_COUNT_KEY = 'auto-ops-count';

export const useQueryAutoOpsCount = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [AUTO_OPS_COUNT_KEY, params],
    queryFn: async () => {
      return autoOpsCountFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAutoOpsCount = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [AUTO_OPS_COUNT_KEY, params],
    queryFn: async () => {
      return autoOpsCountFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAutoOpsCount = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [AUTO_OPS_COUNT_KEY, params],
    queryFn: async () => {
      return autoOpsCountFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateAutoOpsCount = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [AUTO_OPS_COUNT_KEY]
  });
};
