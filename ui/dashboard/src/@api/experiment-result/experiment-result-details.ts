import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { ExperimentResultResponse } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface ExperimentResultDetailsFetcherParams {
  experimentId: string;
  environmentId: string;
}

export const experimentResultDetailsFetcher = async (
  params?: ExperimentResultDetailsFetcherParams
): Promise<ExperimentResultResponse> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<ExperimentResultResponse>(`/v1/experiment_result?${requestParams}`)
    .then(response => response.data);
};
