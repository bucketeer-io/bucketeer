import { SortingState } from '@tanstack/react-table';
import { OrderBy, OrderDirection, Account } from '@types';

export interface MembersFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  disabled?: boolean;
  searchQuery: string;
  organizationId?: string;
  organizationRole?: number;
  tags?: string[];
}

export interface CollectionProps {
  isLoading?: boolean;
  onSortingChange: (v: SortingState) => void;
  projects: Account[];
}

export type MemberActionsType =
  | 'EDIT'
  | 'DETAILS'
  | 'DELETE'
  | 'DISABLE'
  | 'ENABLE';
