import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { InsightSourceId, InsightsMonthlySummaryResponse } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface InsightsMonthlySummaryFetcherParams {
  environmentIds?: string[];
  sourceIds?: InsightSourceId[];
}

export const insightsMonthlySummaryFetcher = async (
  params?: InsightsMonthlySummaryFetcherParams
): Promise<InsightsMonthlySummaryResponse> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));
  return axiosClient
    .get<InsightsMonthlySummaryResponse>(
      `/v1/insights/monthly_summary?${requestParams}`
    )
    .then(response => response.data);
};
