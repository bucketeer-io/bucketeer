import { OperationStatus } from './goal';

export interface RolloutCollection {
  progressiveRollouts: Rollout[];
  cursor: string;
  totalCount: string;
}

export type RolloutType = 'MANUAL_SCHEDULE' | 'TEMPLATE_SCHEDULE';
export type RolloutStoppedBy =
  | 'UNKNOWN'
  | 'USER'
  | 'OPS_SCHEDULE'
  | 'OPS_KILL_SWITCH';
export type IntervalType = 'UNKNOWN' | 'HOURLY' | 'DAILY' | 'WEEKLY';
export interface Rollout {
  id: string;
  featureId: string;
  clause: RolloutClause;
  status: OperationStatus;
  createdAt: string;
  updatedAt: string;
  type: RolloutType;
  stoppedBy: RolloutStoppedBy;
  stoppedAt: string;
}

export interface RolloutClause {
  schedules: RolloutSchedule[];
  interval?: IntervalType;
  increments?: string;
  variationId: string;
}

export interface RolloutSchedule {
  scheduleId: string;
  executeAt: string;
  weight: number;
  triggeredAt?: string;
}
export interface RolloutManualScheduleClause {
  schedules: RolloutSchedule[];
  variationId: string;
}

export interface RolloutTemplateScheduleClause {
  schedules: RolloutSchedule[];
  interval: IntervalType;
  increments: string;
  variationId: string;
}
