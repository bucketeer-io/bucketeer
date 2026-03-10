import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import {
  InsightApiId,
  InsightSourceId,
  InsightsTimeSeriesResponse
} from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface InsightsTimeSeriesFetcherParams {
  environmentIds?: string[];
  sourceIds?: InsightSourceId[];
  apiIds?: InsightApiId[];
  startAt: string;
  endAt: string;
}

export const insightsLatencyFetcher = async (
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));
  return axiosClient
    .get<InsightsTimeSeriesResponse>(`/v1/insights/latency?${requestParams}`)
    .then(response => response.data);
};

export const insightsRequestsFetcher = async (
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));
  return axiosClient
    .get<InsightsTimeSeriesResponse>(`/v1/insights/requests?${requestParams}`)
    .then(response => response.data);
};

export const insightsEvaluationsFetcher = async (
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));
  return axiosClient
    .get<InsightsTimeSeriesResponse>(
      `/v1/insights/evaluations?${requestParams}`
    )
    .then(response => response.data);
};

export const insightsErrorRatesFetcher = async (
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));
  return axiosClient
    .get<InsightsTimeSeriesResponse>(
      `/v1/insights/error_rates?${requestParams}`
    )
    .then(response => response.data);
};
