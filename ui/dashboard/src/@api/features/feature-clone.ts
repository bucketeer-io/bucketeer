import axiosClient from '@api/axios-client';
import { Feature } from '@types';

export interface FeatureCloneParams {
  id: string;
  environmentId: string;
  targetEnvironmentId: string;
}

export interface FeatureResponse {
  feature: Feature;
}

export const featureClone = async (
  params?: FeatureCloneParams
): Promise<FeatureResponse> => {
  return axiosClient
    .post<FeatureResponse>(`/v1/feature/clone`, params)
    .then(response => response.data);
};
