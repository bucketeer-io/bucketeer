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

export interface FeatureBulkCloneParams {
  id: string;
  environmentId: string;
  targetEnvironmentIds: string[];
}

export type BulkCloneFeatureStatus =
  | 'BULK_CLONE_FEATURE_STATUS_UNSPECIFIED'
  | 'BULK_CLONE_FEATURE_STATUS_SUCCESS'
  | 'BULK_CLONE_FEATURE_STATUS_ALREADY_EXISTS'
  | 'BULK_CLONE_FEATURE_STATUS_FAILED';

export interface BulkCloneFeatureResult {
  environmentId: string;
  status: BulkCloneFeatureStatus;
  error: string;
}

export interface FeatureBulkCloneResponse {
  results: BulkCloneFeatureResult[];
}

export const featureBulkClone = async (
  params: FeatureBulkCloneParams
): Promise<FeatureBulkCloneResponse> => {
  return axiosClient
    .post<FeatureBulkCloneResponse>(`/v1/feature/bulk_clone`, params)
    .then(response => response.data);
};
