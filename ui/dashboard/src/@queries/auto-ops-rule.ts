import {
  autoOpsRuleFetcher,
  AutoOpsRuleFetcherParams,
  AutoOpsRuleFetcherResponse
} from '@api/auto-ops';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<AutoOpsRuleFetcherResponse> & {
  params: AutoOpsRuleFetcherParams;
};

export const AUTO_OPS_RULE_KEY = 'auto-ops-rule';

export const useQueryAutoOpsRule = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [AUTO_OPS_RULE_KEY, params],
    queryFn: async () => {
      return autoOpsRuleFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAutoOpsRule = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [AUTO_OPS_RULE_KEY, params],
    queryFn: async () => {
      return autoOpsRuleFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAutoOpsRule = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [AUTO_OPS_RULE_KEY, params],
    queryFn: async () => {
      return autoOpsRuleFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateAutoOpsRule = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [AUTO_OPS_RULE_KEY]
  });
};
