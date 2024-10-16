import axiosClient from '@api/axios-client';
import { CollectionParams, ProjectCollection } from '@types';

export interface ProjectFetcherParams extends CollectionParams {
  organizationIds: string[];
}

export const projectsFetcher = async (
  params?: ProjectFetcherParams
): Promise<ProjectCollection> => {
  return axiosClient
    .post<ProjectCollection>('/v1/environment/list_projects_v2', params)
    .then(response => response.data);
};
