export enum FeatureRuleClauseOperator {
  EQUALS = 'EQUALS',
  GREATER_OR_EQUAL = 'GREATER_OR_EQUAL',
  GREATER = 'GREATER',
  LESS_OR_EQUAL = 'LESS_OR_EQUAL',
  LESS = 'LESS',
  IN = 'IN',
  PARTIALLY_MATCH = 'PARTIALLY_MATCH',
  STARTS_WITH = 'STARTS_WITH',
  ENDS_WITH = 'ENDS_WITH',
  BEFORE = 'BEFORE',
  AFTER = 'AFTER',
  SEGMENT = 'SEGMENT',
  FEATURE_FLAG = 'FEATURE_FLAG'
}

export type FeatureVariationType = 'STRING' | 'BOOLEAN' | 'NUMBER' | 'JSON';

export type FeatureChangeType = 'UNSPECIFIED' | 'CREATE' | 'UPDATE' | 'DELETE';

export interface FeatureVariation {
  id: string;
  value: string;
  name: string;
  description?: string;
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

export enum StrategyType {
  FIXED = 'FIXED',
  ROLLOUT = 'ROLLOUT',
  MANUAL = 'MANUAL'
}
export enum DefaultRuleStrategyType {
  FIXED = 'FIXED',
  MANUAL = 'MANUAL',
  ROLLOUT = 'ROLLOUT'
}

export interface FeatureRuleStrategy {
  type: StrategyType;
  fixedStrategy: {
    variation: string;
  };
  rolloutStrategy: {
    variations: RuleStrategyVariation[];
    audience: {
      percentage: number;
      defaultVariation: string;
    };
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
  enabled: boolean;
  deleted: boolean;
  evaluationUndelayable: boolean;
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
  archived: boolean;
  prerequisites: FeaturePrerequisite[];
  samplingSeed: string;
  autoOpsSummary: AutoOpsSummary;
}

export interface AutoOpsSummary {
  progressiveRolloutCount: number;
  scheduleCount: number;
  killSwitchCount: number;
}

export interface FeatureCountByStatus {
  active: number;
  inactive: number;
  total: number;
}
export interface FeatureCollection {
  featureCountByStatus: FeatureCountByStatus;
  features: Array<Feature>;
  cursor: string;
  totalCount: string;
}

interface ChangeType {
  changeType: FeatureChangeType;
}
export interface VariationChange extends ChangeType {
  variation: FeatureVariation;
}

export interface FeatureRuleChange {
  id: string;
  strategy: Partial<FeatureRuleStrategy>;
  clauses: FeatureRuleClause[];
}

export interface RuleChange extends ChangeType {
  rule: FeatureRuleChange;
}

export interface PrerequisiteChange extends ChangeType {
  prerequisite: FeaturePrerequisite;
}

export interface TargetChange extends ChangeType {
  target: FeatureTarget;
}

export interface TagChange extends ChangeType {
  tag: string;
}

export interface FeatureUpdaterParams {
  comment: string;
  environmentId: string;
  id: string;
  name: string;
  description: string;
  enabled: boolean;
  archived: boolean;
  defaultStrategy: Partial<FeatureRuleStrategy>;
  offVariation: string;
  resetSamplingSeed: boolean;
  applyScheduleUpdate: boolean;
  variationChanges: VariationChange[];
  ruleChanges: RuleChange[];
  prerequisiteChanges: PrerequisiteChange[];
  targetChanges: TargetChange[];
  tagChanges: TagChange[];
}

export interface AttributeKeysResponse {
  userAttributeKeys: string[];
}
