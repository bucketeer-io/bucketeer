export enum TriggerType {
  UNKNOWN = 'Type_UNKNOWN',
  WEBHOOK = 'Type_WEBHOOK'
}

export enum TriggerActionType {
  UNKNOWN = 'Action_UNKNOWN',
  ON = 'Action_ON',
  OFF = 'Action_OFF'
}

export interface TriggerCollection {
  flagTriggers: TriggerItemType[];
  cursor: string;
  totalCount: string;
}

export interface TriggerItemType {
  flagTrigger: Trigger;
  url: string;
}

export interface Trigger {
  id: string;
  featureId: string;
  type: TriggerType;
  action: TriggerActionType;
  description: string;
  triggerCount: number;
  lastTriggeredAt: string;
  token: string;
  disabled: true;
  createdAt: string;
  updatedAt: string;
  environmentId: string;
}
