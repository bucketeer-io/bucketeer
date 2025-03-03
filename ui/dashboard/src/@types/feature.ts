export type FeatureRuleClauseOperator =
  | 'EQUALS'
  | 'IN'
  | 'ENDS_WITH'
  | 'STARTS_WITH'
  | 'SEGMENT'
  | 'GREATER'
  | 'GREATER_OR_EQUAL'
  | 'LESS'
  | 'LESS_OR_EQUAL'
  | 'BEFORE'
  | 'AFTER'
  | 'FEATURE_FLAG'
  | 'PARTIALLY_MATCH';

export type FeatureVariationType = 'STRING' | 'BOOLEAN' | 'NUMBER' | 'JSON';

export interface FeatureVariation {
  id: string;
  value: string;
  name: string;
  description: string;
}

export interface FeatureTarget {
  variation: string;
  users: string[];
}

export interface FeaturePrerequisite {
  featureId: string;
  variationId: string;
}

export interface RuleStrategyVariation {
  variation: string;
  weight: number;
}

export interface FeatureRuleStrategy {
  type: 'FIXED' | 'ROLLOUT';
  fixedStrategy: {
    variation: string;
  };
  rolloutStrategy: {
    variations: RuleStrategyVariation[];
  };
}

export interface FeatureRuleClause {
  id: string;
  attribute: string;
  operator: FeatureRuleClauseOperator;
  values: string[];
}

export interface FeatureRule {
  id: string;
  strategy: FeatureRuleStrategy;
  clauses: FeatureRuleClause[];
}

export interface Feature {
  id: string;
  name: string;
  description: string;
  enabled: true;
  deleted: true;
  evaluationUndelayable: true;
  ttl: number;
  version: number;
  createdAt: string;
  updatedAt: string;
  variations: FeatureVariation[];
  targets: FeatureTarget[];
  rules: FeatureRule[];
  defaultStrategy: FeatureRuleStrategy;
  offVariation: string;
  tags: string[];
  lastUsedInfo: {
    featureId: string;
    version: 0;
    lastUsedAt: string;
    createdAt: string;
    clientOldestVersion: string;
    clientLatestVersion: string;
  };
  maintainer: string;
  variationType: FeatureVariationType;
  archived: true;
  prerequisites: FeaturePrerequisite[];
  samplingSeed: string;
}

export interface FeatureCollection {
  goals: Array<Feature>;
  cursor: string;
  totalCount: string;
}
