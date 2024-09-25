import {
  EnvironmentsFetcherParams,
  environmentsFetcher
} from '@api/environment';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { EnvironmentCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<EnvironmentCollection> & {
  params?: EnvironmentsFetcherParams;
};

export const ENVIRONMENTS_QUERY_KEY = 'environments';

export const useQueryEnvironments = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [ENVIRONMENTS_QUERY_KEY, params],
    queryFn: async () => {
      return environmentsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchEnvironments = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [ENVIRONMENTS_QUERY_KEY, params],
    queryFn: async () => {
      return environmentsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchEnvironments = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [ENVIRONMENTS_QUERY_KEY, params],
    queryFn: async () => {
      return environmentsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateOrganizations = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [ENVIRONMENTS_QUERY_KEY]
  });
};
