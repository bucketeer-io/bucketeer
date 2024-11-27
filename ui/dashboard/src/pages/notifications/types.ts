import { OrderBy, OrderDirection } from '@types';

export interface NotificationFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  disabled?: boolean;
  searchQuery: string;
}

export type NotificationActionsType = 'EDIT' | 'ENABLE' | 'DISABLE';
