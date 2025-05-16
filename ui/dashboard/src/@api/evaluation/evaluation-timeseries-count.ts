import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { EvaluationCollection, EvaluationTimeRange } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface EvaluationTimeseriesCountParams {
  environmentId: string;
  featureId: string;
  timeRange: EvaluationTimeRange;
}

export const evaluationTimeseriesCount = async (
  _params?: EvaluationTimeseriesCountParams
): Promise<EvaluationCollection> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<EvaluationCollection>(
      `/v1/evaluation_timeseries_count?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
