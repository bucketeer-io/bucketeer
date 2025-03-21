import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { ScheduleFlagCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface FeaturesScheduleFetcherParams {
  environmentId: string;
  featureId: string;
}

export const featuresScheduleFetcher = async (
  params?: FeaturesScheduleFetcherParams
): Promise<ScheduleFlagCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<ScheduleFlagCollection>(`/v1/schedule_flags?${requestParams}`)
    .then(response => response.data);
};
