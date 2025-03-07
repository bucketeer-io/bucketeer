import { ExperimentStatus } from './experiment';

export type ConnectionType = 'UNKNOWN' | 'EXPERIMENT' | 'OPERATION';
export type OperationStatus = 'WAITING' | 'RUNNING' | 'STOPPED' | 'FINISHED';

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
  experiments: GoalExperiment[];
  autoOpsRules: GoalAutoOpsRule[];
}

export interface GoalAutoOpsRule {
  id: string;
  featureId: string;
  featureName: string;
  autoOpsStatus: OperationStatus;
}

export interface GoalExperiment {
  id: string;
  name: string;
  featureId: string;
  featureName: string;
  status: ExperimentStatus;
}

export interface GoalCollection {
  goals: Array<Goal>;
  cursor: string;
  totalCount: string;
}
