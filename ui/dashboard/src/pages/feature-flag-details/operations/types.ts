import { AutoOpsRule, RecurrenceFrequency, Rollout } from '@types';

export enum ActionTypeMap {
  UNKNOWN = 'UNKNOWN',
  ENABLE = 'ENABLE',
  DISABLE = 'DISABLE'
}

export type OperationActionType =
  | 'NEW'
  | 'UPDATE'
  | 'DETAILS'
  | 'STOP'
  | 'DELETE'
  | 'CLONE';

export enum OpsTypeMap {
  // TYPE_UNKNOWN = 'TYPE_UNKNOWN',
  SCHEDULE = 'SCHEDULE',
  EVENT_RATE = 'EVENT_RATE',
  ROLLOUT = 'ROLLOUT'
}

export enum IntervalMap {
  UNKNOWN = 'UNKNOWN',
  HOURLY = 'HOURLY',
  DAILY = 'DAILY',
  WEEKLY = 'WEEKLY'
}

export enum RolloutTypeMap {
  MANUAL_SCHEDULE = 'MANUAL_SCHEDULE',
  TEMPLATE_SCHEDULE = 'TEMPLATE_SCHEDULE'
}

export enum OperationTab {
  ACTIVE = 'ACTIVE',
  FINISHED = 'FINISHED'
}

export enum ScheduleType {
  ONE_TIME = 'ONE_TIME',
  RECURRING = 'RECURRING'
}

export enum EndConditionType {
  NEVER = 'NEVER',
  ON_DATE = 'ON_DATE',
  AFTER = 'AFTER'
}

export type OperationCombinedType = Rollout & AutoOpsRule;

export interface ScheduleItem {
  scheduleId?: string;
  executeAt: Date;
  weight: number;
  triggeredAt?: string;
}

export const DAYS_OF_WEEK = [0, 1, 2, 3, 4, 5, 6] as const;
export const DAY_LABELS_SHORT_KEYS = [
  'feature-flags.day-short-sun',
  'feature-flags.day-short-mon',
  'feature-flags.day-short-tue',
  'feature-flags.day-short-wed',
  'feature-flags.day-short-thu',
  'feature-flags.day-short-fri',
  'feature-flags.day-short-sat'
] as const;
export const DAY_LABELS_FULL = [
  'sunday',
  'monday',
  'tuesday',
  'wednesday',
  'thursday',
  'friday',
  'saturday'
] as const;

export const FREQUENCY_OPTIONS: RecurrenceFrequency[] = [
  'DAILY',
  'WEEKLY',
  'MONTHLY'
];
