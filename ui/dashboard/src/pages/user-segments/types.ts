import { SortingState } from '@tanstack/react-table';
import { OrderBy, OrderDirection, Account } from '@types';

export type UserSegments = {
  id: string;
  name: string;
  description: string;
  rules: string;
  createdAt: string;
  updatedAt: string;
  version: string;
  deleted: boolean;
  includedUserCount: number;
  excludedUserCount: number;
  status: string;
  isInUseStatus: boolean;
  features: unknown;
  connections: number;
};

export interface UserSegmentsFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  disabled?: boolean;
  searchQuery: string;
  organizationId?: string;
}

export interface CollectionProps {
  isLoading?: boolean;
  onSortingChange: (v: SortingState) => void;
  projects: Account[];
}

export type UserSegmentsActionsType = 'EDIT' | 'DELETE' | 'FLAG';
