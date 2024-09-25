import {
  DetailsFetcherParams,
  ProjectDetailsCollection,
  projectDetailsFetcher
} from '@api/project';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<ProjectDetailsCollection> & {
  params?: DetailsFetcherParams;
};

export const PROJECT_DETAILS_QUERY_KEY = 'project_details';

export const useQueryProjectDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [PROJECT_DETAILS_QUERY_KEY, params],
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
    queryKey: [PROJECT_DETAILS_QUERY_KEY, params],
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
    queryKey: [PROJECT_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return projectDetailsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateProjectDetails = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [PROJECT_DETAILS_QUERY_KEY]
  });
};
