import {
  experimentResultDetailsFetcher,
  ExperimentResultDetailsFetcherParams
} from '@api/experiment-result';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { ExperimentResultResponse, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<ExperimentResultResponse> & {
  params?: ExperimentResultDetailsFetcherParams;
};

export const EXPERIMENT_RESULT_DETAILS_QUERY_KEY = 'experiment-result-details';

export const useQueryExperimentResultDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [EXPERIMENT_RESULT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return experimentResultDetailsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchExperimentResultDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [EXPERIMENT_RESULT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return experimentResultDetailsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchExperimentResultDetails = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [EXPERIMENT_RESULT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return experimentResultDetailsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateExperimentResultDetails = (
  queryClient: QueryClient,
  params: ExperimentResultDetailsFetcherParams
) => {
  queryClient.invalidateQueries({
    queryKey: [EXPERIMENT_RESULT_DETAILS_QUERY_KEY, params]
  });
};
