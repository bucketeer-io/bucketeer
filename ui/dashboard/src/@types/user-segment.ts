export interface UserSegmentRule {
  id: string;
  strategy: UserSegmentRuleStrategy;
  clauses: UserSegmentRuleClauses[];
}

export interface RolloutStrategy {
  variations: RolloutStrategyVariations[];
}

export interface RolloutStrategyVariations {
  variation: string;
  weight: number;
}

export type UserSegmentRuleStrategy = {
  type: 'FIXED' | 'ROLLOUT';
  fixedStrategy: { variation: string };
  rolloutStrategy: RolloutStrategy;
};

export type ClauseOperator =
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

export interface UserSegmentRuleClauses {
  id: string;
  attribute: string;
  operator: ClauseOperator;
  values: string[];
}

export type FeatureSegmentStatus =
  | 'INITIAL'
  | 'UPLOADING'
  | 'SUCEEDED'
  | 'FAILED';

export interface FeatureVariations {
  id: string;
  value: string;
  name: string;
  description: string;
}

export interface FeatureTarget {
  variation: string;
  users: string[];
}

export interface FeatureLastUsedInfo {
  featureId: string;
  version: number;
  lastUsedAt: string;
  createdAt: string;
  clientOldestVersion: string;
  clientLatestVersion: string;
}

export interface FeaturePrerequisites {
  featureId: string;
  variationId: string;
}

export interface UserSegmentFeature {
  id: string;
  name: string;
  description: string;
  enabled: true;
  deleted: true;
  evaluationUndelayable: true;
  ttl: 0;
  version: 0;
  createdAt: string;
  updatedAt: string;
  variations: FeatureVariations[];
  targets: FeatureTarget[];
  rules: UserSegmentRule[];
  defaultStrategy: UserSegmentRuleStrategy;
  offVariation: string;
  tags: string[];
  lastUsedInfo: FeatureLastUsedInfo;
  maintainer: string;
  variationType: string;
  archived: true;
  prerequisites: FeaturePrerequisites[];
  samplingSeed: string;
}

export interface UserSegment {
  id: string;
  name: string;
  description: string;
  rules: UserSegmentRule[];
  createdAt: string;
  updatedAt: string;
  version: string;
  deleted: true;
  includedUserCount: string;
  excludedUserCount: string;
  status: FeatureSegmentStatus;
  isInUseStatus: boolean;
  features: UserSegmentFeature[];
}

export interface UserSegmentCollection {
  segments: Array<UserSegment>;
  cursor: string;
  totalCount: string;
}
