export enum TimeseriesUnit {
  HOUR = 'HOUR',
  DAY = 'DAY',
  MINUTE = 'MINUTE'
}

export enum EvaluationTimeRange {
  UNKNOWN = 'UNKNOWN',
  TWENTY_FOUR_HOURS = 'TWENTY_FOUR_HOURS',
  SEVEN_DAYS = 'SEVEN_DAYS',
  FOURTEEN_DAYS = 'FOURTEEN_DAYS',
  THIRTY_DAYS = 'THIRTY_DAYS',
  ONE_HOUR = 'ONE_HOUR',
  SIXTY_DAYS = 'SIXTY_DAYS',
  NINETY_DAYS = 'NINETY_DAYS'
}

export interface EvaluationCollection {
  userCounts: EvaluationCounter[];
  eventCounts: EvaluationCounter[];
}

export interface EvaluationCounter {
  variationId: string;
  timeseries: EventCounterTimeseries;
}

export interface EventCounterTimeseries {
  timestamps: string[];
  values: number[];
  unit: TimeseriesUnit;
  totalCounts: string;
}
