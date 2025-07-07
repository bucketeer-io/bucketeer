export enum RuleClauseType {
  COMPARE = 'compare',
  SEGMENT = 'segment',
  DATE = 'date',
  FEATURE_FLAG = 'feature-flag'
}

export interface PrerequisiteSchema {
  featureId: string;
  variationId: string;
  id?: string;
}

export enum RuleCategory {
  PREREQUISITE = 'PREREQUISITE',
  INDIVIDUAL = 'INDIVIDUAL',
  CUSTOM = 'CUSTOM'
}

export interface IndividualRuleItem {
  variationId: string;
  name?: string;
  users: string[];
  id?: string;
}

export enum DiscardChangesType {
  PREREQUISITE = 'prerequisite',
  INDIVIDUAL = 'individual',
  CUSTOM = 'custom'
}

export enum DiscardChangesType {
  PREREQUISITE = 'prerequisite',
  INDIVIDUAL = 'individual',
  CUSTOM = 'custom'
}
