import {
  autoOpsFetcher,
  AutoOpsFetcherParams
} from '@api/auto-ops/auto-ops-fetcher';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { AutoOpsRuleCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<AutoOpsRuleCollection> & {
  params?: AutoOpsFetcherParams;
};

export const AUTO_OPS_KEY = 'auto-ops';

export const useQueryAutoOps = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [AUTO_OPS_KEY, params],
    queryFn: async () => {
      return autoOpsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAutoOps = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [AUTO_OPS_KEY, params],
    queryFn: async () => {
      return autoOpsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAutoOps = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [AUTO_OPS_KEY, params],
    queryFn: async () => {
      return autoOpsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateAutoOps = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [AUTO_OPS_KEY]
  });
};
