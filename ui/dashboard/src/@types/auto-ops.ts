import { OperationStatus } from './goal';

export type AutoOpsType = 'TYPE_UNKNOWN' | 'SCHEDULE' | 'EVENT_RATE';
export type ClauseActionType = 'UNKNOWN' | 'ENABLE' | 'DISABLE';
export type OpsEventRateClauseOperator = 'GREATER_OR_EQUAL' | 'LESS_OR_EQUAL';
export type AutoOpsChangeType = 'UNSPECIFIED' | 'CREATE' | 'UPDATE' | 'DELETE';
export type RecurrenceFrequency =
  | 'FREQUENCY_UNSPECIFIED'
  | 'ONCE'
  | 'DAILY'
  | 'WEEKLY'
  | 'MONTHLY';

export interface RecurrenceRule {
  frequency: RecurrenceFrequency;
  daysOfWeek: number[];
  dayOfMonth: number;
  startDate: string;
  endDate: string;
  maxOccurrences: number;
  timezone: string;
}

export interface AutoOpsRule {
  id: string;
  featureId: string;
  opsType: AutoOpsType;
  clauses: AutoOpsRuleClause[];
  createdAt: string;
  updatedAt: string;
  deleted: boolean;
  autoOpsStatus: OperationStatus;
  featureName: string;
}

export interface AutoOpsRuleCollection {
  autoOpsRules: AutoOpsRule[];
  cursor: string;
}

export interface AutoOpsRuleClause {
  id: string;
  clause: OpsEventRateClause | DatetimeClause;
  actionType: ClauseActionType;
  executedAt: string;
  isRecurring?: boolean;
}

export interface OpsEventRateClause {
  variationId: string;
  goalId: string;
  minCount: string;
  threadsholdRate: number;
  operator: OpsEventRateClauseOperator;
  actionType: ClauseActionType;
}

export interface DatetimeClause {
  time: string;
  actionType: ClauseActionType;
  recurrence?: RecurrenceRule;
  lastExecutedAt?: string;
  nextExecutionAt?: string;
  executionCount?: number;
}

export interface AutoOpsCountCollection {
  cursor: string;
  opsCounts: AutoOpsCount[];
}

export interface AutoOpsCount {
  id: string;
  autoOpsRuleId: string;
  clauseId: string;
  updatedAt: string;
  opsEventCount: string;
  evaluationCount: string;
  featureId: string;
}
