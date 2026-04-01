import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import {
  InsightApiId,
  InsightSourceId,
  InsightsTimeSeriesResponse
} from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';
import mockData from './mock-data.json';

export interface InsightsTimeSeriesFetcherParams {
  environmentIds?: string[];
  sourceIds?: InsightSourceId[];
  apiIds?: InsightApiId[];
  startAt: string;
  endAt: string;
}

// Re-stamp mock data timestamps to span [startAt, endAt] so the x-axis always
// matches the selected time range filter (hour or day granularity).
const withDynamicTimestamps = (
  response: InsightsTimeSeriesResponse,
  params: InsightsTimeSeriesFetcherParams
): InsightsTimeSeriesResponse => {
  const start = Number(params.startAt);
  const end = Number(params.endAt);
  const durationSecs = end - start;

  // Choose step: <=24h → hourly (3600s), otherwise daily (86400s)
  const step = durationSecs <= 24 * 3600 ? 3600 : 86400;
  const count = Math.round(durationSecs / step);

  return {
    timeseries: response.timeseries.map(series => {
      const srcValues = series.data.map(d => d.value);
      return {
        ...series,
        data: Array.from({ length: count }, (_, i) => ({
          timestamp: String(start + i * step),
          value: srcValues[i % srcValues.length]
        }))
      };
    })
  };
};

const fetchTimeSeries = async (
  endpoint: string,
  params: InsightsTimeSeriesFetcherParams,
  mockFallback: InsightsTimeSeriesResponse
): Promise<InsightsTimeSeriesResponse> => {
  const requestParams = pickBy(params, v => isNotEmpty(v));

  try {
    const response = await axiosClient.get<InsightsTimeSeriesResponse>(
      `/v1/insights/${endpoint}?${stringifyParams(requestParams)}`
    );
    if (response.data) {
      return response.data;
    }
  } catch {
    // fall through to mock
  }

  return withDynamicTimestamps(mockFallback, params);
};

export const insightsLatencyFetcher = async (
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> => {
  return fetchTimeSeries(
    'latency',
    params,
    mockData.latency as InsightsTimeSeriesResponse
  );
};

export const insightsRequestsFetcher = async (
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> => {
  return fetchTimeSeries(
    'requests',
    params,
    mockData.requests as InsightsTimeSeriesResponse
  );
};

export const insightsEvaluationsFetcher = async (
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> => {
  return fetchTimeSeries(
    'evaluations',
    params,
    mockData.evaluations as InsightsTimeSeriesResponse
  );
};

export const insightsErrorRatesFetcher = async (
  params: InsightsTimeSeriesFetcherParams
): Promise<InsightsTimeSeriesResponse> => {
  return fetchTimeSeries(
    'error_rates',
    params,
    mockData.errorRates as InsightsTimeSeriesResponse
  );
};
