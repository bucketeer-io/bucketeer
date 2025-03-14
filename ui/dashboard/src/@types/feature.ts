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
  autoOpsSummary: AutoOpsSummary;
}

export interface AutoOpsSummary {
  progressiveRolloutCount: number;
  scheduleCount: number;
  killSwitchCount: number;
}

export interface FeatureCollection {
  features: Array<Feature>;
  cursor: string;
  totalCount: string;
}

interface ChangeType {
  changeType: FeatureChangeType;
}
interface VariationChange extends ChangeType {
  variation: FeatureVariation;
}

interface RuleChange extends ChangeType {
  rule: FeatureRule;
}

interface PrerequisiteChange extends ChangeType {
  prerequisite: FeaturePrerequisite;
}

interface TargetChange extends ChangeType {
  target: FeatureTarget;
}

interface TagChange extends ChangeType {
  tag: string;
}

export interface FeatureUpdaterParams {
  comment: string;
  environmentId: string;
  id: string;
  name: string;
  description: string;
  tags: {
    values: string[];
  };
  enabled: boolean;
  archived: boolean;
  variations: {
    values: FeatureVariation[];
  };
  prerequisites: {
    values: FeaturePrerequisite[];
  };
  targets: {
    values: FeatureTarget[];
  };
  rules: {
    values: FeatureRule[];
  };
  defaultStrategy: FeatureRuleStrategy;
  offVariation: string;
  resetSamplingSeed: boolean;
  applyScheduleUpdate: boolean;
  variationChanges: VariationChange[];
  ruleChanges: RuleChange[];
  prerequisiteChanges: PrerequisiteChange[];
  targetChanges: TargetChange[];
  tagChanges: TagChange[];
}
