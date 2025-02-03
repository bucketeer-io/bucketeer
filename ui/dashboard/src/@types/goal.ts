import { Experiment } from "./experiment";

export type ConnectionType = 'UNKNOWN' | 'EXPERIMENT' | 'OPERATION';

export type OpsType = 'TYPE_UNKNOWN' | 'SCHEDULE' | 'EVENT_RATE';
export type AutoOpsRuleStatus = 'WAITING' | 'RUNNING' | 'FINISHED' | 'STOPPED';

export interface Goal {
  id: string;
  name: string;
  description: string;
  deleted: boolean;
  createdAt: string;
  updatedAt: string;
  isInUseStatus: boolean;
  archived: boolean;
  connectionType: ConnectionType;
  experiments: Experiment[];
  autoOpsRules: AutoOpsRules[];
}

export interface AutoOpsRules {
  id: string;
  featureId: string;
  opsType: OpsType;
  clauses: {
    id: string;
    clause: unknown;
    actionType: 'UNKNOWN' | 'ENABLE' | 'DISABLE';
  }[];
  createdAt: string;
  updatedAt: string;
  deleted: boolean;
  autoOpsStatus: AutoOpsRuleStatus;
}

export interface GoalCollection {
  goals: Array<Goal>;
  cursor: string;
  totalCount: string;
}
