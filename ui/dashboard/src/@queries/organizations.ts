import {
  organizationsFetcher,
  OrganizationsFetcherParams
} from '@api/organization';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { OrganizationsCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<OrganizationsCollection> & {
  params?: OrganizationsFetcherParams;
};

export const ORGANIZATIONS_QUERY_KEY = 'organizations';

export const useQueryOrganizations = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [ORGANIZATIONS_QUERY_KEY, params],
    queryFn: async () => {
      return organizationsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchOrganizations = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [ORGANIZATIONS_QUERY_KEY, params],
    queryFn: async () => {
      return organizationsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchOrganizations = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [ORGANIZATIONS_QUERY_KEY, params],
    queryFn: async () => {
      return organizationsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateOrganizations = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [ORGANIZATIONS_QUERY_KEY]
  });
};
