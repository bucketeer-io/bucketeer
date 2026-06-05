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

/**
 * Optimistic in-place update of a single feature-details cache entry.
 *
 * The axios response interceptor invalidates the query for free on every
 * non-GET response, but invalidation triggers a refetch. Use this helper
 * when the mutation response already contains the new feature payload and
 * you want the UI to update without waiting for the round-trip.
 */
export const updateFeatureCache = (
  queryClient: QueryClient,
  params: FeatureFetcherParams,
  updatedFeature: FeatureResponse['feature']
) => {
  queryClient.setQueryData<FeatureResponse>([FEATURE_QUERY_KEY, params], old =>
    old ? { ...old, feature: updatedFeature } : old
  );
};
