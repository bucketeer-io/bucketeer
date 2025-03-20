import { RuleSchema } from './form-schema';

export type SituationType =
  | 'compare'
  | 'user-segment'
  | 'date'
  | 'feature-flag';

export type RuleCategory =
  | 'target-segments'
  | 'target-individuals'
  | 'set-prerequisites';

export interface SegmentConditionType {
  situation: SituationType;
  conditioner: string;
  firstValue?: string;
  secondValue?: string;
  value?: string;
  date?: string;
  flagId?: string;
  variation?: string;
}

export interface SegmentRuleItem {
  conditions: SegmentConditionType[];
  variation: boolean;
}

export interface IndividualRuleItem {
  variationId: string;
  name?: string;
  users: string[];
}

export interface PrerequisiteRuleType {
  featureFlag: string;
  variation: string;
}

export interface TargetSegmentItem {
  index: number;
  rules: SegmentRuleItem[];
}

export interface TargetPrerequisiteItem {
  index: number;
  rules: PrerequisiteRuleType[];
}

export interface TargetingForm {
  rules?: RuleSchema[];
  prerequisitesRules: TargetPrerequisiteItem[];
  targetIndividualRules: IndividualRuleItem[];
}
