import { featuresFetcher, FeaturesFetcherParams } from '@api/features';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { FeatureCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<FeatureCollection> & {
  params?: FeaturesFetcherParams;
};

export const FEATURES_QUERY_KEY = 'features';

export const useQueryFeatures = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [FEATURES_QUERY_KEY, params],
    queryFn: async () => {
      return featuresFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchFeatures = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [FEATURES_QUERY_KEY, params],
    queryFn: async () => {
      return featuresFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchFeatures = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [FEATURES_QUERY_KEY, params],
    queryFn: async () => {
      return featuresFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateFeatures = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [FEATURES_QUERY_KEY]
  });
};
