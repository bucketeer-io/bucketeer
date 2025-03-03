import { Feature, FeatureRule } from './feature';

export type FeatureSegmentStatus =
  | 'INITIAL'
  | 'UPLOADING'
  | 'SUCEEDED'
  | 'FAILED';

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
