import { SortingState } from '@tanstack/react-table';
import { OrderBy, OrderDirection, Project } from '@types';

export interface ProjectFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchQuery: string;
  disabled?: boolean;
}

export interface CollectionProps {
  isLoading?: boolean;
  onSortingChange: (v: SortingState) => void;
  projects: Project[];
}
