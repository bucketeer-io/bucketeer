import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { Project } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface ProjectDetailsFetcherParams {
  id: string;
}
export interface ProjectDetailsResponse {
  project: Project;
}

export const projectDetailsFetcher = async (
  _params?: ProjectDetailsFetcherParams
): Promise<ProjectDetailsResponse> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<ProjectDetailsResponse>(
      `/v1/environment/get_project?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
