import { autoOpsRulesFetcher, AutoOpsRulesFetcherParams } from '@api/auto-ops';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { AutoOpsRuleCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<AutoOpsRuleCollection> & {
  params: AutoOpsRulesFetcherParams;
};

export const AUTO_OPS_RULES_KEY = 'auto-ops-rules';

export const useQueryAutoOpsRules = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [AUTO_OPS_RULES_KEY, params],
    queryFn: async () => {
      return autoOpsRulesFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAutoOpsRules = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [AUTO_OPS_RULES_KEY, params],
    queryFn: async () => {
      return autoOpsRulesFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAutoOpsRules = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [AUTO_OPS_RULES_KEY, params],
    queryFn: async () => {
      return autoOpsRulesFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateAutoOpsRules = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [AUTO_OPS_RULES_KEY]
  });
};
