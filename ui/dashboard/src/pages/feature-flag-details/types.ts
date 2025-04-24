export interface TabItem {
  readonly title: string;
  readonly to: string;
}

export type OperationActionType =
  | 'NEW'
  | 'UPDATE'
  | 'DETAILS'
  | 'STOP'
  | 'DELETE';
