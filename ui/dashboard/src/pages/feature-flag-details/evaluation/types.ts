import { EvaluationTimeRange } from '@types';

export enum EvaluationTab {
  EVENT_COUNT = 'EVENT_COUNT',
  USER_COUNT = 'USER_COUNT'
}

export interface EvaluationFilters {
  period: EvaluationTimeRange;
  tab: EvaluationTab;
}

export interface TimeRangeOption {
  label: string;
  value: EvaluationTimeRange;
}

export interface RawPoint {
  x: Date;
  y: number;
  raw: number;
}
