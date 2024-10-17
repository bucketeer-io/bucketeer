import axiosClient from '@api/axios-client';
import { CollectionParams, ProjectCollection } from '@types';

export interface ProjectsFetcherParams extends CollectionParams {
  organizationIds?: string[];
  archived?: boolean;
}

export const projectsFetcher = async (
  params?: ProjectsFetcherParams
): Promise<ProjectCollection> => {
  return axiosClient
    .post<ProjectCollection>('/v1/environment/list_projects_v2', params)
    .then(response => response.data);
};
