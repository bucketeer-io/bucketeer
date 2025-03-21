import {
  FeatureResponse,
  featureFetcher,
  FeatureFetcherParams
} from '@api/features';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<FeatureResponse> & {
  params?: FeatureFetcherParams;
};

export const FEATURE_QUERY_KEY = 'feature-details';

export const useQueryFeature = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [FEATURE_QUERY_KEY, params],
    queryFn: async () => {
      return featureFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchFeature = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [FEATURE_QUERY_KEY, params],
    queryFn: async () => {
      return featureFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchFeature = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [FEATURE_QUERY_KEY, params],
    queryFn: async () => {
      return featureFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateFeature = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [FEATURE_QUERY_KEY]
  });
};
