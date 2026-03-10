import {
  InsightsMonthlySummaryFetcherParams,
  InsightsTimeSeriesFetcherParams,
  insightsEvaluationsFetcher,
  insightsErrorRatesFetcher,
  insightsLatencyFetcher,
  insightsMonthlySummaryFetcher,
  insightsRequestsFetcher
} from '@api/insight';
import { useQuery } from '@tanstack/react-query';
import type {
  InsightsMonthlySummaryResponse,
  InsightsTimeSeriesResponse,
  QueryOptionsRespond
} from '@types';

// --- Query Keys ---

export const INSIGHTS_MONTHLY_SUMMARY_QUERY_KEY = 'insights-monthly-summary';
export const INSIGHTS_LATENCY_QUERY_KEY = 'insights-latency';
export const INSIGHTS_REQUESTS_QUERY_KEY = 'insights-requests';
export const INSIGHTS_EVALUATIONS_QUERY_KEY = 'insights-evaluations';
export const INSIGHTS_ERROR_RATES_QUERY_KEY = 'insights-error-rates';

// --- Monthly Summary ---

type MonthlySummaryQueryOptions =
  QueryOptionsRespond<InsightsMonthlySummaryResponse> & {
    params?: InsightsMonthlySummaryFetcherParams;
  };

export const useQueryInsightsMonthlySummary = (
  options?: MonthlySummaryQueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  return useQuery({
    queryKey: [INSIGHTS_MONTHLY_SUMMARY_QUERY_KEY, params],
    queryFn: () => insightsMonthlySummaryFetcher(params),
    ...queryOptions
  });
};

// --- Time Series ---

type TimeSeriesQueryOptions =
  QueryOptionsRespond<InsightsTimeSeriesResponse> & {
    params: InsightsTimeSeriesFetcherParams;
  };

export const useQueryInsightsLatency = (options: TimeSeriesQueryOptions) => {
  const { params, ...queryOptions } = options;
  return useQuery({
    queryKey: [INSIGHTS_LATENCY_QUERY_KEY, params],
    queryFn: () => insightsLatencyFetcher(params),
    ...queryOptions
  });
};

export const useQueryInsightsRequests = (options: TimeSeriesQueryOptions) => {
  const { params, ...queryOptions } = options;
  return useQuery({
    queryKey: [INSIGHTS_REQUESTS_QUERY_KEY, params],
    queryFn: () => insightsRequestsFetcher(params),
    ...queryOptions
  });
};

export const useQueryInsightsEvaluations = (
  options: TimeSeriesQueryOptions
) => {
  const { params, ...queryOptions } = options;
  return useQuery({
    queryKey: [INSIGHTS_EVALUATIONS_QUERY_KEY, params],
    queryFn: () => insightsEvaluationsFetcher(params),
    ...queryOptions
  });
};

export const useQueryInsightsErrorRates = (options: TimeSeriesQueryOptions) => {
  const { params, ...queryOptions } = options;
  return useQuery({
    queryKey: [INSIGHTS_ERROR_RATES_QUERY_KEY, params],
    queryFn: () => insightsErrorRatesFetcher(params),
    ...queryOptions
  });
};
