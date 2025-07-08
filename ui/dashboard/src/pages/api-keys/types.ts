import { APIKeyRole, OrderBy, OrderDirection } from '@types';

export interface APIKeysFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  disabled?: boolean;
  searchQuery: string;
  environment: string;
  environmentIds: string[];
}

export interface APIKeyOption {
  id: string;
  label: string;
  description: string;
  value: APIKeyRole;
}

export type APIKeyActionsType = 'EDIT' | 'ENABLE' | 'DISABLE';
