import { accountsFetcher, AccountsFetcherParams } from '@api/account';
import {
  QueryClient,
  useInfiniteQuery,
  useQuery,
  useQueryClient
} from '@tanstack/react-query';
import type { AccountCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<AccountCollection> & {
  params?: AccountsFetcherParams;
};

export const ACCOUNTS_QUERY_KEY = 'accounts';

export const useQueryAccounts = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [ACCOUNTS_QUERY_KEY, params],
    queryFn: async () => {
      return accountsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAccounts = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [ACCOUNTS_QUERY_KEY, params],
    queryFn: async () => {
      return accountsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAccounts = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [ACCOUNTS_QUERY_KEY, params],
    queryFn: async () => {
      return accountsFetcher(params);
    },
    ...queryOptions
  });
};

type InfiniteQueryOptions = {
  params?: Omit<AccountsFetcherParams, 'cursor'>;
  enabled?: boolean;
};

export const useInfiniteQueryAccounts = ({
  params,
  enabled = true
}: InfiniteQueryOptions = {}) =>
  useInfiniteQuery({
    queryKey: [ACCOUNTS_QUERY_KEY, 'infinite', params],
    queryFn: ({ pageParam = '0' }) =>
      accountsFetcher({ ...params, cursor: pageParam as string }),
    initialPageParam: '0',
    getNextPageParam: (lastPage: AccountCollection) => {
      const nextCursor = Number(lastPage.cursor);
      const total = Number(lastPage.totalCount);
      return nextCursor < total ? String(nextCursor) : undefined;
    },
    enabled
  });
