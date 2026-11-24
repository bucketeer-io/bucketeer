import { FeatureVariation } from './feature';

export type ExperimentStatus =
  'WAITING' | 'RUNNING' | 'STOPPED' | 'FORCE_STOPPED' | 'NOT_STARTED';

// HOLDOUT is an internal experiment that measures a holdout group.
// Its featureId is the holdout group ID, not a feature flag ID.
export type ExperimentKind = 'KIND_EXPERIMENT' | 'HOLDOUT';

export interface Experiment {
  id: string;
  goalId: string;
  featureId: string;
  featureVersion: number;
  variations: FeatureVariation[];
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
  kind: ExperimentKind;
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
