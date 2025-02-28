export type ExperimentStatus =
  | 'WAITING'
  | 'RUNNING'
  | 'STOPPED'
  | 'FORCE_STOPPED'
  | 'NOT_STARTED';

export interface Experiment {
  id: string;
  goalId: string;
  featureId: string;
  featureVersion: number;
  variations: {
    id: string;
    value: string;
    name: string;
    description: string;
  }[];
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
  goals: {
    id: string;
    name: string;
  }[];
}

export interface ExperimentCollection {
  experiments: Array<Experiment>;
  cursor: string;
  totalCount: string;
  summary: {
    totalWaitingCount: string;
    totalRunningCount: string;
    totalStoppedCount: string;
  };
}
