import { FeatureVariation } from './feature';

export interface DebuggerCollection {
  evaluations: Evaluation[];
  archivedFeatureIds: string[];
}

export interface Evaluation {
  id: string;
  featureId: string;
  featureVersion: number;
  userId: string;
  variationId: string;
  variation: FeatureVariation;
  reason: EvaluationReason;
  variationValue: string;
  variationName: string;
}

export interface EvaluationReason {
  type: ReasonType;
  ruleId: string;
}

export type ReasonType =
  | 'TARGET'
  | 'RULE'
  | 'DEFAULT'
  | 'CLIENT'
  | 'OFF_VARIATION'
  | 'PREREQUISITE'
  | 'ERROR_NO_EVALUATIONS'
  | 'ERROR_FLAG_NOT_FOUND'
  | 'ERROR_WRONG_TYPE'
  | 'ERROR_USER_ID_NOT_SPECIFIED'
  | 'ERROR_FEATURE_FLAG_ID_NOT_SPECIFIED'
  | 'ERROR_EXCEPTION';
