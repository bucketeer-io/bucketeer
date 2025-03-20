export type ScheduleChangeType =
  | 'CHANGE_UNSPECIFIED'
  | 'CHANGE_CREATE'
  | 'CHANGE_UPDATE'
  | 'CHANGE_DELETE';

export type ScheduleFieldType =
  | 'UNSPECIFIED'
  | 'PREREQUISITES'
  | 'TARGETS'
  | 'RULES'
  | 'DEFAULT_STRATEGY'
  | 'OFF_VARIATION'
  | 'VARIATIONS';

export interface ScheduleChangeItem {
  id: string;
  changeType: ScheduleChangeType;
  fieldType: ScheduleFieldType;
  fieldValue: string;
}

export interface ScheduleFlagListItem {
  id: string;
  featureId: string;
  environmentId: string;
  scheduledAt: string;
  createdAt: string;
  updatedAt: string;
  changes: ScheduleChangeItem[];
}

export interface ScheduleFlagCollection {
  scheduledFlagUpdates: ScheduleFlagListItem[];
}
