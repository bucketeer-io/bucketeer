import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';
import { FeatureResponse } from './feature-clone';

export interface FeatureCloneParams {
  id: string;
  environmentId: string;
}

export const featureFetcher = async (
  params?: FeatureCloneParams
): Promise<FeatureResponse> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<FeatureResponse>(`/v1/feature?${requestParams}`)
    .then(response => response.data);
};
