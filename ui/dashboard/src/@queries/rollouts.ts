import {
  progressiveRolloutsFetcher,
  RolloutsFetcherParams
} from '@api/rollouts';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { RolloutCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<RolloutCollection> & {
  params?: RolloutsFetcherParams;
};

export const ROLLOUTS_KEY = 'rollouts';

export const useQueryRollouts = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [ROLLOUTS_KEY, params],
    queryFn: async () => {
      return progressiveRolloutsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchRollouts = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [ROLLOUTS_KEY, params],
    queryFn: async () => {
      return progressiveRolloutsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchRollouts = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [ROLLOUTS_KEY, params],
    queryFn: async () => {
      return progressiveRolloutsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateRollouts = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [ROLLOUTS_KEY]
  });
};
