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
  reason: Reason;
  variationValue: string;
  variationName: string;
}

export interface Reason {
  type: string;
  ruleId: string;
}
