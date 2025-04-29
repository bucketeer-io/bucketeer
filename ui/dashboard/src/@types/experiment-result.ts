export interface ExperimentResultResponse {
  experimentResult: ExperimentResult;
}

export interface ExperimentResult {
  id: string;
  experimentId: string;
  updatedAt: string;
  goalResults: GoalResult[];
}

export interface GoalResult {
  goalId: string;
  variationResults: VariationResult[];
  summary: GoalResultSummary;
}

export interface GoalResultSummary {
  bestVariations: BestVariation[];
  totalEvaluationUserCount: string;
  totalGoalUserCount: string;
}

export interface BestVariation {
  id: string;
  probability: number;
  isBest: boolean;
}

export interface VariationResult {
  variationId: string;
  variationName?: string;
  experimentCount: Count;
  evaluationCount: Count;
  cvrProbBest: CvrProb;
  cvrProbBeatBaseline: CvrProb;
  cvrProb: CvrProb;
  evaluationUserCountTimeseries: Timeseries;
  evaluationEventCountTimeseries: Timeseries;
  goalUserCountTimeseries: Timeseries;
  goalEventCountTimeseries: Timeseries;
  goalValueSumTimeseries: Timeseries;
  cvrMedianTimeseries: Timeseries;
  cvrPercentile025Timeseries: Timeseries;
  cvrPercentile975Timeseries: Timeseries;
  cvrTimeseries: Timeseries;
  goalValueSumPerUserTimeseries: Timeseries;
  goalValueSumPerUserProb: CvrProb;
  goalValueSumPerUserProbBest: CvrProb;
  goalValueSumPerUserProbBeatBaseline: CvrProb;
  goalValueSumPerUserMedianTimeseries: Timeseries;
  goalValueSumPerUserPercentile025Timeseries: Timeseries;
  goalValueSumPerUserPercentile975Timeseries: Timeseries;
  conversionRate: number;
  expectedLoss: number;
  cvrSamples: number[];
}

export interface Timeseries {
  timestamps: string[];
  values: number[];
  unit: Unit;
  totalCounts: string;
}

export enum Unit {
  Hour = 'HOUR'
}

export interface CvrProb {
  mean: number;
  sd: number;
  rhat: number;
  histogram: Histogram;
  median: number;
  percentile025: number;
  percentile975: number;
}

export interface Histogram {
  hist: string[];
  bins: number[];
}

export interface Count {
  variationId: string;
  userCount: string;
  eventCount: string;
  valueSum: number;
  createdAt: string;
  variationValue: string;
  valueSumPerUserMean: number;
  valueSumPerUserVariance: number;
}
