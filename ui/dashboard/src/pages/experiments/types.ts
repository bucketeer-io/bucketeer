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
  statuses?: ExperimentStatus;
  maintainer?: string;
}

export type ExperimentActionsType = 'EDIT' | 'STOP' | 'ARCHIVE' | 'UNARCHIVE';
export type ExperimentTab = 'ACTIVE' | 'ARCHIVED' | 'COMPLETE';
