import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, ProjectCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface ProjectsFetcherParams extends CollectionParams {
  organizationId?: string;
}

export const projectsFetcher = async (
  _params?: ProjectsFetcherParams
): Promise<ProjectCollection> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<ProjectCollection>(
      `/v1/environment/list_projects_v2?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
