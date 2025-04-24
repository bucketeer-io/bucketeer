import {
  rolloutFetcher,
  RolloutFetcherParams,
  RolloutFetcherResponse
} from '@api/rollouts';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<RolloutFetcherResponse> & {
  params?: RolloutFetcherParams;
};

export const ROLLOUT_KEY = 'rollout';

export const useQueryRollout = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [ROLLOUT_KEY, params],
    queryFn: async () => {
      return rolloutFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchRollout = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [ROLLOUT_KEY, params],
    queryFn: async () => {
      return rolloutFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchRollout = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [ROLLOUT_KEY, params],
    queryFn: async () => {
      return rolloutFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateRollout = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [ROLLOUT_KEY]
  });
};
