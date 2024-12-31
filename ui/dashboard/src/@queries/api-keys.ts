import { apiKeysFetcher, APIKeysFetcherParams } from '@api/api-key';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { APIKeyCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<APIKeyCollection> & {
  params?: APIKeysFetcherParams;
};

export const API_KEYS_QUERY_KEY = 'api-keys';

export const useQueryAPIKeys = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [API_KEYS_QUERY_KEY, params],
    queryFn: async () => {
      return apiKeysFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAPIKeys = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [API_KEYS_QUERY_KEY, params],
    queryFn: async () => {
      return apiKeysFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAPIKeys = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [API_KEYS_QUERY_KEY, params],
    queryFn: async () => {
      return apiKeysFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateAPIKeys = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [API_KEYS_QUERY_KEY]
  });
};
