import { OrderBy, OrderDirection } from '@types';

export interface APIKeysFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  disabled?: boolean;
  searchQuery: string;
}

export type APIKeyActionsType = 'EDIT' | 'ENABLE' | 'DISABLE';
