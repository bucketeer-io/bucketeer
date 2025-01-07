import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, ProjectCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';

export interface ProjectsFetcherParams extends CollectionParams {
  organizationId?: string;
}

export const projectsFetcher = async (
  params?: ProjectsFetcherParams
): Promise<ProjectCollection> => {
  return axiosClient
    .post<ProjectCollection>(
      '/v1/environment/list_projects_v2',
      pickBy(params, v => isNotEmpty(v))
    )
    .then(response => response.data);
};
