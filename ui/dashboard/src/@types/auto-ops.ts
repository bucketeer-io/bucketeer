import { OperationStatus } from './goal';

export type AutoOpsType = 'TYPE_UNKNOWN' | 'SCHEDULE' | 'EVENT_RATE';
export type ClauseActionType = 'UNKNOWN' | 'ENABLE' | 'DISABLE';
export type OpsEventRateClauseOperator = 'GREATER_OR_EQUAL' | 'LESS_OR_EQUAL';

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
}
