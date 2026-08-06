import { Feature, FeatureRule } from './feature';

export type FeatureSegmentStatus =
  'INITIAL' | 'UPLOADING' | 'SUCEEDED' | 'FAILED';

export interface UserSegment {
  id: string;
  name: string;
  description: string;
  rules: FeatureRule[];
  createdAt: string;
  updatedAt: string;
  version: string;
  deleted: true;
  includedUserCount: string;
  excludedUserCount: string;
  status: FeatureSegmentStatus;
  isInUseStatus: boolean;
  features: Feature[];
}

export interface UserSegmentCollection {
  segments: Array<UserSegment>;
  cursor: string;
  totalCount: string;
}

export interface SegmentRuleClausePayload {
  id: string;
  attribute: string;
  operator: string;
  values: string[];
}

// Segment rules carry no strategy: matching a rule simply means
// "the user is in the segment".
export interface SegmentRulePayload {
  id: string;
  clauses: SegmentRuleClausePayload[];
}

// proto RuleListValue: absent = rules unchanged, present = full replacement.
export interface SegmentRuleListValue {
  values: SegmentRulePayload[];
}
