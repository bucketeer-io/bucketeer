import {
  evaluationTimeseriesCount,
  EvaluationTimeseriesCountParams
} from '@api/evaluation';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { EvaluationCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<EvaluationCollection> & {
  params?: EvaluationTimeseriesCountParams;
};

export const EVALUATION_QUERY_KEY = 'evaluation-timeseries';

export const useQueryEvaluation = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [EVALUATION_QUERY_KEY, params],
    queryFn: async () => {
      return evaluationTimeseriesCount(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchEvaluation = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [EVALUATION_QUERY_KEY, params],
    queryFn: async () => {
      return evaluationTimeseriesCount(params);
    },
    ...queryOptions
  });
};

export const prefetchEvaluation = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [EVALUATION_QUERY_KEY, params],
    queryFn: async () => {
      return evaluationTimeseriesCount(params);
    },
    ...queryOptions
  });
};

export const invalidateEvaluation = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [EVALUATION_QUERY_KEY]
  });
};
