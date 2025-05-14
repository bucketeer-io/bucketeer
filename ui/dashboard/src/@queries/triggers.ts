import { triggersFetcher, TriggersFetcherParams } from '@api/trigger';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { TriggerCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<TriggerCollection> & {
  params?: TriggersFetcherParams;
};

export const TRIGGERS_KEY = 'triggers';

export const useQueryTriggers = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [TRIGGERS_KEY, params],
    queryFn: async () => {
      return triggersFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchTriggers = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [TRIGGERS_KEY, params],
    queryFn: async () => {
      return triggersFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchTriggers = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [TRIGGERS_KEY, params],
    queryFn: async () => {
      return triggersFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateTriggers = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [TRIGGERS_KEY]
  });
};
