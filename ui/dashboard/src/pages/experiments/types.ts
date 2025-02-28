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
  filterByTab?: boolean;
  filterBySummary?: SummaryType;
}

export type SummaryType = 'scheduled' | 'running' | 'stopped' | 'not-started';

export type ExperimentActionsType =
  | 'EDIT'
  | 'STOP'
  | 'START'
  | 'ARCHIVE'
  | 'UNARCHIVE'
  | 'GOALS-CONNECTION';
export type ExperimentTab = 'ACTIVE' | 'ARCHIVED' | 'FINISHED';
