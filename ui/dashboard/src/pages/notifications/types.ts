import { OrderBy, OrderDirection, SourceType } from '@types';

export interface NotificationFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  disabled?: boolean;
  searchQuery: string;
}

export type NotificationActionsType = 'EDIT' | 'ENABLE' | 'DISABLE' | 'DELETE';

export interface NotificationOption {
  value: SourceType;
  label: string;
  description: string;
}
