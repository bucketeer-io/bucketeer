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
  firstValue?: string;
  conditioner: string;
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
  on: string[];
  off: string[];
}

export interface PrerequisiteRuleType {
  featureFlag: string;
  variation: string;
}

export interface TargetSegmentItem {
  index: number;
  rules: SegmentRuleItem[];
}

export interface TargetIndividualItem {
  index: number;
  rules: IndividualRuleItem[];
}

export interface TargetPrerequisiteItem {
  index: number;
  rules: PrerequisiteRuleType[];
}
