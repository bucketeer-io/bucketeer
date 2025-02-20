import {
  goalDetailsFetcher,
  GoalDetailsFetcherParams,
  GoalDetailsResponse
} from '@api/goal';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<GoalDetailsResponse> & {
  params?: GoalDetailsFetcherParams;
};

export const GOAL_DETAILS_QUERY_KEY = 'goal-details';

export const useQueryGoalDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [GOAL_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return goalDetailsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchGoalDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [GOAL_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return goalDetailsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchGoalDetails = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [GOAL_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return goalDetailsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateGoalDetails = (
  queryClient: QueryClient,
  params: GoalDetailsFetcherParams
) => {
  queryClient.invalidateQueries({
    queryKey: [GOAL_DETAILS_QUERY_KEY, params]
  });
};
