import axiosClient from '@api/axios-client';
import { CollectionParams, EnvironmentCollection } from '@types';

export interface EnvironmentsFetcherParams extends CollectionParams {
  projectId?: string;
  organizationId?: string;
  archived?: boolean;
}

export const environmentsFetcher = async (
  params?: EnvironmentsFetcherParams
): Promise<EnvironmentCollection> => {
  return axiosClient
    .post<EnvironmentCollection>('/v1/environment/list_environments', params)
    .then(response => response.data);
};
