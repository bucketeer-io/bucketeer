import axiosClient from '@api/axios-client';
import { Feature, FeatureVariation, FeatureVariationType } from '@types';

export interface FeatureCreatorParams {
  environmentId: string;
  id: string;
  name: string;
  variations: FeatureVariation[];
  tags: string[];
  defaultOnVariationIndex: number;
  defaultOffVariationIndex: number;
  variationType: FeatureVariationType;
  description?: string;
}

export interface FeatureCreatorResponse {
  feature: Feature;
}

export const featureCreator = async (
  params?: FeatureCreatorParams
): Promise<FeatureCreatorResponse> => {
  return axiosClient
    .post<FeatureCreatorResponse>('/v1/feature', params)
    .then(response => response.data);
};
