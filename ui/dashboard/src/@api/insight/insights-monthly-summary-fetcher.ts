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
  _params?: InsightsMonthlySummaryFetcherParams
): Promise<InsightsMonthlySummaryResponse> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  const response = await axiosClient.get<InsightsMonthlySummaryResponse>(
    `/v1/insights/monthly_summary?${stringifyParams(params)}`
  );
  return response.data;
};
