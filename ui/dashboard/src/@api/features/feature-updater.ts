import axiosClient from '@api/axios-client';
import { FeatureUpdaterParams } from '@types';
import { FeatureResponse } from './feature-clone';

export const featureUpdater = async (
  params?: Partial<FeatureUpdaterParams>
): Promise<FeatureResponse> => {
  return axiosClient
    .patch<FeatureResponse>('/v1/feature', params)
    .then(response => response.data);
};
