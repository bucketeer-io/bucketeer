import { SortingState } from '@tanstack/react-table';
import { CollectionStatusType, OrderBy, OrderDirection, Project } from '@types';

export interface ProjectsFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchQuery: string;
  disabled?: boolean;
  status?: CollectionStatusType;
  organizationIds?: string[];
}

export interface CollectionProps {
  isLoading?: boolean;
  onSortingChange: (v: SortingState) => void;
  projects: Project[];
}
