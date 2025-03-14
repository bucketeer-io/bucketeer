import axiosClient from '@api/axios-client';
import { Feature } from '@types';

export interface FeatureFetcherParams {
  id: string;
  environmentId: string;
}

export interface FeatureResponse {
  feature: Feature;
}

export const featureClone = async (
  params?: FeatureFetcherParams
): Promise<FeatureResponse> => {
  return axiosClient
    .post<FeatureResponse>(`/v1/feature/clone`, params)
    .then(response => response.data);
};
