import { SortingState } from '@tanstack/react-table';
import { OrderBy, OrderDirection, Organization } from '@types';

export interface OrganizationFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchQuery: string;
  disabled?: boolean;
}

export interface CollectionProps {
  isLoading?: boolean;
  onSortingChange: (v: SortingState) => void;
  organizations: Organization[];
}
