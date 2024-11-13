import { SortingState } from '@tanstack/react-table';
import { CollectionStatusType, OrderBy, OrderDirection, Project } from '@types';

export interface EnvironmentFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchQuery: string;
  status: CollectionStatusType;
}

export interface CollectionProps {
  isLoading?: boolean;
  onSortingChange: (v: SortingState) => void;
  projects: Project[];
}

export type EnvironmentActionsType = 'EDIT' | 'ARCHIVE' | 'UNARCHIVE';
