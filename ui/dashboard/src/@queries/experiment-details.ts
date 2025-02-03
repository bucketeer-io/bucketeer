import {
  experimentFetcher,
  ExperimentFetcherParams,
  ExperimentResponse
} from '@api/experiment';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<ExperimentResponse> & {
  params?: ExperimentFetcherParams;
};

export const EXPERIMENT_DETAILS_QUERY_KEY = 'experiment-details';

export const useQueryExperimentDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [EXPERIMENT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return experimentFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchExperimentDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [EXPERIMENT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return experimentFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchExperimentDetails = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [EXPERIMENT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return experimentFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateExperimentDetails = (
  queryClient: QueryClient,
  params: ExperimentFetcherParams
) => {
  queryClient.invalidateQueries({
    queryKey: [EXPERIMENT_DETAILS_QUERY_KEY, params]
  });
};
