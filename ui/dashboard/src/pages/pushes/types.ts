import { OrderBy, OrderDirection } from '@types';

export interface PushFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  disabled?: boolean;
  searchQuery: string;
}

export type PushActionsType = 'EDIT' | 'ENABLE' | 'DISABLE';
