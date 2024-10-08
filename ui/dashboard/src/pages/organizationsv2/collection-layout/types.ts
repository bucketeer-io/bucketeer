import { SortingState } from '@tanstack/react-table';
import { Organization } from '@types';

export interface ItemProps {
  organization: Organization;
}

// export interface ItemLayoutProps {
// 	layout: CollectionLayoutType;
// 	organization: Organization;
// }

export interface CollectionProps {
  isLoading?: boolean;
  onSortingChange: (v: SortingState) => void;
  organizations: Organization[];
}
