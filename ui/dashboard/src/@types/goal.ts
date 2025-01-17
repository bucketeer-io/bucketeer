export type ConnectionType = 'UNKNOWN' | 'EXPERIMENT' | 'OPERATION';
export type ExperimentStatus =
  | 'WAITING'
  | 'RUNNING'
  | 'STOPPED'
  | 'FORCE_STOPPED';
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

export interface Experiment {
  id: string;
  goalId: string;
  featureId: string;
  featureVersion: number;
  variations: [
    {
      id: string;
      value: string;
      name: string;
      description: string;
    }
  ];
  startAt: string;
  stopAt: string;
  stopped: boolean;
  stoppedAt: string;
  createdAt: string;
  updatedAt: string;
  deleted: boolean;
  goalIds: string[];
  name: string;
  description: string;
  baseVariationId: string;
  status: ExperimentStatus;
  maintainer: string;
  archived: boolean;
}

export interface GoalCollection {
  goals: Array<Goal>;
  cursor: string;
  totalCount: string;
}
