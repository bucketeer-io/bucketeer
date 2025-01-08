export type UserSegmentRule = {
  id: string;
  strategy: UserSegmentRuleStrategy;
  clauses: UserSegmentRuleClauses[];
};

export type FixedStrategy = {
  variation: string;
};
export type RolloutStrategy = {
  variations: RolloutStrategyVariations[];
};
export type RolloutStrategyVariations = FixedStrategy & {
  weight: number;
};
export type UserSegmentRuleStrategy = {
  type: 'FIXED' | 'ROLLOUT';
  fixedStrategy: FixedStrategy;
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

export type UserSegmentRuleClauses = {
  id: string;
  attribute: string;
  operator: ClauseOperator;
  values: string[];
};

export type FeatureSegmentStatus =
  | 'INITIAL'
  | 'UPLOADING'
  | 'SUCEEDED'
  | 'FAILED';

export type FeatureVariations = {
  id: string;
  value: string;
  name: string;
  description: string;
};

export type FeatureTarget = FixedStrategy & {
  users: string[];
};

export type FeatureLastUsedInfo = {
  featureId: string;
  version: number;
  lastUsedAt: string;
  createdAt: string;
  clientOldestVersion: string;
  clientLatestVersion: string;
};

export type FeaturePrerequisites = {
  featureId: string;
  variationId: string;
};

export type UserSegmentFeature = {
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
};

export type UserSegment = {
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
};

export interface UserSegmentCollection {
  segments: Array<UserSegment>;
  cursor: string;
  totalCount: string;
}