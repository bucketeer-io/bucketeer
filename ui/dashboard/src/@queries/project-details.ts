import {
  projectDetailsFetcher,
  ProjectDetailsFetcherParams,
  ProjectDetailsResponse
} from '@api/project';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<ProjectDetailsResponse> & {
  params?: ProjectDetailsFetcherParams;
};

export const ORGANIZATION_DETAILS_QUERY_KEY = 'project-details';

export const useQueryProjectDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [ORGANIZATION_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return projectDetailsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchProjectDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [ORGANIZATION_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return projectDetailsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchProjectDetails = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [ORGANIZATION_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return projectDetailsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateProjectDetails = (
  queryClient: QueryClient,
  params: ProjectDetailsFetcherParams
) => {
  queryClient.invalidateQueries({
    queryKey: [ORGANIZATION_DETAILS_QUERY_KEY, params]
  });
};
