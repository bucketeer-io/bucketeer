import { accountsFetcher, AccountsFetcherParams } from '@api/account';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
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

export const invalidateAccounts = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [ACCOUNTS_QUERY_KEY]
  });
};
