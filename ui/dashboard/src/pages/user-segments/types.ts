import { SortingState } from '@tanstack/react-table';
import { OrderBy, OrderDirection, Account, FeatureSegmentStatus } from '@types';

export interface UserSegmentsFilters {
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  disabled?: boolean;
  status?: FeatureSegmentStatus;
  isInUseStatus?: boolean;
  environmentId: string;
}

export interface CollectionProps {
  isLoading?: boolean;
  onSortingChange: (v: SortingState) => void;
  projects: Account[];
}

export type UserSegmentsActionsType = 'EDIT' | 'DELETE' | 'DOWNLOAD' | 'FLAG';

export type UserSegmentForm = {
  id?: string;
  name: string;
  description?: string;
  userIds?: string;
  file?: unknown;
};
