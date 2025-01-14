import {
  accountByOrgFetcher,
  AccountByOrgFetcherParams,
  AccountByOrgFetcherResponse
} from '@api/account';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<AccountByOrgFetcherResponse> & {
  params?: AccountByOrgFetcherParams;
};

export const ACCOUNT_ORGANIZATION_DETAILS_QUERY_KEY =
  'account-by-organization-details';

export const useQueryAccountOrganizationDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [ACCOUNT_ORGANIZATION_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return accountByOrgFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAccountOrganizationDetails = (
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [ACCOUNT_ORGANIZATION_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return accountByOrgFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAccountOrganizationDetails = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [ACCOUNT_ORGANIZATION_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return accountByOrgFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateAccountOrganizationDetails = (
  queryClient: QueryClient,
  params: AccountByOrgFetcherParams
) => {
  queryClient.invalidateQueries({
    queryKey: [ACCOUNT_ORGANIZATION_DETAILS_QUERY_KEY, params]
  });
};
