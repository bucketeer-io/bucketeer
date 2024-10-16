import { OrderBy, OrderDirection } from '@types';

export interface OrganizationProjectFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchQuery: string;
}

export interface OrganizationUsersFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchQuery: string;
  organizationId?: string;
}
