import { CollectionStatusType, OrderBy, OrderDirection } from '@types';

export type GoalActions = 'ARCHIVE' | 'CONNECTION';

export interface GoalFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchQuery: string;
  status: CollectionStatusType;
  environmentId: string;
  isInUseStatus?: boolean;
  archived?: boolean;
}
