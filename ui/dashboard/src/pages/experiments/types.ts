import { ExperimentStatus, OrderBy, OrderDirection } from '@types';

export interface ExperimentFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  archived?: boolean;
  searchQuery: string;
  featureId?: string;
  featureVersion?: number;
  from?: string;
  to?: string;
  status?: ExperimentTab;
  statuses?: ExperimentStatus[];
  maintainer?: string;
  isFilter?: boolean;
}

export type ExperimentActionsType =
  | 'EDIT'
  | 'STOP'
  | 'ARCHIVE'
  | 'UNARCHIVE'
  | 'GOALS-CONNECTION';
export type ExperimentTab = 'ACTIVE' | 'ARCHIVED' | 'FINISHED';
