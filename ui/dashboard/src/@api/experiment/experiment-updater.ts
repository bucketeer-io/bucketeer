import axiosClient from '@api/axios-client';
import { ExperimentStatus } from '@types';
import { ExperimentCreateUpdateResponse } from './experiment-creator';

export interface ExperimentUpdaterParams {
  id?: string;
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

export const experimentUpdater = async (
  params?: ExperimentUpdaterParams
): Promise<ExperimentCreateUpdateResponse> => {
  return axiosClient
    .patch<ExperimentCreateUpdateResponse>('/v1/experiment', params)
    .then(response => response.data);
};
