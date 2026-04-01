export type InsightSourceId =
  | 'UNKNOWN'
  | 'ANDROID'
  | 'IOS'
  | 'WEB'
  | 'GO_SERVER'
  | 'NODE_SERVER'
  | 'JAVASCRIPT'
  | 'FLUTTER'
  | 'REACT'
  | 'REACT_NATIVE'
  | 'OPEN_FEATURE_KOTLIN'
  | 'OPEN_FEATURE_SWIFT'
  | 'OPEN_FEATURE_JAVASCRIPT'
  | 'OPEN_FEATURE_GO'
  | 'OPEN_FEATURE_NODE'
  | 'OPEN_FEATURE_REACT'
  | 'OPEN_FEATURE_REACT_NATIVE';

export type InsightApiId =
  | 'UNKNOWN_API'
  | 'GET_EVALUATION'
  | 'GET_EVALUATIONS'
  | 'REGISTER_EVENTS'
  | 'GET_FEATURE_FLAGS'
  | 'GET_SEGMENT_USERS'
  | 'SDK_GET_VARIATION';

// --- Monthly Summary ---

export interface MonthlySummaryDataPoint {
  yearmonth: string; // e.g. "202601"
  mau: string; // int64 as string
  requests: string; // int64 as string
}

export interface MonthlySummarySeries {
  environmentId: string;
  sourceId: InsightSourceId;
  environmentName: string;
  projectName: string;
  data: MonthlySummaryDataPoint[]; // sorted by yearmonth ascending
}

export interface InsightsMonthlySummaryResponse {
  series: MonthlySummarySeries[];
}

// --- Time Series ---

export interface InsightsDataPoint {
  timestamp: string; // int64 as string (Unix seconds)
  value: number;
}

export interface InsightsTimeSeries {
  environmentId: string;
  sourceId: InsightSourceId;
  apiId: InsightApiId;
  data: InsightsDataPoint[]; // sorted by timestamp ascending
  labels: Record<string, string>; // e.g. { evaluation_type: "diff" }
}

export interface InsightsTimeSeriesResponse {
  timeseries: InsightsTimeSeries[];
}
