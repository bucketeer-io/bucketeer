import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { Experiment } from '@types';
import { isNotEmpty } from 'utils/data-type';

export interface ExperimentFetcherParams {
  id: string;
  environmentId: string;
}

export interface ExperimentResponse {
  experiment: Experiment;
}

export const experimentFetcher = async (
  params?: ExperimentFetcherParams
): Promise<ExperimentResponse> => {
  return axiosClient
    .get<ExperimentResponse>('/v1/experiment', {
      params: pickBy(params, v => isNotEmpty(v))
    })
    .then(response => response.data);
};
