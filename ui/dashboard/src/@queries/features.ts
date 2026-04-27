import { featuresFetcher, FeaturesFetcherParams } from '@api/features';
import {
  QueryClient,
  useInfiniteQuery,
  useQuery,
  useQueryClient
} from '@tanstack/react-query';
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

type UseInfiniteQueryFeaturesParams = Pick<
  FeaturesFetcherParams,
  'environmentId' | 'archived'
> & {
  searchKeyword?: string;
  pageSize?: number;
};

export const useInfiniteQueryFeatures = ({
  environmentId,
  searchKeyword = '',
  pageSize,
  archived
}: UseInfiniteQueryFeaturesParams) =>
  useInfiniteQuery({
    queryKey: [
      FEATURES_QUERY_KEY,
      'infinite',
      { environmentId, searchKeyword, pageSize }
    ],
    queryFn: ({ pageParam = 0 }) =>
      featuresFetcher({
        cursor: String(pageParam),
        pageSize,
        searchKeyword,
        environmentId,
        archived
      }),
    getNextPageParam: (lastPage, allPages) => {
      const fetched = allPages.reduce((sum, p) => sum + p.features.length, 0);
      const total = Number(lastPage.totalCount);
      return fetched < total ? fetched : undefined;
    },
    initialPageParam: 0
  });
