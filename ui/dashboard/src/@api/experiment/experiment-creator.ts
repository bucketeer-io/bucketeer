import axiosClient from '@api/axios-client';
import { Experiment } from '@types';

export interface ExperimentCreatorParams {
  environmentId: string;
  featureId: string;
  startAt: string;
  stopAt: string;
  goalIds: string[];
  name: string;
  baseVariationId: string;
  description?: string;
}

export interface ExperimentCreatorResponse {
  experiment: Array<Experiment>;
}

export const experimentCreator = async (
  params?: ExperimentCreatorParams
): Promise<ExperimentCreatorResponse> => {
  return axiosClient
    .post<ExperimentCreatorResponse>('/v1/experiment', params)
    .then(response => response.data);
};
