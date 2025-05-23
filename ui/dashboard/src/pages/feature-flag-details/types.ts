export interface TabItem {
  readonly title: string;
  readonly to: string;
}

export type OperationActionType =
  | 'NEW'
  | 'UPDATE'
  | 'DETAILS'
  | 'STOP'
  | 'DELETE'
  | 'CLONE';

export interface ScheduleItem {
  scheduleId?: string;
  executeAt: Date;
  weight: number;
  triggeredAt?: string;
}
