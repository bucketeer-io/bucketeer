import {
  autoOpsGoalUserCountFetcher,
  AutoOpsGoalUserCountFetcherParams,
  AutoOpsGoalUserCountResponse
} from '@api/auto-ops/goal-user-count-fetcher';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<AutoOpsGoalUserCountResponse> & {
  params?: AutoOpsGoalUserCountFetcherParams;
};

export const OPS_GOAL_USER_COUNT_KEY = 'ops-goal-user-count';

export const useQueryAutoOpsGoalUserCount = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [OPS_GOAL_USER_COUNT_KEY, params],
    queryFn: async () => {
      return autoOpsGoalUserCountFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAutoOpsGoalUserCount = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [OPS_GOAL_USER_COUNT_KEY, params],
    queryFn: async () => {
      return autoOpsGoalUserCountFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAutoOpsGoalUserCount = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [OPS_GOAL_USER_COUNT_KEY, params],
    queryFn: async () => {
      return autoOpsGoalUserCountFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateAutoOpsGoalUserCount = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [OPS_GOAL_USER_COUNT_KEY]
  });
};
