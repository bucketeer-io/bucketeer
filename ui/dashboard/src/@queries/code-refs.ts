import { codeRefsFetcher, CodeRefsFetcherParams } from '@api/code-refs';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { CodeRefCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<CodeRefCollection> & {
  params?: CodeRefsFetcherParams;
};

export const CODE_REFS_KEY = 'code-refs';

export const useQueryCodeRefs = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [CODE_REFS_KEY, params],
    queryFn: async () => {
      return codeRefsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchCodeRefs = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [CODE_REFS_KEY, params],
    queryFn: async () => {
      return codeRefsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchCodeRefs = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [CODE_REFS_KEY, params],
    queryFn: async () => {
      return codeRefsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateCodeRefs = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [CODE_REFS_KEY]
  });
};
