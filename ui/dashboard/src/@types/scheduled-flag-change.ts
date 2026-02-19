import type {
  FeatureRuleStrategy,
  FeaturePrerequisite,
  FeatureTarget,
  FeatureVariation,
  FeatureRuleChange,
  FeatureChangeType
} from './feature';

export const ScheduledFlagChangeStatuses = {
  UNSPECIFIED: 'SCHEDULED_FLAG_CHANGE_STATUS_UNSPECIFIED',
  PENDING: 'SCHEDULED_FLAG_CHANGE_STATUS_PENDING',
  EXECUTED: 'SCHEDULED_FLAG_CHANGE_STATUS_EXECUTED',
  FAILED: 'SCHEDULED_FLAG_CHANGE_STATUS_FAILED',
  CANCELLED: 'SCHEDULED_FLAG_CHANGE_STATUS_CANCELLED',
  CONFLICT: 'SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT'
} as const;

export type ScheduledFlagChangeStatus =
  (typeof ScheduledFlagChangeStatuses)[keyof typeof ScheduledFlagChangeStatuses];

export const ScheduledChangeCategories = {
  UNSPECIFIED: 'SCHEDULED_CHANGE_CATEGORY_UNSPECIFIED',
  TARGETING: 'SCHEDULED_CHANGE_CATEGORY_TARGETING',
  VARIATIONS: 'SCHEDULED_CHANGE_CATEGORY_VARIATIONS',
  SETTINGS: 'SCHEDULED_CHANGE_CATEGORY_SETTINGS',
  MIXED: 'SCHEDULED_CHANGE_CATEGORY_MIXED'
} as const;

export type ScheduledChangeCategory =
  (typeof ScheduledChangeCategories)[keyof typeof ScheduledChangeCategories];

export type ScheduledChangeConflictType =
  | 'CONFLICT_TYPE_UNSPECIFIED'
  | 'CONFLICT_TYPE_VERSION_MISMATCH'
  | 'CONFLICT_TYPE_OVERLAPPING_SCHEDULE'
  | 'CONFLICT_TYPE_DEPENDENCY_MISSING'
  | 'CONFLICT_TYPE_INVALID_REFERENCE';

export interface ScheduledChangePayloadVariationChange {
  changeType: FeatureChangeType;
  variation: FeatureVariation;
}

export interface ScheduledChangePayloadRuleChange {
  changeType: FeatureChangeType;
  rule: FeatureRuleChange;
}

export interface ScheduledChangePayloadPrerequisiteChange {
  changeType: FeatureChangeType;
  prerequisite: FeaturePrerequisite;
}

export interface ScheduledChangePayloadTargetChange {
  changeType: FeatureChangeType;
  target: FeatureTarget;
}

export interface ScheduledChangePayloadTagChange {
  changeType: FeatureChangeType;
  tag: string;
}

export interface ScheduledChangePayload {
  ruleChanges?: ScheduledChangePayloadRuleChange[];
  targetChanges?: ScheduledChangePayloadTargetChange[];
  prerequisiteChanges?: ScheduledChangePayloadPrerequisiteChange[];
  defaultStrategy?: Partial<FeatureRuleStrategy>;
  variationChanges?: ScheduledChangePayloadVariationChange[];
  offVariation?: string;
  enabled?: boolean;
  name?: string;
  description?: string;
  tagChanges?: ScheduledChangePayloadTagChange[];
  archived?: boolean;
  resetSamplingSeed?: boolean;
  maintainer?: string;
}

export interface ChangeSummary {
  messageKey: string;
  values: Record<string, string>;
}

export interface ScheduledChangeConflict {
  type: ScheduledChangeConflictType;
  description: string;
  conflictingScheduleId: string;
  conflictingField: string;
  detectedAt: string;
}

export interface ScheduledFlagChange {
  id: string;
  featureId: string;
  environmentId: string;
  scheduledAt: string;
  timezone: string;
  payload: ScheduledChangePayload;
  comment: string;
  status: ScheduledFlagChangeStatus;
  failureReason: string;
  flagVersionAtCreation: number;
  conflicts: ScheduledChangeConflict[];
  createdBy: string;
  createdAt: string;
  updatedBy: string;
  updatedAt: string;
  executedAt: string;
  category: ScheduledChangeCategory;
  changeSummaries: ChangeSummary[];
}

export interface ScheduledFlagChangeSummary {
  featureId: string;
  pendingCount: number;
  conflictCount: number;
  nextScheduledAt: string;
  nextCategory: ScheduledChangeCategory;
}

export interface ScheduledFlagChangeCollection {
  scheduledFlagChanges: ScheduledFlagChange[];
  cursor: string;
  totalCount: string;
}

export interface CreateScheduledFlagChangeResponse {
  scheduledFlagChange: ScheduledFlagChange;
  detectedConflicts: ScheduledChangeConflict[];
}

export interface UpdateScheduledFlagChangeResponse {
  scheduledFlagChange: ScheduledFlagChange;
  detectedConflicts: ScheduledChangeConflict[];
}

export interface GetScheduledFlagChangeResponse {
  scheduledFlagChange: ScheduledFlagChange;
}

export interface ExecuteScheduledFlagChangeResponse {
  scheduledFlagChange: ScheduledFlagChange;
}

export interface GetScheduledFlagChangeSummaryResponse {
  summary: ScheduledFlagChangeSummary;
}
