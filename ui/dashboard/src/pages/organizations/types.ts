import { SortingState } from '@tanstack/react-table';
import {
  CollectionStatusType,
  OrderBy,
  OrderDirection,
  Organization
} from '@types';

export interface OrganizationFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  status: CollectionStatusType;
  searchQuery: string;
  disabled: boolean;
}

export interface CollectionProps {
  isLoading?: boolean;
  onSortingChange: (v: SortingState) => void;
  organizations: Organization[];
}

export type OrganizationActionsType = 'EDIT' | 'ARCHIVE' | 'UNARCHIVE';
