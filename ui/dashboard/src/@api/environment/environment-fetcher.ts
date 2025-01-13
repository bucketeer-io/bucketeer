import axiosClient from '@api/axios-client';
import { pickBy } from 'lodash';
import { isNotEmpty } from 'utils/data-type';
import { EnvironmentResponse } from './environment-creator';

export interface EnvironmentFetcherParams {
  id: string;
}

export const environmentFetcher = async (
  params?: EnvironmentFetcherParams
): Promise<EnvironmentResponse> => {
  return axiosClient
    .get<EnvironmentResponse>('/v1/environment/get_environment', {
      params: pickBy(params, v => isNotEmpty(v))
    })
    .then(response => response.data);
};
