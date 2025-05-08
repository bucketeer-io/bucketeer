import { historiesFetcher, HistoriesFetcherParams } from '@api/histories';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { AuditLogCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<AuditLogCollection> & {
  params?: HistoriesFetcherParams;
};

export const HISTORIES_QUERY_KEY = 'histories';

export const useQueryHistories = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [HISTORIES_QUERY_KEY, params],
    queryFn: async () => {
      return historiesFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchHistories = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [HISTORIES_QUERY_KEY, params],
    queryFn: async () => {
      return historiesFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchHistories = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [HISTORIES_QUERY_KEY, params],
    queryFn: async () => {
      return historiesFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateHistories = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [HISTORIES_QUERY_KEY]
  });
};
