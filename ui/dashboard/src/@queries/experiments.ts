import { ExperimentsFetcherParams, experimentsFetcher } from '@api/experiment';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { ExperimentCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<ExperimentCollection> & {
  params?: ExperimentsFetcherParams;
};

export const EXPERIMENTS_QUERY_KEY = 'experiments';

export const useQueryExperiments = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [EXPERIMENTS_QUERY_KEY, params],
    queryFn: async () => {
      return experimentsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchExperiments = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [EXPERIMENTS_QUERY_KEY, params],
    queryFn: async () => {
      return experimentsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchExperiments = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [EXPERIMENTS_QUERY_KEY, params],
    queryFn: async () => {
      return experimentsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateExperiments = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [EXPERIMENTS_QUERY_KEY]
  });
};
