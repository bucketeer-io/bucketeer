import {
  apiKeyFetcher,
  APIKeyFetcherParams,
  APIKeyResponse
} from '@api/api-key';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<APIKeyResponse> & {
  params?: APIKeyFetcherParams;
};

export const API_KEY_DETAILS_QUERY_KEY = 'api-keys-details';

export const useQueryAPIKey = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [API_KEY_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return apiKeyFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAPIKey = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [API_KEY_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return apiKeyFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAPIKey = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [API_KEY_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return apiKeyFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateAPIKey = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [API_KEY_DETAILS_QUERY_KEY]
  });
};
