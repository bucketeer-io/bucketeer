import axiosClient from '@api/axios-client';
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

const fetchTimeSeries = async (
  endpoint: string,
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> => {
  const requestParams = Object.fromEntries(
    Object.entries(params).filter(([, v]) => isNotEmpty(v))
  );
  const response = await axiosClient.get<InsightsTimeSeriesResponse>(
    `/v1/insights/${endpoint}?${stringifyParams(requestParams)}`
  );
  return response.data;
};

export const insightsLatencyFetcher = async (
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> => fetchTimeSeries('latency', params);

export const insightsRequestsFetcher = async (
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> => fetchTimeSeries('requests', params);

export const insightsEvaluationsFetcher = async (
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> =>
  fetchTimeSeries('evaluations', params);

export const insightsErrorRatesFetcher = async (
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> =>
  fetchTimeSeries('error_rates', params);
