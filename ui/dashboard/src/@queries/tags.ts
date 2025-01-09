import { tagsFetcher, TagsFetcherParams } from '@api/tag';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { TagCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<TagCollection> & {
  params?: TagsFetcherParams;
};

export const TAGS_QUERY_KEY = 'tags';

export const useQueryTags = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [TAGS_QUERY_KEY, params],
    queryFn: async () => {
      return tagsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchTags = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [TAGS_QUERY_KEY, params],
    queryFn: async () => {
      return tagsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchTags = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [TAGS_QUERY_KEY, params],
    queryFn: async () => {
      return tagsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateTags = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [TAGS_QUERY_KEY]
  });
};
