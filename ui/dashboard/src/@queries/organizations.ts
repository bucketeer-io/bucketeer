import {
  organizationsFetcher,
  OrganizationsFetcherParams
} from '@api/organization';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import pickBy from 'lodash/pickby';
import type { OrganizationCollection, QueryOptionsRespond } from '@types';
import { isNotEmpty } from 'utils/data-type';

type QueryOptions = QueryOptionsRespond<OrganizationCollection> & {
  params?: OrganizationsFetcherParams;
};

export const ORGANIZATIONS_QUERY_KEY = 'organizations';

export const useQueryOrganizations = (options?: QueryOptions) => {
  const { params: _params, ...queryOptions } = options || {};
  const params = pickBy(_params, v => isNotEmpty(v));
  const query = useQuery({
    queryKey: [ORGANIZATIONS_QUERY_KEY, params],
    queryFn: async () => {
      return organizationsFetcher(_params);
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
