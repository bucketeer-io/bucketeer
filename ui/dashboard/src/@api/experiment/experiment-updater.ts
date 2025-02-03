import axiosClient from '@api/axios-client';
import { Experiment, ExperimentStatus } from '@types';

export interface ExperimentUpdaterParams {
  id: string;
  environmentId: string;
  name?: string;
  description?: string;
  startAt?: string;
  stopAt?: string;
  archived?: boolean;
  status?: {
    status: ExperimentStatus;
  };
}

export interface ExperimentUpdaterResponse {
  experiment: Experiment;
}

export const experimentUpdater = async (
  params?: ExperimentUpdaterParams
): Promise<ExperimentUpdaterResponse> => {
  return axiosClient
    .patch<ExperimentUpdaterResponse>('/v1/experiment', params)
    .then(response => response.data);
};
