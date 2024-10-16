import {
  organizationDetailsFetcher,
  OrganizationDetailsFetcherParams,
  OrganizationDetailsResponse
} from '@api/organization';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<OrganizationDetailsResponse> & {
  params?: OrganizationDetailsFetcherParams;
};

export const ORGANIZATION_DETAILS_QUERY_KEY = 'organization-details';

export const useQueryOrganizationDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [ORGANIZATION_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return organizationDetailsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchOrganizationDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [ORGANIZATION_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return organizationDetailsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchOrganizationDetails = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [ORGANIZATION_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return organizationDetailsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateOrganizationDetails = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [ORGANIZATION_DETAILS_QUERY_KEY]
  });
};
