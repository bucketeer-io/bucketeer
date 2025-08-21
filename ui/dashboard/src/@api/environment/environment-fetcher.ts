import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { EnvironmentResponse } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface EnvironmentFetcherParams {
  id: string;
}

export const environmentFetcher = async (
  _params?: EnvironmentFetcherParams
): Promise<EnvironmentResponse> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<EnvironmentResponse>(
      `/v1/environment/get_environment?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
