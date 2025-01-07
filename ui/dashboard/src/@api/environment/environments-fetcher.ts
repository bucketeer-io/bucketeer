import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, EnvironmentCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface EnvironmentsFetcherParams extends CollectionParams {
  projectId?: string;
  organizationId?: string;
  archived?: boolean;
}

export const environmentsFetcher = async (
  _params?: EnvironmentsFetcherParams
): Promise<EnvironmentCollection> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<EnvironmentCollection>(
      `/v1/environment/list_environments?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
