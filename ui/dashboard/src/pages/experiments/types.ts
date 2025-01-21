import { OrderBy, OrderDirection } from '@types';

export interface ExperimentFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  disabled?: boolean;
  searchQuery: string;
}

export type ExperimentActionsType = 'EDIT' | 'STOP';
