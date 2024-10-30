import { OrderBy, OrderDirection } from '@types';

export interface OrganizationProjectFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchQuery: string;
}

export interface OrganizationMembersFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchQuery: string;
  organizationId?: string;
}

export interface TabItem {
  readonly title: string;
  readonly to: string;
}
