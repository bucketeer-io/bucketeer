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
  CUSTOM = 'CUSTOM',
  DEFAULT = 'DEFAULT'
}

export interface IndividualRuleItem {
  variationId: string;
  name?: string;
  users: string[];
}

export interface VariationPercent {
  variation: string;
  weight?: number | null;
  variationIndex?: number;
  variationId?: string;
}
export interface VariationFeatures {
  label: string;
  value: string;
}
export interface RuleOrders {
  variations: VariationPercent[][];
  labels: string[][];
  isNewRules?: boolean[];
}

export interface DiscardChangesStateData {
  labelType: 'ADD' | 'UPDATE' | 'REMOVE' | 'REORDER';
  label: string;
  groupLabel?: string[];
  featureId?: string;
  variationIndex: number;
  valueLabel?: string;
  changeType?:
    | 'value'
    | 'clause'
    | 'strategy'
    | 'default-strategy'
    | 'default-audience'
    | 'audience'
    | 'new-rule';
  ruleOrders?: RuleOrders;
  variation?: FeatureVariation;
  ruleIndex?: number;
  audienceExcluded?: VariationPercent;
  audienceIncluded?: VariationPercent[];
  variationPercent?: VariationPercent[];
  isAddNew?: boolean;
}

export enum DiscardChangesType {
  PREREQUISITE = 'prerequisite',
  INDIVIDUAL = 'individual',
  CUSTOM = 'custom',
  DEFAULT = 'default'
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

export interface ClauseLabel {
  fullLabel: string;
  valueLabel: string;
}
