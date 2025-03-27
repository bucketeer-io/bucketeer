import { AnyObject } from 'yup';
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

export interface Rollout {
  id: string;
  featureId: string;
  clause: AnyObject;
  status: OperationStatus;
  createdAt: string;
  updatedAt: string;
  type: RolloutType;
  stoppedBy: RolloutStoppedBy;
  stoppedAt: string;
}
