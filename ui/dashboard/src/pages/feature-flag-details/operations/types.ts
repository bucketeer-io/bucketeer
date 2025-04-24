import { AutoOpsRule, Rollout } from '@types';

export enum ActionTypeMap {
  UNKNOWN = 'UNKNOWN',
  ENABLE = 'ENABLE',
  DISABLE = 'DISABLE'
}

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
  COMPLETED = 'COMPLETED'
}

export type OperationCombinedType = Rollout & AutoOpsRule;
