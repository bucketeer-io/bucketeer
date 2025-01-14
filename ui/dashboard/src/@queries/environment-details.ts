import {
  environmentFetcher,
  EnvironmentFetcherParams,
  EnvironmentResponse
} from '@api/environment';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<EnvironmentResponse> & {
  params?: EnvironmentFetcherParams;
};

export const ENVIRONMENT_DETAILS_QUERY_KEY = 'environment-details';

export const useQueryEnvironmentDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [ENVIRONMENT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return environmentFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchEnvironmentDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [ENVIRONMENT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return environmentFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchEnvironmentDetails = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [ENVIRONMENT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return environmentFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateEnvironmentDetails = (
  queryClient: QueryClient,
  params: EnvironmentFetcherParams
) => {
  queryClient.invalidateQueries({
    queryKey: [ENVIRONMENT_DETAILS_QUERY_KEY, params]
  });
};
