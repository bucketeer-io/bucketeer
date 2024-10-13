import { CollectionStatusType, OrderBy, OrderDirection } from '@types';

export interface OrganizationFilters {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  status: CollectionStatusType;
  searchQuery: string;
  disabled?: boolean;
}
