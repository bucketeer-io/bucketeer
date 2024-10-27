import { projectsFetcher, ProjectsFetcherParams } from '@api/project';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { ProjectCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<ProjectCollection> & {
  params?: ProjectsFetcherParams;
};

export const PROJECTS_QUERY_KEY = 'projects';

export const useQueryProjects = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [PROJECTS_QUERY_KEY, params],
    queryFn: async () => {
      return projectsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchProjects = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [PROJECTS_QUERY_KEY, params],
    queryFn: async () => {
      return projectsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchProjects = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [PROJECTS_QUERY_KEY, params],
    queryFn: async () => {
      return projectsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateProjects = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [PROJECTS_QUERY_KEY]
  });
};
