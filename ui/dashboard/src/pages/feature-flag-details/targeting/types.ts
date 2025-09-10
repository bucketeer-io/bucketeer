import { FeatureVariation, PrerequisiteChange } from '@types';

export enum RuleClauseType {
  COMPARE = 'compare',
  SEGMENT = 'segment',
  DATE = 'date',
  FEATURE_FLAG = 'feature-flag'
}

export interface PrerequisiteSchema {
  featureId: string;
  variationId: string;
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
}

export interface DiscardChangesStateData {
  labelType: 'ADD' | 'UPDATE' | 'REMOVE';
  label: string;
  featureId?: string;
  variationIndex: number;
  // rule segment
  valueLabel?: string;
  changeField?: 'value' | 'clause' | 'strategy';
  variation?: FeatureVariation;
  ruleIndex?: number;
  variationPercent?: { variation: string; weight?: number }[];
}

export enum DiscardChangesType {
  PREREQUISITE = 'prerequisite',
  INDIVIDUAL = 'individual',
  CUSTOM = 'custom'
}

export type DiscardFeaturePrerequisiteChange = {
  deleted?: PrerequisiteChange;
  created?: PrerequisiteChange;
};

export interface DiscardChangesState {
  type: DiscardChangesType | undefined;
  isOpen: boolean;
  data: DiscardChangesStateData[];
  ruleIndex?: number;
}
