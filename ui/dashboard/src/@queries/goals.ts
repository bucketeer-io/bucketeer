import { goalsFetcher, GoalsFetcherParams } from '@api/goal';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { GoalCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<GoalCollection> & {
  params?: GoalsFetcherParams;
};

export const GOALS_QUERY_KEY = 'goals';

export const useQueryGoals = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [GOALS_QUERY_KEY, params],
    queryFn: async () => {
      return goalsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchGoals = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [GOALS_QUERY_KEY, params],
    queryFn: async () => {
      return goalsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchGoals = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [GOALS_QUERY_KEY, params],
    queryFn: async () => {
      return goalsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateGoals = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [GOALS_QUERY_KEY]
  });
};
