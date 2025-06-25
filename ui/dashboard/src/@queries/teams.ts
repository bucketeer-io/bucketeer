import { teamsFetcher, TeamsFetcherParams } from '@api/team';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond, TeamCollection } from '@types';

type QueryOptions = QueryOptionsRespond<TeamCollection> & {
  params?: TeamsFetcherParams;
};

export const TEAMS_QUERY_KEY = 'teams';

export const useQueryTeams = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [TEAMS_QUERY_KEY, params],
    queryFn: async () => {
      return teamsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchTeams = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [TEAMS_QUERY_KEY, params],
    queryFn: async () => {
      return teamsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchTeams = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [TEAMS_QUERY_KEY, params],
    queryFn: async () => {
      return teamsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateTeams = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [TEAMS_QUERY_KEY]
  });
};
