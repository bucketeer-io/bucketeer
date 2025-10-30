import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { EnvironmentResponse } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface EnvironmentResultDetailsFetcherParams {
  id: string;
}

export const environmentResultDetailsFetcher = async (
  params?: EnvironmentResultDetailsFetcherParams
): Promise<EnvironmentResponse> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));
  return axiosClient
    .get<EnvironmentResponse>(
      `/v1/environment/get_environment?${requestParams}`
    )
    .then(response => response.data);
};
